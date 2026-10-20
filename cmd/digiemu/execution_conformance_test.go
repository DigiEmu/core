package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// This file is a test-only execution-conformance prototype. It exercises the
// real CLI verification path (runVerifyWithIO -> Operator.Verify) against
// temporary copies of the measured demo fixture and compares independently
// produced observations with declared expectations. It is NOT an independent
// verifier implementation and does not replace the structural conformance
// runner in internal/conformance.

var snapshotHashPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

// execExpectation declares the v2 contract fields an execution-conformance
// case expects the real verifier to produce. WantExit is the expected CLI
// exit code. FailurePath empty means the field must be absent.
type execExpectation struct {
	Result      string
	ReasonCode  string
	FailurePath string
	WantExit    int
}

// observeVerify executes the real CLI verification path in-process and
// returns the exit code plus the decoded v2 document. The emitted document
// is validated against schemas/verify_result_v2.schema.json.
func observeVerify(t *testing.T, args []string) (int, map[string]any) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	var doc map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("stdout is not valid JSON (exit=%d): %v\nstdout=%q\nstderr=%q", exit, err, stdout.String(), stderr.String())
	}
	if err := compileV2Schema(t).Validate(doc); err != nil {
		t.Fatalf("emitted v2 document failed schema validation: %v\nstdout=%s", err, stdout.String())
	}
	return exit, doc
}

// compareObserved reports a mismatch between the independently observed
// document (and its CLI exit code) and the declared expectation. Observed
// values are never copied into the expectation.
func compareObserved(exit int, doc map[string]any, want execExpectation) error {
	if exit != want.WantExit {
		return fmt.Errorf("exit code: observed %d, expected %d", exit, want.WantExit)
	}
	if got := doc["result"]; got != want.Result {
		return fmt.Errorf("result: observed %v, expected %v", got, want.Result)
	}
	if got := doc["reason_code"]; got != want.ReasonCode {
		return fmt.Errorf("reason_code: observed %v, expected %v", got, want.ReasonCode)
	}
	if want.FailurePath == "" {
		if v, ok := doc["failure_path"]; ok {
			return fmt.Errorf("failure_path: unexpected value %v", v)
		}
	} else if got := doc["failure_path"]; got != want.FailurePath {
		return fmt.Errorf("failure_path: observed %v, expected %v", got, want.FailurePath)
	}
	return nil
}

// copyDemoBundle copies the measured demo fixture into a temp directory and
// registers a cleanup asserting the repository fixture stays byte-identical.
func copyDemoBundle(t *testing.T) (srcDir, dstDir string) {
	t.Helper()
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	srcDir = filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")
	srcSnapshot := filepath.Join(srcDir, "snapshot.json")
	before, err := os.ReadFile(srcSnapshot)
	if err != nil {
		t.Fatalf("read source snapshot: %v", err)
	}
	t.Cleanup(func() {
		after, err := os.ReadFile(srcSnapshot)
		if err != nil {
			t.Errorf("re-read source snapshot: %v", err)
			return
		}
		if !bytes.Equal(before, after) {
			t.Errorf("repository fixture snapshot.json was modified during test")
		}
	})
	dstDir = filepath.Join(t.TempDir(), "snapshots", "demo")
	copyDir(t, srcDir, dstDir)
	return srcDir, dstDir
}

// setExpectedHashV1 rewrites only expected_hash_v1 in a temporary snapshot.
func setExpectedHashV1(t *testing.T, snapshotPath, hash string) {
	t.Helper()
	b, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("read temp snapshot: %v", err)
	}
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("unmarshal temp snapshot: %v", err)
	}
	obj["expected_hash_v1"] = hash
	nb, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatalf("marshal temp snapshot: %v", err)
	}
	if err := os.WriteFile(snapshotPath, nb, 0o644); err != nil {
		t.Fatalf("write temp snapshot: %v", err)
	}
}

// readSnapshotObj loads a snapshot.json into a map for semantic comparison.
func readSnapshotObj(t *testing.T, path string, out *map[string]any) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
}

func v2VerifyArgs(bundle string) []string {
	return []string{"--bundle", bundle, "--verify-result-version=v2-draft", "--json=canonical"}
}

