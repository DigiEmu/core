package main

import (
	"bytes"
	"testing"
)

func TestRunExport_MissingSubcommand_IsUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runExportWithIO([]string{}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit=%d, want 2", exit)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr")
	}
}

func TestRunExport_UnitMissingUnitFlag_IsUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runExportWithIO([]string{"unit"}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("exit=%d, want 2", exit)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout")
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr")
	}
}
