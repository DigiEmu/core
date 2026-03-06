package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runImportBundleWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu import bundle <bundle-dir>")
		return 2
	}

	dir := args[0]

	info, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(stderr, "import bundle: %v\n", err)
		return 1
	}
	if !info.IsDir() {
		fmt.Fprintln(stderr, "import bundle: bundle dir must be a directory")
		return 1
	}

	required := []string{"bundle.json", "snapshot.json", "meta.json", "trace.json"}
	for _, name := range required {
		p := filepath.Join(dir, name)
		fi, err := os.Stat(p)
		if err != nil {
			fmt.Fprintf(stderr, "import bundle: missing %s\n", name)
			return 1
		}
		if fi.IsDir() {
			fmt.Fprintf(stderr, "import bundle: expected file, got directory %s\n", name)
			return 1
		}
	}

	fmt.Fprintf(stdout, "OK import bundle %s\n", dir)
	return 0
}
