package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestExperimentalConformanceSmoke(t *testing.T) {
	// run handler-level invocation with captured output
	var out bytes.Buffer
	var errb bytes.Buffer

	// args: conformance run <path>
	// target the basic_pass case so the smoke test expects success
	repoPath := filepath.Join("..", "..", "testdata", "core_2_conformance", "basic_pass")
	args := []string{"conformance", "run", repoPath}

	code := runExperimentalWithIO(args, &out, &errb)
	if code != 0 {
		t.Fatalf("experimental conformance run exit code %d; stderr=%s", code, errb.String())
	}
	if got := out.String(); got == "" {
		t.Fatalf("expected summary output, got empty")
	}
	if want := "Conformance run summary"; !bytes.Contains(out.Bytes(), []byte(want)) {
		t.Fatalf("expected output to contain %q, got: %s", want, out.String())
	}
}
