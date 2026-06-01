package main

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExperimentalTopLevelDispatch(t *testing.T) {
	t.Parallel()

	// run: go run ./cmd/digiemu experimental conformance run <path>
	repoPath := filepath.Join("..", "..", "testdata", "core_2_conformance", "basic_pass")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// run in-package using 'go run .' since tests execute from the package dir
	cmd := exec.CommandContext(ctx, "go", "run", ".", "experimental", "conformance", "run", repoPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v\noutput:\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "Conformance run summary") {
		t.Fatalf("unexpected output, missing summary: %s", string(out))
	}
}
