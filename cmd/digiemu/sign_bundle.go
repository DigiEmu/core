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

type signatureDoc struct {
	Algorithm    string `json:"algorithm"`
	Identity     string `json:"identity"`
	BundleSHA256 string `json:"bundle_sha256"`
	Signature    string `json:"signature"`
}

func runSignBundleWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu sign bundle <bundle.json>")
		return 2
	}

	bundlePath := args[0]

	data, err := os.ReadFile(bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "sign bundle: %v\n", err)
		return 1
	}

	meta, pub, priv, err := ensureLocalIdentity()
	if err != nil {
		fmt.Fprintf(stderr, "sign bundle: ensure identity: %v\n", err)
		return 1
	}
	_ = pub

	sum := sha256.Sum256(data)
	hashHex := hex.EncodeToString(sum[:])

	sig := ed25519.Sign(priv, data)

	doc := signatureDoc{
		Algorithm:    meta.Algorithm,
		Identity:     meta.Name,
		BundleSHA256: hashHex,
		Signature:    hex.EncodeToString(sig),
	}

	sigBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "sign bundle: marshal signature: %v\n", err)
		return 1
	}

	dir := filepath.Dir(bundlePath)

	if err := os.WriteFile(filepath.Join(dir, "signature.json"), sigBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "sign bundle: write signature.json: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "OK sign bundle %s identity=%s\n", bundlePath, meta.Name)
	return 0
}
