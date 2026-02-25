package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bun "digiemu-core/internal/bundle"
	"digiemu-core/internal/verify"
	"digiemu-core/pkg/snapshot"
)

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	refStr := fs.String("ref", "", "Snapshot hash reference to verify (required)")
	dataDir := fs.String("data", "./data", "Data directory containing snapshots")
	fixtureRoot := fs.String("fixture-root", "data/test-fixtures", "Fixture root directory (contains snapshots/<ref>/snapshot.json)")
	preferData := fs.Bool("prefer-data", false, "If true, prefer --data over --fixture-root when both bundles exist")
	jsonOut := fs.Bool("json", false, "Emit stable JSON output (machine readable)")
	strict := fs.Bool("strict", false, "Exit non-zero if verification fails")
	writeExpected := fs.Bool("write-expected", false, "Write expected_hash_v1 only if snapshot has a placeholder; never overwrites.")
	bundlePath := fs.String("bundle", "", "Optional bundle root path to load directly (bypass ref lookup)")
	_ = fs.Parse(args)

	// Allow --ref to be omitted when --bundle is provided. Derive ref from bundle root name.
	if strings.TrimSpace(*bundlePath) != "" && strings.TrimSpace(*refStr) == "" {
		root := filepath.Clean(strings.TrimSpace(*bundlePath))
		derivedRef := filepath.Base(root)
		*refStr = derivedRef
	}

	if strings.TrimSpace(*refStr) == "" {
		fmt.Fprintln(os.Stderr, "--ref is required")
		fs.Usage()
		os.Exit(4)
	}

	ref := snapshot.Ref{Hash: snapshot.Hash(strings.TrimSpace(*refStr))}
	if err := ref.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid ref: %v\n", err)
		os.Exit(4)
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
			fmt.Fprintln(os.Stderr, "--bundle must point to a snapshots/<ref> directory containing snapshot.json")
			os.Exit(4)
		}
		derivedRoot := filepath.Dir(filepath.Dir(root))
		op.FixtureRoot = derivedRoot
		op.DataDir = derivedRoot
		op.PreferData = false
	}

	res, err := op.Verify(ref)
	if err != nil {
		// Only treat truly unexpected errors as internal. Typed write-policy errors
		// should still produce normal output and deterministic exit codes.
		if !(errors.Is(err, bun.ErrExpectedAlreadySet) ||
			errors.Is(err, bun.ErrSnapshotNotFound) ||
			errors.Is(err, bun.ErrSnapshotInvalidJSON) ||
			errors.Is(err, bun.ErrInvalidNewHash)) {
			fmt.Fprintf(os.Stderr, "verify error: %v\n", err)
			os.Exit(5)
		}
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if e := enc.Encode(res); e != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", e)
			os.Exit(5)
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
			fmt.Fprintf(os.Stdout, "OK %s%s\n", res.Ref, writeSuffix)
		} else {
			msg := res.Message
			fmt.Fprintf(os.Stdout, "FAIL %s %s%s\n", res.Ref, msg, writeSuffix)
		}
	}

	_ = strict // preserved for backward compatibility; exit codes are deterministic now.
	os.Exit(verifyExitCode(res, *writeExpected, err))
}

func verifyExitCode(res verify.ResultV1, writeExpected bool, opErr error) int {
	// write blocked (ErrExpectedAlreadySet) => 3 (only when --write-expected set)
	if writeExpected && opErr != nil && errors.Is(opErr, bun.ErrExpectedAlreadySet) {
		return 3
	}
	// Governance: also return 3 when the operator recorded write_blocked for existing expected,
	// even if verification itself succeeded.
	if writeExpected && res.WriteBlocked && res.WriteReason == verify.WriteReasonExistingExpected {
		return 3
	}

	if res.OK {
		return 0
	}

	// mismatch (hash computed but does not match expected)
	if res.Got != "" && res.Expected != "" && res.Got != res.Expected {
		return 2
	}

	// invalid snapshot/json/file => 4
	if opErr != nil {
		if errors.Is(opErr, bun.ErrSnapshotNotFound) || errors.Is(opErr, bun.ErrSnapshotInvalidJSON) || errors.Is(opErr, bun.ErrInvalidNewHash) {
			return 4
		}
	}
	if res.WriteBlocked {
		return 4
	}
	if len(res.Errors) > 0 {
		return 4
	}

	// internal error
	return 5
}
