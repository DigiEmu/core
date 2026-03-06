package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	err = json.Unmarshal(data, &bundle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify bundle: invalid json\n")
		return 1
	}

	if bundle.Version == "" || bundle.Ref == "" {
		fmt.Fprintf(os.Stderr, "verify bundle: invalid schema\n")
		return 1
	}

	fmt.Printf("OK bundle %s ref=%s\n", bundle.Version, bundle.Ref)

	return 0
}
