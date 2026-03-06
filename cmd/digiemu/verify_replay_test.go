package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	pkgreplay "digiemu-core/pkg/replay"
)

func TestRunVerifyReplayWithIO_Success_NoExpected(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), []byte(`{"hello":"world"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "trace.json"), []byte(`{"sha256":"abc"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bundle.json"), []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := runVerifyReplayWithIO([]string{filepath.Join(dir, "bundle.json")}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
}

func TestRunVerifyReplayWithIO_Success_WithExpected(t *testing.T) {
	dir := t.TempDir()

	snapshotPath := filepath.Join(dir, "snapshot.json")
	if err := os.WriteFile(snapshotPath, []byte(`{"hello":"world"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	snapshotBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatal(err)
	}
	// verify replay uses FromBytes(snapshot file bytes) to avoid source-dependent hashes.
	res, err := pkgreplay.FromBytes(snapshotBytes)
	if err != nil {
		t.Fatal(err)
	}
	out, err := res.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(out)
	expected := hex.EncodeToString(sum[:])

	traceBytes, err := json.Marshal(map[string]any{"replay_sha256": expected})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "trace.json"), traceBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bundle.json"), []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := runVerifyReplayWithIO([]string{filepath.Join(dir, "bundle.json")}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
}

func TestRunVerifyReplayWithIO_Mismatch(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), []byte(`{"hello":"world"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "trace.json"), []byte(`{"replay_sha256":"deadbeef"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bundle.json"), []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := runVerifyReplayWithIO([]string{filepath.Join(dir, "bundle.json")}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestRunVerifyReplayWithIO_Usage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runVerifyReplayWithIO(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
}
