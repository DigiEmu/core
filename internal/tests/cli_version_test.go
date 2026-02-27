package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// internal/tests -> repo root = ../../
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestCLIVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/digiemu", "--version")
	cmd.Dir = repoRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v\n%s", err, string(out))
	}
	s := strings.TrimSpace(string(out))
	if !strings.HasPrefix(s, "digiemu ") {
		t.Fatalf("unexpected version output: %q", s)
	}
}
