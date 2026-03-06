package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pkgreplay "digiemu-core/pkg/replay"
)

type replayBundleDoc struct {
	Version  string `json:"version"`
	Ref      string `json:"ref"`
	Snapshot string `json:"snapshot"`
	Meta     string `json:"meta"`
	Trace    string `json:"trace"`
}

func runReplayBundleWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu replay bundle <path>")
		return 2
	}

	bundlePath := args[0]

	data, err := os.ReadFile(bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "replay bundle: %v\n", err)
		return 1
	}

	var bundle replayBundleDoc
	if err := json.Unmarshal(data, &bundle); err != nil {
		fmt.Fprintln(stderr, "replay bundle: invalid json")
		return 1
	}

	if bundle.Version == "" || bundle.Ref == "" || bundle.Snapshot == "" {
		fmt.Fprintln(stderr, "replay bundle: invalid schema")
		return 1
	}

	baseDir := filepath.Dir(bundlePath)
	snapshotPath := filepath.Join(baseDir, bundle.Snapshot)

	result, err := pkgreplay.FromFile(snapshotPath)
	if err != nil {
		fmt.Fprintf(stderr, "replay bundle: %v\n", err)
		return 1
	}

	out, err := result.Marshal()
	if err != nil {
		fmt.Fprintf(stderr, "replay bundle: marshal result: %v\n", err)
		return 1
	}

	if _, err := fmt.Fprintln(stdout, string(out)); err != nil {
		fmt.Fprintf(stderr, "replay bundle: write output: %v\n", err)
		return 1
	}

	return 0
}
