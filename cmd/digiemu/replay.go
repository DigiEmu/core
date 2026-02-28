package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"digiemu-core/internal/verify"
)

func runReplay(args []string) {
	os.Exit(runReplayWithIO(args, os.Stdout, os.Stderr))
}

func runReplayWithIO(args []string, stdout, stderr io.Writer) int {
	args = normalizeJSONFlagArgs(args)

	flagset := flag.NewFlagSet("replay", flag.ContinueOnError)
	flagset.SetOutput(stderr)
	bundlePath := flagset.String("bundle", "", "Path to snapshots/<ref> directory containing snapshot.json")
	jsonOut := flagset.String("json", "", "Emit JSON output to stdout (pretty or canonical). Use --json or --json=canonical")
	if err := flagset.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	mode, err := parseJSONMode(*jsonOut)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		flagset.Usage()
		return 2
	}

	root := filepath.Clean(strings.TrimSpace(*bundlePath))
	if root == "." || strings.TrimSpace(*bundlePath) == "" {
		fmt.Fprintln(stderr, "--bundle is required")
		flagset.Usage()
		return 2
	}

	parent := filepath.Base(filepath.Dir(root))
	if parent != "snapshots" {
		fmt.Fprintln(stderr, "--bundle must point to a snapshots/<ref> directory containing snapshot.json")
		return 2
	}

	b, trace, err := verify.LoadBundleRootV1(root)
	if err != nil {
		fmt.Fprintf(stderr, "replay load error: %v\n", err)
		// Contract: load errors are treated as IO.
		// Use a conservative check for missing/permission-related errors.
		if errorsIsAny(err, fs.ErrNotExist, os.ErrNotExist) {
			return 3
		}
		return 3
	}

	state, err := verify.ReplayV1(b, trace)
	if err != nil {
		fmt.Fprintf(stderr, "replay error: %v\n", err)
		return 4
	}

	if mode != jsonModeNone {
		// JSON mode contract: stdout is JSON only.
		var encErr error
		switch mode {
		case jsonModePretty:
			encErr = writePrettyJSON(stdout, state)
		case jsonModeCanonical:
			encErr = writeCanonicalJSON(stdout, state)
		default:
			encErr = fmt.Errorf("unexpected json mode: %q", mode)
		}
		if encErr != nil {
			fmt.Fprintf(stderr, "encode error: %v\n", encErr)
			return 4
		}
		return 0
	}

	// human summary
	ref := filepath.Base(root)
	fmt.Fprintf(stdout, "OK replay %s units=%d versions=%d claims=%d meaning=%d uncertainty=%d\n",
		ref, len(state.Units), len(state.Versions), len(state.Claims), len(state.Meaning), len(state.Uncertainty))
	return 0
}

func errorsIsAny(err error, targets ...error) bool {
	for _, t := range targets {
		if t != nil && errors.Is(err, t) {
			return true
		}
	}
	return false
}
