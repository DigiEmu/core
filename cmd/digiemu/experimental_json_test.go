package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestExperimentalJSONOutput(t *testing.T) {
	t.Parallel()

	repoPath := filepath.Join("..", "..", "testdata", "core_2_conformance")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".", "experimental", "conformance", "run", repoPath, "--json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\noutput:\n%s", err, string(out))
	}

	var rep map[string]any
	// the test harness and other packages may emit text; find the JSON object
	s := string(out)
	i := -1
	for idx, ch := range s {
		if ch == '{' {
			i = idx
			break
		}
	}
	if i < 0 {
		t.Fatalf("no JSON object found in output:\n%s", s)
	}
	if err := json.Unmarshal([]byte(s[i:]), &rep); err != nil {
		t.Fatalf("invalid json output: %v\njson fragment:\n%s\nfull output:\n%s", err, s[i:], s)
	}
	// basic assertions
	if rep["status"] != "PASS" {
		t.Fatalf("expected status PASS, got %v", rep["status"])
	}
	if int(rep["total"].(float64)) != 10 {
		t.Fatalf("expected total 10, got %v", rep["total"])
	}
	if int(rep["passed"].(float64)) != 10 {
		t.Fatalf("expected passed 10, got %v", rep["passed"])
	}
	if int(rep["failed"].(float64)) != 0 {
		t.Fatalf("expected failed 0, got %v", rep["failed"])
	}
	// cases array
	cases, ok := rep["cases"].([]any)
	if !ok || len(cases) != 10 {
		t.Fatalf("expected 10 cases, got %v", rep["cases"])
	}

	// Optional: assert presence of known case names
	want := map[string]bool{
		"basic_pass":                       true,
		"hash_mismatch_fail":               true,
		"invalid_schema_error":             true,
		"missing_required_field_error":     true,
		"malformed_json_error":             true,
		"unknown_reason_code_error":        true,
		"unsupported_hash_algorithm_error": true,
		"wrong_profile_fail":               true,
		"inside_payload_mutation_detected": true,
		"outside_metadata_ignored":         true,
	}
	for _, c := range cases {
		if m, ok := c.(map[string]any); ok {
			if name, ok := m["name"].(string); ok {
				delete(want, name)
			}
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing expected cases: %v", want)
	}
}
