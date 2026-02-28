package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) { return 0, errors.New("write failed") }

func TestRunReplay_MissingBundle_IsUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runReplayWithIO([]string{}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit=%d, want 2", exit)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr usage")
	}
}

func TestRunReplay_BadBundlePath_IsUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runReplayWithIO([]string{"--bundle", filepath.Join(os.TempDir(), "not-snapshots", "demo")}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit=%d, want 2", exit)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr")
	}
}

func TestRunReplay_JSONEncodeError_IsInternal(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "snapshots", "demo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Minimal valid JSON for snapshot.json; ReplayV1 only needs it to be JSON.
	if err := os.WriteFile(filepath.Join(root, "snapshot.json"), []byte("{\"expected_hash_v1\":\"x\",\"foo\":1}\n"), 0o644); err != nil {
		t.Fatalf("write snapshot.json: %v", err)
	}

	var stderr bytes.Buffer
	exit := runReplayWithIO([]string{"--bundle", root, "--json=canonical"}, errWriter{}, &stderr)
	if exit != 4 {
		t.Fatalf("exit=%d, want 4", exit)
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr")
	}
}
