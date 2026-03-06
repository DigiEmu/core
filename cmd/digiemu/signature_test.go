package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSignBundleWithIO_AndVerifySignatureWithIO_Success(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	bundlePath := filepath.Join(tmp, "bundle.json")

	if err := os.WriteFile(bundlePath, []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var signOut bytes.Buffer
	var signErr bytes.Buffer

	code := runSignBundleWithIO([]string{bundlePath}, &signOut, &signErr)
	if code != 0 {
		t.Fatalf("expected sign exit 0, got %d, stderr=%s", code, signErr.String())
	}

	var verifyOut bytes.Buffer
	var verifyErr bytes.Buffer

	code = runVerifySignatureWithIO([]string{bundlePath}, &verifyOut, &verifyErr)
	if code != 0 {
		t.Fatalf("expected verify exit 0, got %d, stderr=%s", code, verifyErr.String())
	}
	if !strings.Contains(verifyOut.String(), "OK verify signature") {
		t.Fatalf("unexpected stdout: %s", verifyOut.String())
	}
	if !strings.Contains(verifyOut.String(), "identity=local") {
		t.Fatalf("unexpected stdout: %s", verifyOut.String())
	}
}

func TestRunVerifySignatureWithIO_MissingSignatureFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	bundlePath := filepath.Join(tmp, "bundle.json")

	if err := os.WriteFile(bundlePath, []byte(`{"version":"snapshot-bundle-v0.1","ref":"abc123"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runVerifySignatureWithIO([]string{bundlePath}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected 1, got %d", code)
	}
}

func TestRunSignBundleWithIO_Usage(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runSignBundleWithIO(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected 2, got %d", code)
	}
}
