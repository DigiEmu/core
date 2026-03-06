package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type snapshotFileDoc struct {
	Canonical string `json:"canonical"`
	SHA256    string `json:"sha256"`
}

type snapshotMetaDoc struct {
	Source    string `json:"source"`
	InputPath string `json:"input_path"`
	Size      int    `json:"size"`
}

type snapshotTraceDoc struct {
	Source    string `json:"source"`
	InputPath string `json:"input_path"`
	SHA256    string `json:"sha256"`
	CreatedAt string `json:"created_at"`
	Mode      string `json:"mode"`
}

type snapshotBundleDoc struct {
	Version  string `json:"version"`
	Ref      string `json:"ref"`
	Snapshot string `json:"snapshot"`
	Meta     string `json:"meta"`
	Trace    string `json:"trace"`
}

func runSnapshot(args []string) {
	os.Exit(runSnapshotWithIO(args, os.Stdout, os.Stderr))
}

func runSnapshotWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "file" {
		return runSnapshotFileWithIO(args[1:], stdout, stderr)
	}

	fmt.Fprintln(stderr, "usage: digiemu snapshot file <path>")
	return 2
}

func runSnapshotFileWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu snapshot file <path>")
		return 2
	}

	path := args[0]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "snapshot file: %v\n", err)
		return 1
	}

	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	ref := hash[:12]

	outDir := filepath.Join("snapshots", ref)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "snapshot file: create output dir: %v\n", err)
		return 1
	}

	now := time.Now().UTC().Format(time.RFC3339)

	snapshotDoc := snapshotFileDoc{
		Canonical: string(data),
		SHA256:    hash,
	}
	metaDoc := snapshotMetaDoc{
		Source:    "file",
		InputPath: path,
		Size:      len(data),
	}
	traceDoc := snapshotTraceDoc{
		Source:    "file",
		InputPath: path,
		SHA256:    hash,
		CreatedAt: now,
		Mode:      "snapshot-file-v0.1",
	}
	bundleDoc := snapshotBundleDoc{
		Version:  "snapshot-bundle-v0.1",
		Ref:      ref,
		Snapshot: "snapshot.json",
		Meta:     "meta.json",
		Trace:    "trace.json",
	}

	snapshotBytes, err := json.MarshalIndent(snapshotDoc, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "snapshot file: marshal snapshot: %v\n", err)
		return 1
	}
	metaBytes, err := json.MarshalIndent(metaDoc, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "snapshot file: marshal meta: %v\n", err)
		return 1
	}
	traceBytes, err := json.MarshalIndent(traceDoc, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "snapshot file: marshal trace: %v\n", err)
		return 1
	}
	bundleBytes, err := json.MarshalIndent(bundleDoc, "", "  ")
	if err != nil {
		fmt.Fprintf(stderr, "snapshot file: marshal bundle: %v\n", err)
		return 1
	}

	if err := os.WriteFile(filepath.Join(outDir, "snapshot.json"), snapshotBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "snapshot file: write snapshot.json: %v\n", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outDir, "meta.json"), metaBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "snapshot file: write meta.json: %v\n", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outDir, "trace.json"), traceBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "snapshot file: write trace.json: %v\n", err)
		return 1
	}
	if err := os.WriteFile(filepath.Join(outDir, "bundle.json"), bundleBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "snapshot file: write bundle.json: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "OK snapshot %s\n", outDir)
	return 0
}
