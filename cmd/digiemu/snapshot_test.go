package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSnapshotFileWithIO_Success(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	defer func() { _ = os.Chdir(wd) }()

	in := filepath.Join(tmp, "input.json")
	if err := os.WriteFile(in, []byte(`{"hello":"world"}`), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runSnapshotWithIO([]string{"file", in}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "OK snapshot snapshots\\") && !strings.Contains(stdout.String(), "OK snapshot snapshots/") {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}

	matches, err := filepath.Glob(filepath.Join("snapshots", "*", "snapshot.json"))
	if err != nil {
		t.Fatalf("glob snapshot: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 snapshot.json, got %d", len(matches))
	}

	metaMatches, err := filepath.Glob(filepath.Join("snapshots", "*", "meta.json"))
	if err != nil {
		t.Fatalf("glob meta: %v", err)
	}
	if len(metaMatches) != 1 {
		t.Fatalf("expected 1 meta.json, got %d", len(metaMatches))
	}

	traceMatches, err := filepath.Glob(filepath.Join("snapshots", "*", "trace.json"))
	if err != nil {
		t.Fatalf("glob trace: %v", err)
	}
	if len(traceMatches) != 1 {
		t.Fatalf("expected 1 trace.json, got %d", len(traceMatches))
	}
}

func TestRunSnapshotFileWithIO_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runSnapshotWithIO([]string{"file"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "usage: digiemu snapshot file <path>") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunSnapshotFileWithIO_MissingFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runSnapshotWithIO([]string{"file", "missing.json"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "snapshot file:") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}
