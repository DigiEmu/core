package main

import (
	"encoding/json"
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
		jsonOutput := fs.Bool("json", false, "emit machine-readable JSON conformance report")
		// Support both -json and --json (flag package accepts single-dash only).
		var parsedArgs []string
		jsonDetected := false
		for _, a := range args[1:] {
			if a == "--json" {
				jsonDetected = true
				parsedArgs = append(parsedArgs, "-json")
			} else {
				parsedArgs = append(parsedArgs, a)
			}
		}
		_ = fs.Parse(parsedArgs)
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
			if r.Err != "" {
				failed++
			} else {
				passed++
			}
		}

		effectiveJSON := jsonDetected || *jsonOutput
		if effectiveJSON {
			// build machine-readable report
			type caseReport struct {
				Name           string `json:"name"`
				CasePassed     bool   `json:"case_passed"`
				ExpectedResult string `json:"expected_result"`
				ReasonCode     string `json:"reason_code"`
				Error          string `json:"error,omitempty"`
			}
			rep := map[string]interface{}{
				"report_version": "core-2-conformance-report-v1",
				"status":         "PASS",
				"total":          total,
				"passed":         passed,
				"failed":         failed,
			}
			var cases []caseReport
			for _, r := range results {
				cr := caseReport{
					Name:           r.CaseName,
					CasePassed:     r.CasePassed,
					ExpectedResult: r.ExpectedResult,
					ReasonCode:     r.ReasonCode,
					Error:          r.Err,
				}
				cases = append(cases, cr)
			}
			rep["cases"] = cases
			// overall status FAIL if any failed
			if failed > 0 {
				rep["status"] = "FAIL"
			}
			b, _ := json.MarshalIndent(rep, "", "  ")
			fmt.Fprintln(stdout, string(b))
			if failed > 0 {
				return 3
			}
			return 0
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
