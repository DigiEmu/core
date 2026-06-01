package conformance

import (
	"os"
	"path/filepath"
	"testing"
)

func repoTestdataPath() string {
	// tests run with package directory as cwd; reach up to repo testdata
	return filepath.Join("..", "..", "testdata", "core_2_conformance")
}

func TestDiscoverCasesFindsAll(t *testing.T) {
	root := repoTestdataPath()
	cases, err := DiscoverCases(root)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(cases) == 0 {
		t.Fatalf("expected at least one conformance case, got 0")
	}
}

func TestRunAllValidatesExpectedResults(t *testing.T) {
	root := repoTestdataPath()
	results, err := RunAll(root)
	if err != nil {
		t.Fatalf("RunAll: %v", err)
	}
	// ensure each result has either Result set or an error
	found := 0
	for _, r := range results {
		if r.Err != "" {
			t.Fatalf("case %s reported error: %s", r.CaseName, r.Err)
		}
		if r.Result == "" || r.ReasonCode == "" || r.VerifyResultVersion == "" {
			t.Fatalf("case %s missing required fields", r.CaseName)
		}
		found++
	}
	if found == 0 {
		t.Fatalf("no cases validated")
	}
}

func TestValidateCaseFilesMissingProducesError(t *testing.T) {
	dir := t.TempDir()
	// create only input.json, omit expected_verify_result.json
	os.WriteFile(filepath.Join(dir, "input.json"), []byte("{}"), 0o644)
	err := ValidateCaseFiles(dir)
	if err == nil {
		t.Fatalf("expected error for missing expected_verify_result.json")
	}
}

func TestRunCaseInvalidResultValue(t *testing.T) {
	dir := t.TempDir()
	// write expected_verify_result with invalid result value
	os.WriteFile(filepath.Join(dir, "expected_verify_result.json"), []byte(`{"result":"BAD","reason_code":"X","verify_result_version":"v"}`), 0o644)
	// input.json required for discovery but not used for RunCase test
	os.WriteFile(filepath.Join(dir, "input.json"), []byte(`{}`), 0o644)

	r := RunCase(dir)
	if r.Err == "" {
		t.Fatalf("expected error for invalid result value")
	}
}
