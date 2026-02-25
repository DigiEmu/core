package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type verifyJSON struct {
	OK            bool     `json:"ok"`
	Ref           string   `json:"ref"`
	Expected      string   `json:"expected"`
	Got           string   `json:"got"`
	HashAlg       string   `json:"hash_alg"`
	Scope         string   `json:"canonical_scope"`
	Trace         []string `json:"trace"`
	Message       string   `json:"message"`
	WroteExpected bool     `json:"wrote_expected"`
	WriteBlocked  bool     `json:"write_blocked"`
	WriteReason   string   `json:"write_reason"`
}

var (
	buildOnce sync.Once
	builtBin  string
	buildErr  error
)

func buildBinaryOnce(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		repoRoot, err := findRepoRoot()
		if err != nil {
			buildErr = err
			return
		}

		tmp := os.TempDir()
		outDir, err := os.MkdirTemp(tmp, "digiemu-cli-")
		if err != nil {
			buildErr = err
			return
		}

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
			buildErr = errors.New(err.Error() + ": " + stderr.String())
			return
		}

		builtBin = binPath
	})

	if buildErr != nil {
		t.Fatalf("build digiemu-test binary: %v", buildErr)
	}
	return builtBin
}

func runVerifyCLI(t *testing.T, binPath string, args ...string) (exitCode int, stdout string, stderr string) {
	t.Helper()

	cmd := exec.Command(binPath, append([]string{"verify"}, args...)...)
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

	t.Fatalf("run verify: %v\nstdout=%s\nstderr=%s", err, stdout, stderr)
	return 0, stdout, stderr
}

func findRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	d := wd
	for {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d, nil
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	return "", errors.New("could not find repo root (go.mod)")
}

func copyDir(t *testing.T, srcDir, dstDir string) {
	t.Helper()

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		t.Fatalf("copy dir: %v", err)
	}
}

func readVerifyJSON(t *testing.T, stdout string) verifyJSON {
	t.Helper()
	var v verifyJSON
	dec := json.NewDecoder(strings.NewReader(stdout))
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode verify json: %v\nstdout=%s", err, stdout)
	}
	return v
}

func TestVerifyExitCodeOK(t *testing.T) {
	bin := buildBinaryOnce(t)
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code, out, _ := runVerifyCLI(t, bin, "--bundle", bundle, "--json")
	if code != 0 {
		t.Fatalf("exit code=%d, want 0\nstdout=%s", code, out)
	}
	_ = readVerifyJSON(t, out)
}

func TestVerifyExitCodeMismatch(t *testing.T) {
	bin := buildBinaryOnce(t)
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	srcBundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	tmp := t.TempDir()
	dstBundle := filepath.Join(tmp, "snapshots", "demo")
	copyDir(t, srcBundle, dstBundle)

	snapshotPath := filepath.Join(dstBundle, "snapshot.json")
	b, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatal(err)
	}
	b = stripUTF8BOM(b)

	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	obj["expected_hash_v1"] = strings.Repeat("0", 64)

	nb, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.WriteFile(snapshotPath, nb, 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}

	code, out, _ := runVerifyCLI(t, bin, "--bundle", dstBundle, "--json")
	if code != 2 {
		t.Fatalf("exit code=%d, want 2\nstdout=%s", code, out)
	}
	_ = readVerifyJSON(t, out)
}

func TestVerifyExitCodeWriteBlockedExistingExpected(t *testing.T) {
	bin := buildBinaryOnce(t)
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code, out, _ := runVerifyCLI(t, bin, "--bundle", bundle, "--write-expected", "--json")
	if code != 3 {
		t.Fatalf("exit code=%d, want 3\nstdout=%s", code, out)
	}
	v := readVerifyJSON(t, out)
	if !v.WriteBlocked {
		t.Fatalf("write_blocked=%v, want true\nstdout=%s", v.WriteBlocked, out)
	}
	if v.WriteReason != "existing_expected_present" {
		t.Fatalf("write_reason=%q, want %q\nstdout=%s", v.WriteReason, "existing_expected_present", out)
	}
}

func TestVerifyExitCodeSnapshotNotFound(t *testing.T) {
	bin := buildBinaryOnce(t)

	tmp := t.TempDir()
	bundle := filepath.Join(tmp, "snapshots", "does-not-exist")

	code, out, _ := runVerifyCLI(t, bin, "--bundle", bundle, "--json")
	if code != 4 {
		t.Fatalf("exit code=%d, want 4\nstdout=%s", code, out)
	}
	// The CLI still emits JSON for this case.
	_ = readVerifyJSON(t, out)
}

func stripUTF8BOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}
