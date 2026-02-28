package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func runExport(args []string) {
	os.Exit(runExportWithIO(args, os.Stdout, os.Stderr))
}

func runExportWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "export subcommands: unit")
		return 2
	}

	switch args[0] {
	case "unit":
		fs := flag.NewFlagSet("export unit", flag.ContinueOnError)
		fs.SetOutput(stderr)
		unitKey := fs.String("unit", "", "unit key (required)")
		data := fs.String("data", "./data", "data directory")
		withAudit := fs.Bool("audit", false, "include audit events for this unit")
		pretty := fs.Bool("pretty", false, "pretty-print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			if err == flag.ErrHelp {
				return 0
			}
			return 2
		}

		if *unitKey == "" {
			fmt.Fprintln(stderr, "--unit is required")
			fs.Usage()
			return 2
		}

		repo := fsrepo.NewUnitRepo(*data)

		var audit ports.AuditLogByUnitReader
		if *withAudit {
			audit = fsrepo.NewAuditByUnitReader(*data)
		}

		uc := usecases.ExportUnitSnapshot{Repo: repo, Audit: audit}

		out, err := uc.ExportUnitSnapshot(ports.ExportUnitSnapshotRequest{
			UnitKey:      *unitKey,
			IncludeAudit: *withAudit,
		})
		if err != nil {
			fmt.Fprintf(stderr, "export unit: %v\n", err)
			return 4
		}

		var b []byte
		if *pretty {
			b, err = json.MarshalIndent(out, "", "  ")
		} else {
			b, err = json.Marshal(out)
		}
		if err != nil {
			fmt.Fprintf(stderr, "export marshal: %v\n", err)
			return 4
		}
		fmt.Fprintln(stdout, string(b))
		return 0

	default:
		fmt.Fprintln(stderr, "export subcommands: unit")
		return 2
	}
}
