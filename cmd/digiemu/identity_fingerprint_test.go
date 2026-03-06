package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunIdentityFingerprintWithIO_Success(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	_, _, _, err = ensureLocalIdentity()
	if err != nil {
		t.Fatal(err)
	}

	pubBytes, err := os.ReadFile(filepath.Join(".digiemu", "identity.pub"))
	if err != nil {
		t.Fatal(err)
	}
	expected := sha256.Sum256(pubBytes)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityFingerprintWithIO(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}

	got := strings.TrimSpace(stdout.String())
	want := hex.EncodeToString(expected[:])
	if got != want {
		t.Fatalf("unexpected fingerprint: got %s want %s", got, want)
	}
}

func TestRunIdentityFingerprintWithIO_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityFingerprintWithIO([]string{"extra"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "usage: digiemu identity fingerprint") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunIdentityFingerprintWithIO_MissingIdentity(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityFingerprintWithIO(nil, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "identity fingerprint:") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunIdentityFingerprintWithIO_AfterExportImport(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	_, _, _, err = ensureLocalIdentity()
	if err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tmp, "exported")
	if code := runIdentityExportWithIO([]string{outDir}, io.Discard, io.Discard); code != 0 {
		t.Fatalf("export failed: %d", code)
	}
	if code := runIdentityImportWithIO([]string{outDir}, io.Discard, io.Discard); code != 0 {
		t.Fatalf("import failed: %d", code)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityFingerprintWithIO(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
	if len(strings.TrimSpace(stdout.String())) != 64 {
		t.Fatalf("expected 64-char hex fingerprint, got %q", stdout.String())
	}
}
