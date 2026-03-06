package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runIdentity(args []string) {
	os.Exit(runIdentityWithIO(args, os.Stdout, os.Stderr))
}

func runIdentityWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "identity subcommands: show | export | import | fingerprint")
		return 2
	}

	switch args[0] {
	case "show":
		return runIdentityShowWithIO(args[1:], stdout, stderr)
	case "export":
		return runIdentityExportWithIO(args[1:], stdout, stderr)
	case "import":
		return runIdentityImportWithIO(args[1:], stdout, stderr)
	case "fingerprint":
		return runIdentityFingerprintWithIO(args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "identity subcommands: show | export | import | fingerprint")
		return 2
	}
}

func runIdentityShowWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: digiemu identity show")
		return 2
	}

	meta, _, err := loadLocalIdentityPublic()
	if err != nil {
		fmt.Fprintf(stderr, "identity show: %v\n", err)
		return 1
	}

	out, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "identity show: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, string(out))
	return 0
}

func runIdentityExportWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu identity export <out-dir>")
		return 2
	}

	outDir := args[0]

	metaPath := filepath.Join(".digiemu", "identity.json")
	pubPath := filepath.Join(".digiemu", "identity.pub")

	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		fmt.Fprintf(stderr, "identity export: %v\n", err)
		return 1
	}
	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		fmt.Fprintf(stderr, "identity export: %v\n", err)
		return 1
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "identity export: %v\n", err)
		return 1
	}

	if err := os.WriteFile(filepath.Join(outDir, "identity.json"), metaBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "identity export: %v\n", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outDir, "identity.pub"), pubBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "identity export: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "OK identity export %s\n", outDir)
	return 0
}

func runIdentityImportWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu identity import <dir>")
		return 2
	}

	srcDir := args[0]

	metaBytes, err := os.ReadFile(filepath.Join(srcDir, "identity.json"))
	if err != nil {
		fmt.Fprintf(stderr, "identity import: %v\n", err)
		return 1
	}
	pubBytes, err := os.ReadFile(filepath.Join(srcDir, "identity.pub"))
	if err != nil {
		fmt.Fprintf(stderr, "identity import: %v\n", err)
		return 1
	}

	if err := os.MkdirAll(".digiemu", 0o755); err != nil {
		fmt.Fprintf(stderr, "identity import: %v\n", err)
		return 1
	}

	if err := os.WriteFile(filepath.Join(".digiemu", "trusted_identity.json"), metaBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "identity import: %v\n", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(".digiemu", "trusted_identity.pub"), pubBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "identity import: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "OK identity import %s\n", srcDir)
	return 0
}

func runIdentityFingerprintWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: digiemu identity fingerprint")
		return 2
	}

	pubBytes, err := os.ReadFile(filepath.Join(".digiemu", "identity.pub"))
	if err != nil {
		fmt.Fprintf(stderr, "identity fingerprint: %v\n", err)
		return 1
	}

	sum := sha256.Sum256(pubBytes)
	fmt.Fprintln(stdout, hex.EncodeToString(sum[:]))
	return 0
}
