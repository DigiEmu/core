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
	format := fs.String("format", "text", "Output format: text|json")
	strict := fs.Bool("strict", false, "Exit non-zero if verification fails")
	fs.Parse(args)

	if strings.TrimSpace(*refStr) == "" {
		fmt.Fprintln(os.Stderr, "--ref is required")
		fs.Usage()
		os.Exit(2)
	}

	ref := snapshot.Ref{Hash: snapshot.Hash(strings.TrimSpace(*refStr))}
	v := verify.StubVerifier{}
	res, err := v.Verify(ref)

	switch strings.ToLower(*format) {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if e := enc.Encode(res); e != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", e)
			os.Exit(1)
		}
	case "text":
		if res.OK {
			fmt.Fprintf(os.Stdout, "OK\t%s\n", res.Ref.Hash)
		} else {
			msg := res.Message
			if msg == "" && err != nil {
				msg = err.Error()
			}
			fmt.Fprintf(os.Stdout, "FAIL\t%s\t%s\n", res.Ref.Hash, msg)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid --format: %q (use json|text)\n", *format)
		os.Exit(2)
	}

	if *strict && (!res.OK || err != nil) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "verification failed")
		os.Exit(1)
	}
}