func TestExecConformance_ValidBundle_ObservedPASS(t *testing.T) {
	_, bundle := copyDemoBundle(t)
	want := execExpectation{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED", WantExit: 0}

	exit, doc := observeVerify(t, v2VerifyArgs(bundle))
	if err := compareObserved(exit, doc, want); err != nil {
		t.Fatalf("observation mismatch: %v", err)
	}
	sh, _ := doc["snapshot_hash"].(string)
	if !snapshotHashPattern.MatchString(sh) {
		t.Fatalf("snapshot_hash %q is not a measured sha256 digest", sh)
	}
}

func TestExecConformance_MutatedExpectedHash_ObservedFAIL(t *testing.T) {
	_, bundle := copyDemoBundle(t)
	snapshotPath := filepath.Join(bundle, "snapshot.json")

	// Independently measured baseline must be a clean PASS before its
	// snapshot hash is used as evidence.
	baseExit, baseline := observeVerify(t, v2VerifyArgs(bundle))
	baseWant := execExpectation{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED", WantExit: 0}
	if err := compareObserved(baseExit, baseline, baseWant); err != nil {
		t.Fatalf("baseline observation not a clean PASS: %v", err)
	}
	baselineHash, _ := baseline["snapshot_hash"].(string)
	if !snapshotHashPattern.MatchString(baselineHash) {
		t.Fatalf("baseline snapshot_hash %q is not a measured sha256 digest", baselineHash)
	}

	// Capture semantic content before mutation.
	var before map[string]any
	readSnapshotObj(t, snapshotPath, &before)

	mutated := strings.Repeat("0", 64)
	setExpectedHashV1(t, snapshotPath, mutated)

	// Mutation isolation: everything except expected_hash_v1 must be
	// semantically identical after the rewrite.
	var after map[string]any
	readSnapshotObj(t, snapshotPath, &after)
	if before["expected_hash_v1"] == after["expected_hash_v1"] {
		t.Fatalf("expected_hash_v1 was not mutated: still %v", after["expected_hash_v1"])
	}
	delete(before, "expected_hash_v1")
	delete(after, "expected_hash_v1")
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("mutation touched fields other than expected_hash_v1:\nbefore=%v\nafter=%v", before, after)
	}

	want := execExpectation{
		Result:      "FAIL",
		ReasonCode:  "HASH_MISMATCH",
		FailurePath: "/expected_hash_v1",
		WantExit:    1,
	}
	exit, doc := observeVerify(t, v2VerifyArgs(bundle))
	if err := compareObserved(exit, doc, want); err != nil {
		t.Fatalf("observation mismatch: %v", err)
	}
	if doc["snapshot_hash_expected"] != "sha256:"+mutated {
		t.Fatalf("snapshot_hash_expected=%v, want declared mutation", doc["snapshot_hash_expected"])
	}
	if doc["snapshot_hash_actual"] != baselineHash {
		t.Fatalf("snapshot_hash_actual=%v, want independently measured %v", doc["snapshot_hash_actual"], baselineHash)
	}
}

// TestExecConformance_ComparatorRejectsIncorrectExpectation proves the
// comparison discriminates: deliberately incorrect declared expectations are
// rejected, so a passing case cannot be produced by echoing expectations.
func TestExecConformance_ComparatorRejectsIncorrectExpectation(t *testing.T) {
	_, bundle := copyDemoBundle(t)
	exit, doc := observeVerify(t, v2VerifyArgs(bundle)) // observed PASS document, exit 0

	incorrect := []execExpectation{
		{Result: "FAIL", ReasonCode: "HASH_MISMATCH", FailurePath: "/expected_hash_v1", WantExit: 1},
		{Result: "PASS", ReasonCode: "INTERNAL_ERROR", WantExit: 0},
		{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED", FailurePath: "/bogus", WantExit: 0},
		// All document fields correct but a deliberately wrong exit code.
		{Result: "PASS", ReasonCode: "STATE_RECONSTRUCTED", WantExit: 1},
	}
	for _, bad := range incorrect {
		if err := compareObserved(exit, doc, bad); err == nil {
			t.Fatalf("incorrect expectation %+v was accepted against observed PASS document (exit=%d)", bad, exit)
		}
	}
}
