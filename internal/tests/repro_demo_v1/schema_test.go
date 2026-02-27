package repro_demo_v1

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"unicode/utf16"
)

// NOTE:
// We intentionally keep this test dependency-free.
// The schema file is authoritative, but we validate it structurally here
// and then do a strict JSON parse + minimal invariants.
//
// Next YAML will replace this with a real JSON Schema validator IF the repo already has one.
// If not, we keep minimal invariants as the stable guardrail.

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// internal/tests/repro_demo_v1 -> repo root = ../../../
	root := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	return root
}

func readFile(t *testing.T, rel string) []byte {
	t.Helper()
	p := filepath.Join(repoRoot(t), rel)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return b
}

func decodeTextToUTF8(t *testing.T, b []byte) []byte {
	t.Helper()

	// Handle UTF-16LE with BOM (common on Windows when created via PowerShell redirection).
	if len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE {
		if (len(b)-2)%2 != 0 {
			t.Fatalf("utf-16le: odd byte length")
		}
		u16s := make([]uint16, 0, (len(b)-2)/2)
		for i := 2; i < len(b); i += 2 {
			u16s = append(u16s, binary.LittleEndian.Uint16(b[i:i+2]))
		}
		runes := utf16.Decode(u16s)
		return []byte(string(runes))
	}

	// Strip UTF-8 BOM if present.
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}

	return b
}

func TestExpectedVerifyReport_IsValidJSON(t *testing.T) {
	b := decodeTextToUTF8(t, readFile(t, "examples/snapshot_v1_demo/expected_verify_report.json"))
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // strictness: will fail if we decode into a struct. we decode into map below.
	var any interface{}
	// Can't use DisallowUnknownFields with interface{}, so just decode.
	dec = json.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&any); err != nil {
		t.Fatalf("expected_verify_report.json is not valid JSON: %v", err)
	}
}

func TestVerifyResultSchema_IsPresentAndValidJSON(t *testing.T) {
	b := decodeTextToUTF8(t, readFile(t, "schemas/VERIFY_RESULT_SCHEMA_v1.json"))
	var any interface{}
	if err := json.Unmarshal(b, &any); err != nil {
		t.Fatalf("VERIFY_RESULT_SCHEMA_v1.json is not valid JSON: %v", err)
	}
}

// Minimal invariants for VerifyResult v1 (must match current schema semantics).
// This is a guardrail until/if we wire a full JSON Schema validator using an existing internal package.
func TestExpectedVerifyReport_MinimalInvariants(t *testing.T) {
	b := decodeTextToUTF8(t, readFile(t, "examples/snapshot_v1_demo/expected_verify_report.json"))
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode report: %v", err)
	}

	// required fields
	mustStr := []string{"ref", "expected", "got", "hash_alg", "canonical_scope", "message"}
	for _, k := range mustStr {
		v, ok := m[k]
		if !ok {
			t.Fatalf("missing field %q", k)
		}
		if _, ok := v.(string); !ok {
			t.Fatalf("field %q must be string", k)
		}
	}

	// ok boolean
	if v, ok := m["ok"]; !ok {
		t.Fatalf("missing field %q", "ok")
	} else {
		if _, ok := v.(bool); !ok {
			t.Fatalf("field %q must be bool", "ok")
		}
	}

	// trace: allow array or object (schema-dependent)
	if _, ok := m["trace"]; !ok {
		t.Fatalf("missing field %q", "trace")
	}
}

func TestExpectedVerifyReport_ConformsToSchemaSubset(t *testing.T) {
	schema := decodeTextToUTF8(t, readFile(t, "schemas/VERIFY_RESULT_SCHEMA_v1.json"))
	doc := decodeTextToUTF8(t, readFile(t, "examples/snapshot_v1_demo/expected_verify_report.json"))
	if err := ValidateAgainstSchemaSubset(schema, doc); err != nil {
		t.Fatalf("expected report fails schema subset gate: %v", err)
	}
}

func TestRuntimeVerifyOutput_ConformsToSchemaSubset(t *testing.T) {
	schema := decodeTextToUTF8(t, readFile(t, "schemas/VERIFY_RESULT_SCHEMA_v1.json"))

	root := repoRoot(t)
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
	if err := ValidateAgainstSchemaSubset(schema, out); err != nil {
		t.Fatalf("runtime verify output fails schema subset gate: %v", err)
	}
}
