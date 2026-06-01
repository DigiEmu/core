package main

import (
	"flag"
	"fmt"
	"os"

	conftool "digiemu-core/internal/conformance"
)

func runExperimental(args []string) {
	os.Exit(runExperimentalWithIO(args, os.Stdout, os.Stderr))
}

func runExperimentalWithIO(args []string, stdout, stderr interface{}) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "experimental subcommands: conformance")
		return 2
	}

	switch args[0] {
	case "conformance":
		return runExperimentalConformance(args[1:])
	default:
		fmt.Fprintln(os.Stderr, "experimental subcommands: conformance")
		return 2
	}
}

func runExperimentalConformance(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
		return 2
	}

	switch args[0] {
	case "run":
		fs := flag.NewFlagSet("experimental conformance run", flag.ExitOnError)
		_ = fs.Parse(args[1:])
		rem := fs.Args()
		if len(rem) < 1 {
			fmt.Fprintln(os.Stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
			return 2
		}
		root := rem[0]
		// run the internal conformance runner
		results, err := conftool.RunAll(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "conformance run error: %v\n", err)
			return 1
		}

		total := len(results)
		passed := 0
		failed := 0
		for _, r := range results {
			if r.Err != "" || (r.Result != "PASS") {
				failed++
			} else {
				passed++
			}
		}

		fmt.Printf("Conformance run summary: total=%d passed=%d failed=%d\n", total, passed, failed)
		if failed > 0 {
			return 3
		}
		return 0
	default:
		fmt.Fprintln(os.Stderr, "usage: digiemu experimental conformance run <path-to-conformance-dir>")
		return 2
	}
}
