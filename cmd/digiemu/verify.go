package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

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
	_ = fs.Parse(args)

	if strings.TrimSpace(*refStr) == "" {
		fmt.Fprintln(os.Stderr, "--ref is required")
		fs.Usage()
		os.Exit(2)
	}

	ref := snapshot.Ref{Hash: snapshot.Hash(strings.TrimSpace(*refStr))}
	bundlePath := fs.String("bundle", "", "Optional bundle root path to load directly (bypass ref lookup)")
	v := &verify.VerifierV1{DataDir: *dataDir, FixtureRoot: *fixtureRoot, PreferData: *preferData}

	// If --bundle provided, use it as root (bypass lookup)
	if strings.TrimSpace(*bundlePath) != "" {
		// Treat bundle path as root and assemble state then verify against snapshot.json inside
		// We will reuse VerifierV1 by temporarily setting FixtureRoot to the bundle path and PreferData=true
		v.FixtureRoot = *bundlePath
		v.PreferData = true
	}

	res, err := v.Verify(ref)
	if err != nil {
		// Unexpected runtime error
		fmt.Fprintf(os.Stderr, "verify error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if e := enc.Encode(res); e != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", e)
			os.Exit(1)
		}
	} else {
		if res.OK {
			fmt.Fprintf(os.Stdout, "OK %s\n", res.Ref)
		} else {
			msg := res.Message
			fmt.Fprintf(os.Stdout, "FAIL %s %s\n", res.Ref, msg)
		}
	}

	if !res.OK {
		// verification failed (mismatch or bundle error)
		os.Exit(2)
	}

	if *strict && !res.OK {
		os.Exit(2)
	}
}
