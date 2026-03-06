package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunIdentityShowWithIO_Success(t *testing.T) {
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityShowWithIO(nil, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"name": "local"`) {
		t.Fatalf("unexpected stdout: %s", stdout.String())
	}
}

func TestRunIdentityExportWithIO_Success(t *testing.T) {
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityExportWithIO([]string{outDir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}

	for _, name := range []string{"identity.json", "identity.pub"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("missing exported file %s: %v", name, err)
		}
	}
}

func TestRunIdentityImportWithIO_Success(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	srcDir := filepath.Join(tmp, "source")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "identity.json"), []byte(`{"name":"local","algorithm":"ed25519"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "identity.pub"), []byte("abcd"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runIdentityImportWithIO([]string{srcDir}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected 0, got %d, stderr=%s", code, stderr.String())
	}

	for _, name := range []string{"trusted_identity.json", "trusted_identity.pub"} {
		if _, err := os.Stat(filepath.Join(".digiemu", name)); err != nil {
			t.Fatalf("missing imported file %s: %v", name, err)
		}
	}
}
