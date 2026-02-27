package repro_demo_v1

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf16"
)

func readExpectedFromGitOrDisk(t *testing.T, root string, rel string) []byte {
	t.Helper()

	// Prefer reading from the git blob to avoid platform-specific working tree
	// newline conversion (e.g. core.autocrlf). This keeps the byte-equality test
	// stable across OSes.
	cmd := exec.Command("git", "show", "HEAD:"+rel)
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		return out
	}

	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read expected %s: %v", rel, err)
	}
	return b
}

func encodeUTF16LEWithBOM(s string) []byte {
	u16s := utf16.Encode([]rune(s))
	out := make([]byte, 2+len(u16s)*2)
	out[0] = 0xFF
	out[1] = 0xFE
	for i, v := range u16s {
		binary.LittleEndian.PutUint16(out[2+i*2:2+i*2+2], v)
	}
	return out
}

func normalizeStdoutToExpectedBytes(stdout []byte, expected []byte) []byte {
	// If expected is UTF-16LE with BOM (common when created via PowerShell redirection),
	// normalize stdout by converting LF->CRLF and encoding to UTF-16LE+BOM.
	if len(expected) >= 2 && expected[0] == 0xFF && expected[1] == 0xFE {
		s := string(stdout) // stdout is UTF-8 JSON from Go
		s = strings.ReplaceAll(s, "\n", "\r\n")
		return encodeUTF16LEWithBOM(s)
	}

	// If expected uses CRLF, normalize stdout line endings.
	if bytes.Contains(expected, []byte("\r\n")) && bytes.Contains(stdout, []byte("\n")) {
		return []byte(strings.ReplaceAll(string(stdout), "\n", "\r\n"))
	}

	return stdout
}

func TestVerify_ReproducesExpectedReport_ByteExact(t *testing.T) {
	root := repoRoot(t)

	expected := readExpectedFromGitOrDisk(t, root, "examples/snapshot_v1_demo/expected_verify_report.json")

	// Run: go run ./cmd/digiemu verify --bundle <bundle> --json
	// We capture stdout and compare bytes exactly.
	cmd := exec.Command("go", "run", "./cmd/digiemu", "verify",
		"--bundle", "examples/snapshot_v1_demo/snapshots/snapshot_demo_v1",
		"--json",
	)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		// If CLI exits non-zero but still wrote JSON to stdout, we still want to compare bytes.
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("verify failed: %v", err)
		}
		// continue — out may contain stdout even on ExitError
	}

	norm := normalizeStdoutToExpectedBytes(out, expected)
	if !bytes.Equal(norm, expected) {
		t.Fatalf("verify output drifted: stdout != expected_verify_report.json (byte mismatch)")
	}
}

func TestVerify_ReproducesExpectedHash(t *testing.T) {
	// Hash is stored in report field "got" (current behavior observed).
	// This test only ensures the stored hash matches the report's got.
	// If report format changes, update accordingly (but that would be a contract change).
	root := repoRoot(t)

	hashPath := filepath.Join(root, "examples/snapshot_v1_demo/expected_snapshot_hash.txt")
	hashBytes, err := os.ReadFile(hashPath)
	if err != nil {
		t.Fatalf("read expected hash: %v", err)
	}

	expectedHash := bytes.TrimSpace(decodeTextToUTF8(t, hashBytes))

	cmd := exec.Command("go", "run", "./cmd/digiemu", "verify",
		"--bundle", "examples/snapshot_v1_demo/snapshots/snapshot_demo_v1",
		"--json",
	)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			t.Fatalf("verify failed: %v", err)
		}
		// continue — out may contain stdout even on ExitError
	}
	if len(out) == 0 {
		t.Fatalf("verify produced no stdout JSON")
	}

	// Parse "got" from JSON output and compare to expected_snapshot_hash.txt
	// This is part of the public contract: verify result contains "got" as the computed snapshot hash.
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("decode verify output JSON: %v", err)
	}
	gotVal, ok := m["got"]
	if !ok {
		t.Fatalf("verify output missing field %q", "got")
	}
	got, ok := gotVal.(string)
	if !ok {
		t.Fatalf("verify output field %q must be string", "got")
	}
	if !bytes.Equal([]byte(got), expectedHash) {
		t.Fatalf("snapshot hash mismatch: got=%q expected=%q", got, string(expectedHash))
	}
}
