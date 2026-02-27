package main

import "fmt"

// These values are injected at build time via -ldflags.
// Defaults are "unknown" for local/dev runs.
var Version = "unknown"
var Commit = "unknown"
var Date = "unknown"

func cliVersionLine() string {
	return fmt.Sprintf("digiemu %s (%s) %s", Version, Commit, Date)
}
