package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runExportBundleWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 {
		fmt.Fprintln(stderr, "usage: digiemu export bundle <bundle.json> <out-dir>")
		return 2
	}

	bundlePath := args[0]
	outDir := args[1]

	info, err := os.Stat(bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "export bundle: %v\n", err)
		return 1
	}
	if info.IsDir() {
		fmt.Fprintln(stderr, "export bundle: bundle path must be a file")
		return 1
	}

	srcDir := filepath.Dir(bundlePath)
	baseName := filepath.Base(srcDir)
	targetDir := filepath.Join(outDir, baseName)

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "export bundle: create target dir: %v\n", err)
		return 1
	}

	files := []string{"bundle.json", "snapshot.json", "meta.json", "trace.json"}
	for _, name := range files {
		src := filepath.Join(srcDir, name)
		dst := filepath.Join(targetDir, name)

		data, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintf(stderr, "export bundle: read %s: %v\n", name, err)
			return 1
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			fmt.Fprintf(stderr, "export bundle: write %s: %v\n", name, err)
			return 1
		}
	}

	fmt.Fprintf(stdout, "OK export bundle %s\n", targetDir)
	return 0
}
