package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"digiemu-core/internal/verify"
)

func runReplay(args []string) {
	fs := flag.NewFlagSet("replay", flag.ExitOnError)
	bundlePath := fs.String("bundle", "", "Path to snapshots/<ref> directory containing snapshot.json")
	jsonOut := fs.Bool("json", false, "Emit pretty JSON output (machine readable)")
	_ = fs.Parse(args)

	root := filepath.Clean(strings.TrimSpace(*bundlePath))
	if root == "." || strings.TrimSpace(*bundlePath) == "" {
		fmt.Fprintln(os.Stderr, "--bundle is required")
		fs.Usage()
		os.Exit(4)
	}

	parent := filepath.Base(filepath.Dir(root))
	if parent != "snapshots" {
		fmt.Fprintln(os.Stderr, "--bundle must point to a snapshots/<ref> directory containing snapshot.json")
		os.Exit(4)
	}

	b, trace, err := verify.LoadBundleRootV1(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "replay load error: %v\n", err)
		os.Exit(4)
	}

	state, err := verify.ReplayV1(b, trace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "replay error: %v\n", err)
		os.Exit(4)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(state); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
			os.Exit(5)
		}
		os.Exit(0)
	}

	// human summary
	ref := filepath.Base(root)
	fmt.Fprintf(os.Stdout, "OK replay %s units=%d versions=%d claims=%d meaning=%d uncertainty=%d\n",
		ref, len(state.Units), len(state.Versions), len(state.Claims), len(state.Meaning), len(state.Uncertainty))
	os.Exit(0)
}
