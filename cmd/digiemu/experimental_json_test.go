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
	if int(rep["total"].(float64)) != 3 {
		t.Fatalf("expected total 3, got %v", rep["total"])
	}
	if int(rep["passed"].(float64)) != 3 {
		t.Fatalf("expected passed 3, got %v", rep["passed"])
	}
	if int(rep["failed"].(float64)) != 0 {
		t.Fatalf("expected failed 0, got %v", rep["failed"])
	}
	// cases array
	cases, ok := rep["cases"].([]any)
	if !ok || len(cases) != 3 {
		t.Fatalf("expected 3 cases, got %v", rep["cases"])
	}
}
