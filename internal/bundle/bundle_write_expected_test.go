package bundle

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteExpectedHashV1_Placeholder_Succeeds(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")

	// Include a BOM to ensure the reader is tolerant.
	input := "\ufeff{\n  \"expected_hash_v1\": \"REPLACE_AFTER_ASSEMBLE\",\n  \"ref\": \"demo\"\n}"
	if err := os.WriteFile(p, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	newHash := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := WriteExpectedHashV1(p, newHash); err != nil {
		t.Fatalf("WriteExpectedHashV1: %v", err)
	}

	out, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	// Output must not contain UTF-8 BOM.
	if len(out) >= 3 && out[0] == 0xEF && out[1] == 0xBB && out[2] == 0xBF {
		t.Fatalf("output contains UTF-8 BOM")
	}

	// Pretty JSON should contain indentation.
	if !strings.Contains(string(out), "\n  \"") {
		t.Fatalf("output does not look like pretty JSON: %q", string(out))
	}

	var obj map[string]any
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	got, _ := obj["expected_hash_v1"].(string)
	if got != newHash {
		t.Fatalf("expected_hash_v1 not updated: got=%q want=%q", got, newHash)
	}
}

func TestWriteExpectedHashV1_AlreadySet_ReturnsErrExpectedAlreadySet(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")

	input := "{\n  \"expected_hash_v1\": \"NOT_A_PLACEHOLDER\",\n  \"ref\": \"demo\"\n}"
	if err := os.WriteFile(p, []byte(input), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	err := WriteExpectedHashV1(p, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrExpectedAlreadySet) {
		t.Fatalf("expected ErrExpectedAlreadySet, got %v", err)
	}
}

func TestWriteExpectedHashV1_InvalidJSON_ReturnsErrSnapshotInvalidJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")

	if err := os.WriteFile(p, []byte("{not-json"), 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	err := WriteExpectedHashV1(p, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrSnapshotInvalidJSON) {
		t.Fatalf("expected ErrSnapshotInvalidJSON, got %v", err)
	}
}

func TestWriteExpectedHashV1_MissingFile_ReturnsErrSnapshotNotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "does-not-exist.json")

	err := WriteExpectedHashV1(p, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, ErrSnapshotNotFound) {
		t.Fatalf("expected ErrSnapshotNotFound, got %v", err)
	}
}
