package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func hashOf(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func TestRunVerifyBundle_Success_HashIntegrity(t *testing.T) {
	dir := t.TempDir()

	canonical := `{"hello":"world"}`
	expected := hashOf([]byte(canonical))

	snapshotDoc := verifySnapshotDoc{Canonical: canonical, SHA256: expected}
	snapshotBytes, err := json.Marshal(snapshotDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), snapshotBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	traceDoc := verifyTraceDoc{
		Source:    "file",
		InputPath: "x",
		SHA256:    expected,
		CreatedAt: "2026-01-01T00:00:00Z",
		Mode:      "snapshot-file-v0.1",
	}
	traceBytes, err := json.Marshal(traceDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "trace.json"), traceBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	bundle := `{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`
	path := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(path, []byte(bundle), 0o644); err != nil {
		t.Fatal(err)
	}

	code := runVerifyBundle([]string{path})
	if code != 0 {
		t.Fatalf("expected 0, got %d", code)
	}
}

func TestRunVerifyBundle_SHA256Mismatch(t *testing.T) {
	dir := t.TempDir()

	canonical := `{"hello":"world"}`
	snapshotExpected := hashOf([]byte(canonical))

	snapshotDoc := verifySnapshotDoc{Canonical: canonical, SHA256: snapshotExpected}
	snapshotBytes, err := json.Marshal(snapshotDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), snapshotBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	traceDoc := verifyTraceDoc{
		Source:    "file",
		InputPath: "x",
		SHA256:    "deadbeef",
		CreatedAt: "2026-01-01T00:00:00Z",
		Mode:      "snapshot-file-v0.1",
	}
	traceBytes, err := json.Marshal(traceDoc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "trace.json"), traceBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	bundle := `{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"snapshot.json","meta":"meta.json","trace":"trace.json"}`
	path := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(path, []byte(bundle), 0o644); err != nil {
		t.Fatal(err)
	}

	code := runVerifyBundle([]string{path})
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}
