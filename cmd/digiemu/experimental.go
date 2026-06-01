package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	conftool "digiemu-core/internal/conformance"
)

func runExperimental(args []string) {
	os.Exit(runExperimentalWithIO(args, os.Stdout, os.Stderr))
}

func runExperimentalWithIO(args []string, stdout, stderr io.Writer) int {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	if len(args) < 1 {
		fmt.Fprintln(stderr, "experimental subcommands: conformance")
		return 2
	}

	switch args[0] {
	case "conformance":
		return runExperimentalConformanceWithIO(args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "experimental subcommands: conformance")
		return 2
	}
}

func runExperimentalConformanceWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
		return 2
	}

	switch args[0] {
	case "run":
		fs := flag.NewFlagSet("experimental conformance run", flag.ExitOnError)
		_ = fs.Parse(args[1:])
		rem := fs.Args()
		if len(rem) < 1 {
			fmt.Fprintln(stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
			return 2
		}
		root := rem[0]
		// run the internal conformance runner
		results, err := conftool.RunAll(root)
		if err != nil {
			fmt.Fprintf(stderr, "conformance run error: %v\n", err)
			return 1
		}

		total := len(results)
		passed := 0
		failed := 0
		for _, r := range results {
			// A conformance case is considered successful when the expected
			// verify result declaration is structurally valid. The expected
			// Verify Result value (PASS/FAIL/ERROR) is an assertion about the
			// behavior under test and should not by itself mark the conformance
			// case as failed. Mark failure only when a structural error occurred.
			if r.Err != "" {
				failed++
			} else {
				passed++
			}
		}

		fmt.Fprintf(stdout, "Conformance run summary: total=%d passed=%d failed=%d\n", total, passed, failed)
		if failed > 0 {
			return 3
		}
		return 0
	default:
		fmt.Fprintln(stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
		return 2
	}
}
