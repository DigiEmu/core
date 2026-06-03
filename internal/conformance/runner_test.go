package conformance

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRunCaseMalformedInputMatchesExpectedError(t *testing.T) {
	r := RunCase(filepath.Join(repoTestdataPath(), "malformed_json_error"))
	if r.Err != "" {
		t.Fatalf("expected malformed_json_error to pass, got error: %s", r.Err)
	}
	if !r.CasePassed {
		t.Fatalf("expected malformed_json_error case_passed true")
	}
	if r.Result != "ERROR" || r.ReasonCode != "INVALID_SNAPSHOT_SCHEMA" {
		t.Fatalf("expected ERROR/INVALID_SNAPSHOT_SCHEMA, got %s/%s", r.Result, r.ReasonCode)
	}
}

func TestRunCaseMalformedInputMismatchFails(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "expected_verify_result.json"), []byte(`{"result":"PASS","reason_code":"STATE_RECONSTRUCTED","verify_result_version":"v2-draft"}`), 0o644)
	os.WriteFile(filepath.Join(dir, "input.json"), []byte(`{`), 0o644)

	r := RunCase(dir)
	if r.Err == "" {
		t.Fatalf("expected observed result mismatch for malformed input")
	}
	if r.CasePassed {
		t.Fatalf("expected case_passed false for malformed input mismatch")
	}
}

func TestRunCaseMalformedInputReasonCodeMismatchFails(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "expected_verify_result.json"), []byte(`{"result":"ERROR","reason_code":"INTERNAL_ERROR","verify_result_version":"v2-draft"}`), 0o644)
	os.WriteFile(filepath.Join(dir, "input.json"), []byte(`{`), 0o644)

	r := RunCase(dir)
	if r.Err == "" {
		t.Fatalf("expected observed reason_code mismatch for malformed input")
	}
	if !strings.Contains(r.Err, "observed result mismatch: got ERROR/INVALID_SNAPSHOT_SCHEMA, expected ERROR/INTERNAL_ERROR") {
		t.Fatalf("expected malformed input to compare observed reason_code through helper, got: %s", r.Err)
	}
	if r.CasePassed {
		t.Fatalf("expected case_passed false for malformed input reason_code mismatch")
	}
}

func TestCompareObservedExpectedExactMatch(t *testing.T) {
	observed := observedResult{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED"}
	expected := Result{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED"}

	if err := compareObservedExpected(observed, expected); err != nil {
		t.Fatalf("expected exact match, got error: %s", err)
	}
}

func TestCompareObservedExpectedResultMismatch(t *testing.T) {
	observed := observedResult{Result: "ERROR", ReasonCode: "INVALID_SNAPSHOT_SCHEMA"}
	expected := Result{Result: "PASS", ReasonCode: "INVALID_SNAPSHOT_SCHEMA"}

	err := compareObservedExpected(observed, expected)
	if err == nil {
		t.Fatalf("expected result mismatch")
	}
	if !strings.Contains(err.Error(), "observed result mismatch") {
		t.Fatalf("expected deterministic mismatch error, got: %s", err)
	}
}

func TestCompareObservedExpectedReasonCodeMismatch(t *testing.T) {
	observed := observedResult{Result: "ERROR", ReasonCode: "INVALID_SNAPSHOT_SCHEMA"}
	expected := Result{Result: "ERROR", ReasonCode: "INTERNAL_ERROR"}

	err := compareObservedExpected(observed, expected)
	if err == nil {
		t.Fatalf("expected reason_code mismatch")
	}
	if !strings.Contains(err.Error(), "observed result mismatch") {
		t.Fatalf("expected deterministic mismatch error, got: %s", err)
	}
}

func TestRunCaseMissingInputFailsDeterministically(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "expected_verify_result.json"), []byte(`{"result":"ERROR","reason_code":"INVALID_SNAPSHOT_SCHEMA","verify_result_version":"v2-draft"}`), 0o644)

	r := RunCase(dir)
	if r.Err == "" {
		t.Fatalf("expected error for missing input.json")
	}
	if !strings.Contains(r.Err, "read input") {
		t.Fatalf("expected missing input error to mention read input, got: %s", r.Err)
	}
	if r.CasePassed {
		t.Fatalf("expected case_passed false for missing input.json")
	}
	if r.ExpectedResult != "ERROR" || r.ReasonCode != "INVALID_SNAPSHOT_SCHEMA" {
		t.Fatalf("expected parsed expected ERROR/INVALID_SNAPSHOT_SCHEMA, got %s/%s", r.ExpectedResult, r.ReasonCode)
	}
}

func TestRunCaseMissingExpectedVerifyResultFailsDeterministically(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "input.json"), []byte(`{}`), 0o644)

	r := RunCase(dir)
	if r.Err == "" {
		t.Fatalf("expected error for missing expected_verify_result.json")
	}
	if !strings.Contains(r.Err, "read expected_verify_result") {
		t.Fatalf("expected missing expected result error to mention expected_verify_result, got: %s", r.Err)
	}
	if r.CasePassed {
		t.Fatalf("expected case_passed false for missing expected_verify_result.json")
	}
	if r.Result != "" || r.ExpectedResult != "" || r.ReasonCode != "" {
		t.Fatalf("expected no fabricated expected result fields, got result=%q expected=%q reason=%q", r.Result, r.ExpectedResult, r.ReasonCode)
	}
}
