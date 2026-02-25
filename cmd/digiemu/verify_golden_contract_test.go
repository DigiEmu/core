package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func buildDigiemuBinaryTempDir(t *testing.T, repoRoot string) string {
	t.Helper()

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	out := filepath.Join(t.TempDir(), "digiemu-test"+ext)

	cmd := exec.Command("go", "build", "-o", out, "./cmd/digiemu")
	cmd.Dir = repoRoot
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build cli: %v\n%s", err, string(b))
	}

	return out
}

func normalizeVerifyGolden(t *testing.T, repoRoot string, raw []byte) []byte {
	t.Helper()

	var v verifyJSON
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("invalid verify json: %v\nraw=%s", err, string(raw))
	}

	v.Trace = stableTrace(t, repoRoot, v.Trace)

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("marshal normalized verify json: %v", err)
	}
	return append(b, '\n')
}

func TestVerify_JSON_Golden(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	bin := buildDigiemuBinaryTempDir(t, repoRoot)

	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code, stdout, stderr := runCLI(t, bin, "verify", "--bundle", bundle, "--json")
	if code != 0 {
		t.Fatalf("verify failed: exit=%d\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}

	var parsed verifyJSON
	if err := json.Unmarshal([]byte(stdout), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v\nstdout=%s", err, stdout)
	}
	if !parsed.OK {
		t.Fatalf("expected ok=true\nstdout=%s", stdout)
	}
	if parsed.Ref != "demo" {
		t.Fatalf("expected ref=demo, got %q\nstdout=%s", parsed.Ref, stdout)
	}
	if parsed.Expected == "" || parsed.Got == "" {
		t.Fatalf("expected expected/got to be present\nstdout=%s", stdout)
	}
	if parsed.Expected != parsed.Got {
		t.Fatalf("expected expected==got\nstdout=%s", stdout)
	}
	if parsed.Scope != "canonical_json_v1_excluding_expected_hash_v1" {
		t.Fatalf("unexpected canonical_scope=%q\nstdout=%s", parsed.Scope, stdout)
	}
	if parsed.HashAlg != "sha256(canonical_json_v1)" {
		t.Fatalf("unexpected hash_alg=%q\nstdout=%s", parsed.HashAlg, stdout)
	}

	got := normalizeVerifyGolden(t, repoRoot, []byte(stdout))

	goldenPath := filepath.Join(repoRoot, "cmd", "digiemu", "testdata", "verify_demo_golden.json")
	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	want := normalizeVerifyGolden(t, repoRoot, golden)

	if !bytes.Equal(want, got) {
		t.Fatalf("verify output changed — golden mismatch\n--- want ---\n%s\n--- got ---\n%s", string(want), string(got))
	}
}
