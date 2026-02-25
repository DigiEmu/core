package main

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestVerify_JSON_ConformsToSchema(t *testing.T) {
	bin := buildBinaryOnce(t)

	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code, stdout, stderr := runVerifyCLI(t, bin, "--bundle", bundle, "--json")
	if code != 0 {
		t.Fatalf("verify failed: exit=%d\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(stdout), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v\nstdout=%s", err, stdout)
	}

	required := []string{"ok", "ref", "expected", "got", "hash_alg", "canonical_scope", "trace", "message"}
	for _, k := range required {
		if _, ok := parsed[k]; !ok {
			t.Fatalf("missing required field: %s", k)
		}
	}
}
