package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureLocalIdentity_CreatesFiles(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }()

	meta, pub, priv, err := ensureLocalIdentity()
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name == "" || meta.Algorithm != "ed25519" {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	if len(pub) == 0 || len(priv) == 0 {
		t.Fatal("expected keys")
	}

	checks := []string{
		filepath.Join(".digiemu", "identity.key"),
		filepath.Join(".digiemu", "identity.pub"),
		filepath.Join(".digiemu", "identity.json"),
	}
	for _, p := range checks {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing identity file %s: %v", p, err)
		}
	}
}

func TestRunSignBundleWithIO_AndVerifySignatureWithIO_TrustedIdentity(t *testing.T) {
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
	if !strings.Contains(signOut.String(), "identity=local") {
		t.Fatalf("unexpected sign stdout: %s", signOut.String())
	}

	var verifyOut bytes.Buffer
	var verifyErr bytes.Buffer

	code = runVerifySignatureWithIO([]string{bundlePath}, &verifyOut, &verifyErr)
	if code != 0 {
		t.Fatalf("expected verify exit 0, got %d, stderr=%s", code, verifyErr.String())
	}
	if !strings.Contains(verifyOut.String(), "identity=local") {
		t.Fatalf("unexpected verify stdout: %s", verifyOut.String())
	}
}
