package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func buildDigiemuBinary(t *testing.T) string {
	t.Helper()

	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	outDir := t.TempDir()
	exe := "digiemu-test"
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	binPath := filepath.Join(outDir, exe)

	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/digiemu")
	cmd.Dir = repoRoot
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build digiemu-test binary: %v: %s", err, stderr.String())
	}

	return binPath
}

func runCLI(t *testing.T, binPath string, args ...string) (exitCode int, stdout string, stderr string) {
	t.Helper()

	cmd := exec.Command(binPath, args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	stdout = outb.String()
	stderr = errb.String()

	if err == nil {
		return 0, stdout, stderr
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), stdout, stderr
	}

	t.Fatalf("run cli: %v\nstdout=%s\nstderr=%s", err, stdout, stderr)
	return 0, stdout, stderr
}

func TestReplayDeterminismJSON(t *testing.T) {
	bin := buildDigiemuBinary(t)
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code1, out1, err1 := runCLI(t, bin, "replay", "--bundle", bundle, "--json")
	if code1 != 0 {
		t.Fatalf("exit code=%d, want 0\nstdout=%s\nstderr=%s", code1, out1, err1)
	}

	code2, out2, err2 := runCLI(t, bin, "replay", "--bundle", bundle, "--json")
	if code2 != 0 {
		t.Fatalf("exit code=%d, want 0\nstdout=%s\nstderr=%s", code2, out2, err2)
	}

	if out1 != out2 {
		t.Fatalf("replay output not deterministic (stdout differs)\n--- run1 ---\n%s\n--- run2 ---\n%s", out1, out2)
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(out1), &obj); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nstdout=%s", err, out1)
	}

	traceAny, ok := obj["trace"].([]any)
	if !ok {
		t.Fatalf("trace missing or not an array\nstdout=%s", out1)
	}

	foundUsed := false
	for _, v := range traceAny {
		s, ok := v.(string)
		if !ok {
			t.Fatalf("trace element is not a string\nstdout=%s", out1)
		}
		if strings.Contains(s, "used:") {
			foundUsed = true
		}
	}
	if !foundUsed {
		t.Fatalf("trace does not contain used: marker\nstdout=%s", out1)
	}
}
