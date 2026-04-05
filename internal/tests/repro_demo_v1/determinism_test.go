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

func readExpectedFromDisk(t *testing.T, root string, rel string) []byte {
	t.Helper()

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

func trimSingleTrailingLFOrCRLF(b []byte) []byte {
	if len(b) >= 2 && b[len(b)-2] == '\r' && b[len(b)-1] == '\n' {
		return b[:len(b)-2]
	}
	if len(b) >= 1 && b[len(b)-1] == '\n' {
		return b[:len(b)-1]
	}
	return b
}

func normalizeStdoutToExpectedBytes(stdout []byte, expected []byte) []byte {
	// Normalize a single trailing newline difference, because fixtures may be
	// saved with or without final newline depending on editor/platform.
	stdout = trimSingleTrailingLFOrCRLF(stdout)
	expected = trimSingleTrailingLFOrCRLF(expected)

	// If expected is UTF-16LE with BOM, normalize stdout accordingly.
	if len(expected) >= 2 && expected[0] == 0xFF && expected[1] == 0xFE {
		s := string(stdout)
		s = strings.ReplaceAll(s, "\n", "\r\n")
		return encodeUTF16LEWithBOM(s)
	}

	// If expected uses CRLF, normalize stdout line endings.
	if bytes.Contains(expected, []byte("\r\n")) && bytes.Contains(stdout, []byte("\n")) {
		return []byte(strings.ReplaceAll(string(stdout), "\n", "\r\n"))
	}

	return stdout
}

func firstDiffIndex(a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(a) != len(b) {
		return n
	}
	return -1
}

func TestVerify_ReproducesExpectedReport_ByteExact(t *testing.T) {
	root := repoRoot(t)

	expected := readExpectedFromDisk(t, root, "examples/snapshot_v1_demo/expected_verify_report.json")

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
		// continue: verify may intentionally exit non-zero while still emitting JSON
	}

	norm := normalizeStdoutToExpectedBytes(out, expected)
	expNorm := trimSingleTrailingLFOrCRLF(expected)

	if !bytes.Equal(norm, expNorm) {
		i := firstDiffIndex(norm, expNorm)
		t.Fatalf(
			"verify output drifted: stdout != expected_verify_report.json (byte mismatch; idx=%d, got_len=%d, expected_len=%d)",
			i, len(norm), len(expNorm),
		)
	}
}

func TestVerify_ReproducesExpectedHash(t *testing.T) {
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
	}
	if len(out) == 0 {
		t.Fatalf("verify produced no stdout JSON")
	}

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
