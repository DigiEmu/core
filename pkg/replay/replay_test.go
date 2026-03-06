package replay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func expectedSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func TestFromBytes_ReturnsResultWithMetadata(t *testing.T) {
	in := []byte(`{"hello":"world"}`)

	out, err := FromBytes(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(in, out.Canonical) {
		t.Fatalf("expected identical bytes, got %q want %q", out.Canonical, in)
	}

	if out.Size != len(in) {
		t.Fatalf("expected size %d, got %d", len(in), out.Size)
	}

	if out.Source != "bytes" {
		t.Fatalf("expected source %q, got %q", "bytes", out.Source)
	}

	if out.SHA256 != expectedSHA256(in) {
		t.Fatalf("expected sha256 %q, got %q", expectedSHA256(in), out.SHA256)
	}

	if len(out.Canonical) > 0 && &out.Canonical[0] == &in[0] {
		t.Fatalf("expected copy, got same backing array")
	}
}

func TestResult_Marshal(t *testing.T) {
	in := []byte(`{"hello":"world"}`)

	out, err := FromBytes(in)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := out.Marshal()
	if err != nil {
		t.Fatalf("expected no marshal error, got %v", err)
	}

	want := `{"canonical":"{\"hello\":\"world\"}","size":17,"source":"bytes","sha256":"` + expectedSHA256(in) + `"}`

	if string(got) != want {
		t.Fatalf("unexpected marshal output:\n got: %s\nwant: %s", string(got), want)
	}
}

func TestResult_Marshal_EmptyCanonical(t *testing.T) {
	_, err := (Result{
		Size:   1,
		Source: "bytes",
		SHA256: "abc",
	}).Marshal()
	if err == nil {
		t.Fatal("expected error for empty canonical, got nil")
	}
}

func TestFromBytes_EmptyInput(t *testing.T) {
	_, err := FromBytes(nil)
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestFromReader_ReturnsResultWithMetadata(t *testing.T) {
	in := []byte(`{"kind":"bundle"}`)

	out, err := FromReader(bytes.NewReader(in))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(in, out.Canonical) {
		t.Fatalf("expected identical bytes, got %q want %q", out.Canonical, in)
	}

	if out.Size != len(in) {
		t.Fatalf("expected size %d, got %d", len(in), out.Size)
	}

	if out.Source != "reader" {
		t.Fatalf("expected source %q, got %q", "reader", out.Source)
	}

	if out.SHA256 != expectedSHA256(in) {
		t.Fatalf("expected sha256 %q, got %q", expectedSHA256(in), out.SHA256)
	}
}

func TestFromReader_NilReader(t *testing.T) {
	_, err := FromReader(nil)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
}

func TestFromReader_EmptyInput(t *testing.T) {
	_, err := FromReader(bytes.NewReader(nil))
	if err == nil {
		t.Fatal("expected error for empty reader input, got nil")
	}
}

func TestFromFile_ReturnsResultWithMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.json")
	in := []byte(`{"snapshot":"abc123"}`)

	if err := os.WriteFile(path, in, 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	out, err := FromFile(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bytes.Equal(in, out.Canonical) {
		t.Fatalf("expected identical bytes, got %q want %q", out.Canonical, in)
	}

	if out.Size != len(in) {
		t.Fatalf("expected size %d, got %d", len(in), out.Size)
	}

	if out.Source != "file" {
		t.Fatalf("expected source %q, got %q", "file", out.Source)
	}

	if out.SHA256 != expectedSHA256(in) {
		t.Fatalf("expected sha256 %q, got %q", expectedSHA256(in), out.SHA256)
	}
}

func TestFromFile_EmptyPath(t *testing.T) {
	_, err := FromFile("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestFromFile_MissingFile(t *testing.T) {
	_, err := FromFile(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
