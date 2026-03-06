package main

import (
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
	required := []string{
		filepath.Join(baseDir, bundle.Snapshot),
		filepath.Join(baseDir, bundle.Meta),
		filepath.Join(baseDir, bundle.Trace),
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

	fmt.Printf("OK bundle %s ref=%s\n", bundle.Version, bundle.Ref)

	return 0
}
