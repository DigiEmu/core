package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestVerifyV2_ObservedHashSurvivesExpectedHashMutation demonstrates that the
// existing verifier recomputes the observed snapshot hash independently of the
// declared expected_hash_v1: mutating only expected_hash_v1 in a copied bundle
// flips the result to FAIL/HASH_MISMATCH while snapshot_hash_actual remains
// identical. This does not establish an independent verifier implementation.
func TestVerifyV2_ObservedHashSurvivesExpectedHashMutation(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	srcBundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")
	srcSnapshot := filepath.Join(srcBundle, "snapshot.json")

	srcBefore, err := os.ReadFile(srcSnapshot)
	if err != nil {
		t.Fatalf("read source snapshot: %v", err)
	}
	t.Cleanup(func() {
		after, err := os.ReadFile(srcSnapshot)
		if err != nil {
			t.Errorf("re-read source snapshot: %v", err)
			return
		}
		if !bytes.Equal(srcBefore, after) {
			t.Errorf("repository fixture snapshot.json was modified during test")
		}
	})

	tmp := t.TempDir()
	bundle := filepath.Join(tmp, "snapshots", "demo")
	copyDir(t, srcBundle, bundle)

	args := []string{"--bundle", bundle, "--verify-result-version=v2-draft", "--json=canonical"}

	// Baseline: unmodified copy must verify as PASS.
	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("baseline: expected exit=0, got %d\nstderr=%s", exit, stderr.String())
	}
	var pass map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &pass); err != nil {
		t.Fatalf("baseline: stdout not valid json: %v", err)
	}
	if pass["result"] != "PASS" || pass["reason_code"] != "STATE_RECONSTRUCTED" {
		t.Fatalf("baseline: expected PASS/STATE_RECONSTRUCTED, got %v/%v", pass["result"], pass["reason_code"])
	}
	baselineHash, ok := pass["snapshot_hash"].(string)
	if !ok || !strings.HasPrefix(baselineHash, "sha256:") {
		t.Fatalf("baseline: missing or malformed snapshot_hash: %v", pass["snapshot_hash"])
	}

	// Mutate only expected_hash_v1 in the temporary copy.
	mutated := strings.Repeat("0", 64)
	snapshotPath := filepath.Join(bundle, "snapshot.json")
	b, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("read temp snapshot: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal temp snapshot: %v", err)
	}
	obj["expected_hash_v1"] = mutated
	nb, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatalf("marshal temp snapshot: %v", err)
	}
	if err := os.WriteFile(snapshotPath, nb, 0o644); err != nil {
		t.Fatalf("write temp snapshot: %v", err)
	}

	// Re-verify: the independently computed hash must survive the mutation.
	stdout.Reset()
	stderr.Reset()
	exit = runVerifyWithIO(args, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("mutated: expected exit=1, got %d\nstderr=%s", exit, stderr.String())
	}
	var fail map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &fail); err != nil {
		t.Fatalf("mutated: stdout not valid json: %v", err)
	}
	if fail["result"] != "FAIL" || fail["reason_code"] != "HASH_MISMATCH" {
		t.Fatalf("mutated: expected FAIL/HASH_MISMATCH, got %v/%v", fail["result"], fail["reason_code"])
	}
	if fail["snapshot_hash_actual"] != baselineHash {
		t.Fatalf("mutated: snapshot_hash_actual=%v, want baseline %v", fail["snapshot_hash_actual"], baselineHash)
	}
	if want := "sha256:" + mutated; fail["snapshot_hash_expected"] != want {
		t.Fatalf("mutated: snapshot_hash_expected=%v, want %v", fail["snapshot_hash_expected"], want)
	}
	if fail["failure_path"] != "/expected_hash_v1" {
		t.Fatalf("mutated: failure_path=%v, want /expected_hash_v1", fail["failure_path"])
	}
}
