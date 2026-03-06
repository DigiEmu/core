package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReplayBundleWithIO_Success(t *testing.T) {
	dir := t.TempDir()

	snapshot := `{"canonical":"{\"hello\":\"world\"}","sha256":"abc"}`
	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), []byte(snapshot), 0o644); err != nil {
		t.Fatal(err)
	}

	bundle := `{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, []byte(bundle), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayBundleWithIO([]string{bundlePath}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"source":"file"`) {
		t.Fatalf("expected replay output json, got: %s", stdout.String())
	}
}

func TestRunReplayBundleWithIO_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayBundleWithIO(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "usage: digiemu replay bundle <path>") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunReplayBundleWithIO_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, []byte(`{invalid`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayBundleWithIO([]string{bundlePath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestRunReplayBundleWithIO_MissingSnapshot(t *testing.T) {
	dir := t.TempDir()

	bundle := `{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`
	bundlePath := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(bundlePath, []byte(bundle), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runReplayBundleWithIO([]string{bundlePath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}
