package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"digiemu-core/internal/verify"
	derr "digiemu-core/pkg/digiemu"
	"digiemu-core/pkg/snapshot"
)

func runVerify(args []string) {
	os.Exit(runVerifyWithIO(args, os.Stdout, os.Stderr))
}

func runVerifyWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "bundle" {
		return runVerifyBundle(args[1:])
	}
	if len(args) > 0 && args[0] == "signature" {
		return runVerifySignatureWithIO(args[1:], stdout, stderr)
	}
	if len(args) > 0 && args[0] == "replay" {
		return runVerifyReplayWithIO(args[1:], stdout, stderr)
	}

	args = normalizeJSONFlagArgs(args)

	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(stderr)
	refStr := fs.String("ref", "", "Snapshot hash reference to verify (required)")
	dataDir := fs.String("data", "./data", "Data directory containing snapshots")
	fixtureRoot := fs.String("fixture-root", "data/test-fixtures", "Fixture root directory (contains snapshots/<ref>/snapshot.json)")
	preferData := fs.Bool("prefer-data", false, "If true, prefer --data over --fixture-root when both bundles exist")
	jsonOut := fs.String("json", "", "Emit JSON output to stdout (pretty or canonical). Use --json or --json=canonical")
	strict := fs.Bool("strict", false, "Exit non-zero if verification fails")
	writeExpected := fs.Bool("write-expected", false, "Write expected_hash_v1 only if snapshot has a placeholder; never overwrites.")
	bundlePath := fs.String("bundle", "", "Optional bundle root path to load directly (bypass ref lookup)")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	mode, err := parseJSONMode(*jsonOut)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		fs.Usage()
		return 2
	}

	// Allow --ref to be omitted when --bundle is provided. Derive ref from bundle root name.
	if strings.TrimSpace(*bundlePath) != "" && strings.TrimSpace(*refStr) == "" {
		root := filepath.Clean(strings.TrimSpace(*bundlePath))
		derivedRef := filepath.Base(root)
		*refStr = derivedRef
	}

	if strings.TrimSpace(*refStr) == "" {
		fmt.Fprintln(stderr, "--ref is required")
		fs.Usage()
		return 2
	}

	ref := snapshot.Ref{Hash: snapshot.Hash(strings.TrimSpace(*refStr))}
	if err := ref.Validate(); err != nil {
		fmt.Fprintf(stderr, "invalid ref: %v\n", err)
		return 2
	}

	op := &verify.Operator{
		DataDir:     *dataDir,
		FixtureRoot: *fixtureRoot,
		PreferData:  *preferData,
		Options: verify.Options{
			WriteExpected: *writeExpected,
		},
	}

	// If --bundle is provided, treat it as snapshots/<ref> root and derive roots so internal lookup lands there.
	if strings.TrimSpace(*bundlePath) != "" {
		root := filepath.Clean(strings.TrimSpace(*bundlePath))
		parent := filepath.Base(filepath.Dir(root))
		if parent != "snapshots" {
			fmt.Fprintln(stderr, "--bundle must point to a snapshots/<ref> directory containing snapshot.json")
			return 2
		}
		derivedRoot := filepath.Dir(filepath.Dir(root))
		op.FixtureRoot = derivedRoot
		op.DataDir = derivedRoot
		op.PreferData = false
	}

	res, opErr := op.Verify(ref)

	if mode != jsonModeNone {
		// JSON mode contract: stdout is JSON only.
		var encErr error
		switch mode {
		case jsonModePretty:
			encErr = writePrettyJSON(stdout, res)
		case jsonModeCanonical:
			encErr = writeCanonicalJSON(stdout, res)
		default:
			encErr = fmt.Errorf("unexpected json mode: %q", mode)
		}
		if encErr != nil {
			fmt.Fprintf(stderr, "encode error: %v\n", encErr)
			return 4
		}
		// In JSON mode, any operational error text goes to stderr only.
		if opErr != nil {
			fmt.Fprintf(stderr, "verify error: %v\n", opErr)
		}
	} else {
		writeSuffix := ""
		if res.WroteExpected {
			writeSuffix = " (expected written)"
		} else if res.WriteBlocked {
			switch res.WriteReason {
			case verify.WriteReasonExistingExpected:
				writeSuffix = " (write blocked: existing expected present)"
			case verify.WriteReasonSnapshotNotFound:
				writeSuffix = " (write blocked: snapshot not found)"
			case verify.WriteReasonSnapshotInvalidJSON:
				writeSuffix = " (write blocked: snapshot invalid json)"
			case verify.WriteReasonInvalidHash:
				writeSuffix = " (write blocked: invalid hash)"
			default:
				writeSuffix = " (write blocked)"
			}
		}

		if res.OK {
			fmt.Fprintf(stdout, "OK %s%s\n", res.Ref, writeSuffix)
		} else {
			msg := res.Message
			fmt.Fprintf(stdout, "FAIL %s %s%s\n", res.Ref, msg, writeSuffix)
		}
	}

	_ = strict // preserved for backward compatibility; exit codes follow Phase A2 contract.
	_ = writeExpected

	if opErr != nil {
		return exitCodeForError(opErr)
	}
	if res.OK {
		return 0
	}
	// Verification failed without an operational error; classify as verify failure.
	return exitCodeForError(derr.VerifyFailed("cli.verify", res.Ref, nil))
}
