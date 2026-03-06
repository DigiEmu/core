package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runVerifySignatureWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu verify signature <bundle.json>")
		return 2
	}

	bundlePath := args[0]
	dir := filepath.Dir(bundlePath)

	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "verify signature: %v\n", err)
		return 1
	}

	sigBytes, err := os.ReadFile(filepath.Join(dir, "signature.json"))
	if err != nil {
		fmt.Fprintf(stderr, "verify signature: %v\n", err)
		return 1
	}

	var doc signatureDoc
	if err := json.Unmarshal(sigBytes, &doc); err != nil {
		fmt.Fprintln(stderr, "verify signature: invalid signature json")
		return 1
	}

	if doc.Algorithm != "ed25519" {
		fmt.Fprintln(stderr, "verify signature: unsupported algorithm")
		return 1
	}
	if doc.Identity == "" {
		fmt.Fprintln(stderr, "verify signature: missing identity")
		return 1
	}

	identityMeta, pub, err := loadLocalIdentityPublic()
	if err != nil {
		fmt.Fprintf(stderr, "verify signature: %v\n", err)
		return 1
	}
	if identityMeta.Name != doc.Identity {
		fmt.Fprintln(stderr, "verify signature: identity mismatch")
		return 1
	}

	sum := sha256.Sum256(bundleBytes)
	actualHash := hex.EncodeToString(sum[:])
	if doc.BundleSHA256 != actualHash {
		fmt.Fprintln(stderr, "verify signature: bundle hash mismatch")
		return 1
	}

	sig, err := hex.DecodeString(doc.Signature)
	if err != nil {
		fmt.Fprintln(stderr, "verify signature: invalid signature encoding")
		return 1
	}

	if !ed25519.Verify(pub, bundleBytes, sig) {
		fmt.Fprintln(stderr, "verify signature: signature check failed")
		return 1
	}

	fmt.Fprintf(stdout, "OK verify signature %s identity=%s\n", bundlePath, doc.Identity)
	return 0
}
