package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRunVerifyBundle_Success(t *testing.T) {
	dir := t.TempDir()

	canonical := `{"hello":"world"}`
	sum := sha256.Sum256([]byte(canonical))
	expected := hex.EncodeToString(sum[:])

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

func TestRunVerifyBundle_MissingReferencedFile(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{}`), 0o644); err != nil {
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

func TestRunVerifyBundle_ReferencedPathIsDirectory(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "snapshot.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "trace.json"), 0o755); err != nil {
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

func TestRunVerifyBundle_Usage(t *testing.T) {
	code := runVerifyBundle(nil)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
}

func TestRunVerifyBundle_InvalidSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(path, []byte(`{"version":"","ref":""}`), 0o644); err != nil {
		t.Fatal(err)
	}

	code := runVerifyBundle([]string{path})
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestRunVerifyBundle_MissingFileReferences(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(path, []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123","snapshot":"","meta":"","trace":""}`), 0o644); err != nil {
		t.Fatal(err)
	}

	code := runVerifyBundle([]string{path})
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestRunVerifyBundle_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(path, []byte(`{invalid`), 0o644); err != nil {
		t.Fatal(err)
	}

	code := runVerifyBundle([]string{path})
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}
