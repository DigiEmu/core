package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeBundleFixture(t *testing.T, dir string) string {
	t.Helper()

	bundleDir := filepath.Join(dir, "snapshots", "abc123")
	if err := os.MkdirAll(bundleDir, 0o755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"bundle.json":   `{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`,
		"snapshot.json": `{"canonical":"{\"hello\":\"world\"}","sha256":"x"}`,
		"meta.json":     `{"source":"file"}`,
		"trace.json":    `{"sha256":"x","replay_sha256":"y"}`,
	}

	for name, body := range files {
		if err := os.WriteFile(filepath.Join(bundleDir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	return filepath.Join(bundleDir, "bundle.json")
}

func TestRunExportBundleWithIO_Success(t *testing.T) {
	root := t.TempDir()
	bundlePath := writeBundleFixture(t, root)
	outDir := filepath.Join(root, "out")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runExportBundleWithIO([]string{bundlePath, outDir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "OK export bundle") {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}

	checks := []string{"bundle.json", "snapshot.json", "meta.json", "trace.json"}
	for _, name := range checks {
		p := filepath.Join(outDir, "abc123", name)
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing exported file %s: %v", name, err)
		}
	}
}

func TestRunImportBundleWithIO_Success(t *testing.T) {
	root := t.TempDir()
	bundlePath := writeBundleFixture(t, root)
	dir := filepath.Dir(bundlePath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runImportBundleWithIO([]string{dir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "OK import bundle") {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}
}

func TestRunImportBundleWithIO_MissingFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bundle.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runImportBundleWithIO([]string{dir}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "missing") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}
