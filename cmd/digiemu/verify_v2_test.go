package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func compileV2Schema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	absSchemaPath, err := filepath.Abs(filepath.Join(repoRoot, "schemas", "verify_result_v2.schema.json"))
	if err != nil {
		t.Fatalf("abs schema path: %v", err)
	}
	compiler := jsonschema.NewCompiler()
	schemaURL := "file:///" + filepath.ToSlash(absSchemaPath)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile v2 schema: %v", err)
	}
	return sch
}

// writeDemoBundleV2 creates a minimal snapshot bundle under tmp/snapshots/<ref>
// with the given expected_hash_v1 value.
func writeDemoBundleV2(t *testing.T, tmp, ref, expected string) string {
	t.Helper()
	root := filepath.Join(tmp, "snapshots", ref)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := "{\n" +
		"  \"expected_hash_v1\": " + strconv.Quote(expected) + ",\n" +
		"  \"foo\": 1\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join(root, "snapshot.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write snapshot.json: %v", err)
	}
	return root
}

func TestRunVerify_V2Flag_EmitsV2Document(t *testing.T) {
	tmp := t.TempDir()
	writeDemoBundleV2(t, tmp, "demo", "NOPE")

	args := []string{
		"--ref", "demo",
		"--data", tmp,
		"--fixture-root", tmp,
		"--verify-result-version=v2-draft",
	}

	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("expected exit=1 (verify_fail), got %d\nstderr=%s", exit, stderr.String())
	}

	var doc map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("stdout not valid json: %v\nstdout=%q", err, stdout.String())
	}

	// v2 required fields
	for _, k := range []string{"verify_result_version", "result", "reason_code", "canonicalization_profile", "hash_algorithm", "reconstruction_profile", "verifier_version"} {
		if _, ok := doc[k]; !ok {
			t.Fatalf("missing required v2 field %q in %v", k, doc)
		}
	}
	if doc["verify_result_version"] != "v2-draft" {
		t.Fatalf("verify_result_version=%v", doc["verify_result_version"])
	}
	if doc["result"] != "FAIL" || doc["reason_code"] != "HASH_MISMATCH" {
		t.Fatalf("expected FAIL/HASH_MISMATCH, got %v/%v", doc["result"], doc["reason_code"])
	}

	// v1-only fields must not leak into v2 output
	for _, k := range []string{"ok", "ref", "expected", "got", "hash_alg", "canonical_scope", "trace", "wrote_expected", "write_blocked", "write_reason"} {
		if _, ok := doc[k]; ok {
			t.Fatalf("v1 field %q leaked into v2 output", k)
		}
	}
}

func TestRunVerify_V2Flag_Alias20Draft(t *testing.T) {
	tmp := t.TempDir()
	writeDemoBundleV2(t, tmp, "demo", "NOPE")

	args := []string{
		"--ref", "demo",
		"--data", tmp,
		"--fixture-root", tmp,
		"--verify-result-version=2.0-draft",
	}

	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("expected exit=1, got %d", exit)
	}

	var doc map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("stdout not valid json: %v", err)
	}
	// Compatibility alias must normalize to the canonical identifier.
	if doc["verify_result_version"] != "v2-draft" {
		t.Fatalf("expected canonical v2-draft, got %v", doc["verify_result_version"])
	}
}

func TestRunVerify_V2Flag_ShortV2Accepted(t *testing.T) {
	tmp := t.TempDir()
	writeDemoBundleV2(t, tmp, "demo", "NOPE")

	args := []string{
		"--ref", "demo",
		"--data", tmp,
		"--fixture-root", tmp,
		"--verify-result-version=v2",
	}

	var stdout, stderr bytes.Buffer
	if exit := runVerifyWithIO(args, &stdout, &stderr); exit != 1 {
		t.Fatalf("expected exit=1, got %d", exit)
	}
	var doc map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("stdout not valid json: %v", err)
	}
	if doc["verify_result_version"] != "v2-draft" {
		t.Fatalf("expected v2-draft, got %v", doc["verify_result_version"])
	}
}

func TestRunVerify_V2Flag_InvalidValueUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO([]string{"--verify-result-version=v3"}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("expected exit=2 (usage), got %d", exit)
	}
	if stdout.String() != "" {
		t.Fatalf("expected empty stdout on usage error")
	}
	if !strings.Contains(stderr.String(), "verify-result-version") {
		t.Fatalf("expected stderr to mention the flag, got %q", stderr.String())
	}
}

func TestRunVerify_V2Flag_PASS_IncludesSnapshotHash(t *testing.T) {
	// PASS path needs a real matching expected hash. Build via the demo fixture:
	// use --write-expected style snapshot? Simpler: compute via the operator is
	// overkill; reuse an existing passing fixture.
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO([]string{"--bundle", bundle, "--verify-result-version=v2-draft"}, &stdout, &stderr)
	if exit != 0 {
		t.Fatalf("expected exit=0, got %d\nstderr=%s", exit, stderr.String())
	}

	var doc map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("stdout not valid json: %v", err)
	}
	if doc["result"] != "PASS" || doc["reason_code"] != "STATE_RECONSTRUCTED" {
		t.Fatalf("expected PASS/STATE_RECONSTRUCTED, got %v/%v", doc["result"], doc["reason_code"])
	}
	sh, ok := doc["snapshot_hash"].(string)
	if !ok || sh == "" {
		t.Fatalf("PASS v2 result missing snapshot_hash")
	}
	if !strings.HasPrefix(sh, "sha256:") {
		t.Fatalf("expected sha256:-tagged snapshot_hash, got %q", sh)
	}
}

func TestRunVerify_V2Flag_OutputValidatesV2Schema(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	cases := []struct {
		name string
		args []string
	}{
		{"pass pretty", []string{"--bundle", bundle, "--verify-result-version=v2-draft"}},
		{"pass canonical", []string{"--bundle", bundle, "--json=canonical", "--verify-result-version=v2-draft"}},
	}

	sch := compileV2Schema(t)
	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		exit := runVerifyWithIO(tc.args, &stdout, &stderr)
		if exit != 0 {
			t.Fatalf("%s: expected exit=0, got %d\nstderr=%s", tc.name, exit, stderr.String())
		}
		var doc map[string]any
		if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
			t.Fatalf("%s: stdout not valid json: %v", tc.name, err)
		}
		if err := sch.Validate(doc); err != nil {
			t.Fatalf("%s: emitted v2 document failed schema validation: %v\nstdout=%s", tc.name, err, stdout.String())
		}
	}
}

func TestRunVerify_V2Flag_CanonicalModeEmitsCanonicalJSON(t *testing.T) {
	tmp := t.TempDir()
	writeDemoBundleV2(t, tmp, "demo", "NOPE")

	args := []string{
		"--ref", "demo",
		"--data", tmp,
		"--fixture-root", tmp,
		"--json=canonical",
		"--verify-result-version=v2-draft",
	}

	var stdout, stderr bytes.Buffer
	if exit := runVerifyWithIO(args, &stdout, &stderr); exit != 1 {
		t.Fatalf("expected exit=1, got %d", exit)
	}
	out := stdout.String()
	if strings.Contains(out, "\n  ") {
		t.Fatalf("expected compact canonical JSON, got pretty output")
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("stdout not valid json: %v", err)
	}
	if doc["verify_result_version"] != "v2-draft" {
		t.Fatalf("expected v2-draft, got %v", doc["verify_result_version"])
	}
}
