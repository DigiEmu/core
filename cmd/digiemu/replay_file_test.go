package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReplayFileWithIO_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.json")

	if err := os.WriteFile(path, []byte(`{"hello":"world"}`), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayWithIO([]string{"file", path}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, `"source":"file"`) {
		t.Fatalf("expected output to contain source=file, got: %s", out)
	}
	if !strings.Contains(out, `"sha256":"`) {
		t.Fatalf("expected output to contain sha256, got: %s", out)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got: %s", stderr.String())
	}
}

func TestRunReplayFileWithIO_MissingPath(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayWithIO([]string{"file"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "usage: digiemu replay file <path>") {
		t.Fatalf("expected usage message, got: %s", stderr.String())
	}
}

func TestRunReplayFileWithIO_MissingFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayWithIO([]string{"file", filepath.Join(t.TempDir(), "missing.json")}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "replay file:") {
		t.Fatalf("expected replay file error, got: %s", stderr.String())
	}
}
