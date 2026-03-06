package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type verifyBundleDoc struct {
	Version  string `json:"version"`
	Ref      string `json:"ref"`
	Snapshot string `json:"snapshot"`
	Meta     string `json:"meta"`
	Trace    string `json:"trace"`
}

type verifyTraceDoc struct {
	Source    string `json:"source"`
	InputPath string `json:"input_path"`
	SHA256    string `json:"sha256"`
	CreatedAt string `json:"created_at"`
	Mode      string `json:"mode"`
}

type verifySnapshotDoc struct {
	Canonical string `json:"canonical"`
	SHA256    string `json:"sha256"`
}

func runVerifyBundle(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: digiemu verify bundle <path>")
		return 2
	}

	path := args[0]

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify bundle: %v\n", err)
		return 1
	}

	var bundle verifyBundleDoc
	if err := json.Unmarshal(data, &bundle); err != nil {
		fmt.Fprintln(os.Stderr, "verify bundle: invalid json")
		return 1
	}

	if bundle.Version == "" || bundle.Ref == "" {
		fmt.Fprintln(os.Stderr, "verify bundle: invalid schema")
		return 1
	}
	if bundle.Snapshot == "" || bundle.Meta == "" || bundle.Trace == "" {
		fmt.Fprintln(os.Stderr, "verify bundle: missing file references")
		return 1
	}

	baseDir := filepath.Dir(path)
	snapshotPath := filepath.Join(baseDir, bundle.Snapshot)
	metaPath := filepath.Join(baseDir, bundle.Meta)
	tracePath := filepath.Join(baseDir, bundle.Trace)
	required := []string{
		snapshotPath,
		metaPath,
		tracePath,
	}
	for _, requiredPath := range required {
		info, err := os.Stat(requiredPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "verify bundle: missing file %s\n", filepath.Base(requiredPath))
			return 1
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "verify bundle: expected file, got directory %s\n", filepath.Base(requiredPath))
			return 1
		}
	}

	traceBytes, err := os.ReadFile(tracePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify bundle: read trace: %v\n", err)
		return 1
	}
	var trace verifyTraceDoc
	if err := json.Unmarshal(traceBytes, &trace); err != nil {
		fmt.Fprintln(os.Stderr, "verify bundle: invalid trace json")
		return 1
	}
	if trace.SHA256 == "" {
		fmt.Fprintln(os.Stderr, "verify bundle: missing trace sha256")
		return 1
	}

	snapshotBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify bundle: read snapshot: %v\n", err)
		return 1
	}
	var snapshot verifySnapshotDoc
	if err := json.Unmarshal(snapshotBytes, &snapshot); err != nil {
		fmt.Fprintln(os.Stderr, "verify bundle: invalid snapshot json")
		return 1
	}
	if snapshot.SHA256 == "" {
		fmt.Fprintln(os.Stderr, "verify bundle: missing snapshot sha256")
		return 1
	}

	// Integrity: trace sha256 must match what snapshot declares, and must be the sha256 of snapshot.canonical.
	sum := sha256.Sum256([]byte(snapshot.Canonical))
	actual := hex.EncodeToString(sum[:])
	if actual != trace.SHA256 {
		fmt.Fprintf(os.Stderr, "verify bundle: sha256 mismatch expected=%s actual=%s\n", trace.SHA256, actual)
		return 1
	}
	if snapshot.SHA256 != trace.SHA256 {
		fmt.Fprintf(os.Stderr, "verify bundle: sha256 mismatch expected=%s actual=%s\n", trace.SHA256, snapshot.SHA256)
		return 1
	}

	fmt.Printf("OK bundle %s ref=%s\n", bundle.Version, bundle.Ref)

	return 0
}
