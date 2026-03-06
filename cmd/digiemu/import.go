package main

import (
	"io"
	"os"
)

func runImport(args []string) {
	os.Exit(runImportWithIO(args, os.Stdout, os.Stderr))
}

func runImportWithIO(args []string, stdout, stderr io.Writer) int {
	// Keep this thin and consistent with other command wrappers.
	if len(args) < 1 {
		return runImportBundleWithIO(nil, stdout, stderr)
	}

	switch args[0] {
	case "bundle":
		return runImportBundleWithIO(args[1:], stdout, stderr)
	default:
		return runImportBundleWithIO(nil, stdout, stderr)
	}
}
