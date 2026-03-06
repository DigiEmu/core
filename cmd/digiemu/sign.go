package main

import (
	"os"
)

func runSign(args []string) {
	if len(args) < 1 {
		os.Exit(runSignBundleWithIO(nil, os.Stdout, os.Stderr))
	}

	switch args[0] {
	case "bundle":
		os.Exit(runSignBundleWithIO(args[1:], os.Stdout, os.Stderr))
	default:
		os.Exit(runSignBundleWithIO(nil, os.Stdout, os.Stderr))
	}
}
