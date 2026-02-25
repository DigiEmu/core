package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// NOTE:
// - This is an integration test: we build ./cmd/digiemu into a temp binary and execute it.
// - We lock in resolver governance:
//   (1) --bundle is single source of truth and ignores resolver flags
//   (2) when no --bundle: resolver order follows --prefer-data
//   (3) trace MUST contain exactly one final "used:<root>" marker
//   (4) trace MUST contain file paths that were read

func buildCLIBinary(t *testing.T) string {
	t.Helper()

	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}

	tmp := t.TempDir()
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	out := filepath.Join(tmp, "digiemu-test"+ext)

	cmd := exec.Command("go", "build", "-o", out, "./cmd/digiemu")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build cli: %v", err)
	}
	return out
}

func copyDirResolverTest(t *testing.T, src, dst string) {
	t.Helper()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("readdir src: %v", err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatalf("mkdir dst: %v", err)
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			copyDirResolverTest(t, s, d)
			continue
		}
		b, err := os.ReadFile(s)
		if err != nil {
			t.Fatalf("readfile %s: %v", s, err)
		}
		if err := os.WriteFile(d, b, 0o644); err != nil {
			t.Fatalf("writefile %s: %v", d, err)
		}
	}
}

func mustUnmarshalVerifyResult(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(stdout), &m); err != nil {
		t.Fatalf("stdout not valid json: %v\nstdout=%s", err, stdout)
	}
	return m
}

func traceFromResult(t *testing.T, m map[string]any) []string {
	t.Helper()
	raw, ok := m["trace"]
	if !ok {
		t.Fatalf("missing trace in result")
	}
	arr, ok := raw.([]any)
	if !ok {
		t.Fatalf("trace is not array")
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		s, ok := v.(string)
		if !ok {
			t.Fatalf("trace entry is not string: %#v", v)
		}
		out = append(out, s)
	}
	return out
}

func assertTraceGovernance(t *testing.T, trace []string, wantUsedRoot string, wantContains []string) {
	t.Helper()

	// Must contain at least one entry + used marker.
	if len(trace) == 0 {
		t.Fatalf("trace empty")
	}

	// Exactly one used: marker, and it must be the last element.
	usedCount := 0
	usedIdx := -1
	for i, s := range trace {
		if strings.HasPrefix(s, "used:") {
			usedCount++
			usedIdx = i
		}
	}
	if usedCount != 1 {
		t.Fatalf("trace must contain exactly one used: marker, got %d\ntrace=%v", usedCount, trace)
	}
	if usedIdx != len(trace)-1 {
		t.Fatalf("used: marker must be final trace entry, got index %d (len=%d)\ntrace=%v", usedIdx, len(trace), trace)
	}

	gotUsed := strings.TrimPrefix(trace[len(trace)-1], "used:")
	// Normalize path separators for comparison robustness.
	norm := func(p string) string {
		p = filepath.Clean(p)
		p = strings.ReplaceAll(p, "/", string(filepath.Separator))
		p = strings.ReplaceAll(p, "\\", string(filepath.Separator))
		return p
	}
	if norm(gotUsed) != norm(wantUsedRoot) {
		t.Fatalf("used root mismatch\nwant=%s\ngot =%s\ntrace=%v", wantUsedRoot, gotUsed, trace)
	}

	// Must contain expected file paths read.
	all := strings.Join(trace, "\n")
	for _, needle := range wantContains {
		if !strings.Contains(all, needle) {
			t.Fatalf("trace missing expected fragment %q\ntrace=%v", needle, trace)
		}
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("find repo root: %v", err)
	}
	return root
}

func TestVerify_UsesBundleAsSingleSourceOfTruth_IgnoresResolverFlags(t *testing.T) {
	bin := buildCLIBinary(t)
	root := repoRoot(t)

	// Use the tracked demo fixture bundle.
	bundle := filepath.Join(root, "data", "test-fixtures", "snapshots", "demo")

	// Run with nonsense resolver flags; bundle must still be used.
	code, stdout, stderr := runCLI(t, bin,
		"verify",
		"--bundle", bundle,
		"--data", filepath.Join(root, "data"), // should be ignored
		"--fixture-root", filepath.Join(root, "data", "test-fixtures"), // should be ignored
		"--prefer-data",
		"--json",
	)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}
	m := mustUnmarshalVerifyResult(t, stdout)
	tr := traceFromResult(t, m)

	// Trace must end with used:<bundle>
	assertTraceGovernance(
		t,
		tr,
		bundle,
		[]string{
			"snapshot.json",
			string(filepath.Separator) + "claims" + string(filepath.Separator),
		},
	)
}

func TestVerify_ResolverOrder_FollowsPreferData(t *testing.T) {
	bin := buildCLIBinary(t)
	root := repoRoot(t)

	// Prepare two roots in a temp dir:
	// tmpRoot/data/snapshots/demo  (data bundle)
	// tmpRoot/fixtures/snapshots/demo (fixture bundle)
	tmp := t.TempDir()
	dataRoot := filepath.Join(tmp, "data")
	fixtureRoot := filepath.Join(tmp, "fixtures")

	// Copy the tracked demo fixture into BOTH, so verify succeeds in both cases.
	srcDemo := filepath.Join(root, "data", "test-fixtures", "snapshots", "demo")
	dstDataDemo := filepath.Join(dataRoot, "snapshots", "demo")
	dstFixDemo := filepath.Join(fixtureRoot, "snapshots", "demo")
	copyDirResolverTest(t, srcDemo, dstDataDemo)
	copyDirResolverTest(t, srcDemo, dstFixDemo)

	// Make the two bundles distinguishable in trace by adding an extra file in each.
	// Not read by the loader, but we can distinguish by used:<root> selection.
	_ = os.WriteFile(filepath.Join(dstDataDemo, "data_marker.txt"), []byte("data"), 0o644)
	_ = os.WriteFile(filepath.Join(dstFixDemo, "fixture_marker.txt"), []byte("fixture"), 0o644)

	// Case A: prefer-data=false -> fixture first
	codeA, outA, errA := runCLI(t, bin,
		"verify",
		"--ref", "demo",
		"--data", dataRoot,
		"--fixture-root", fixtureRoot,
		"--json",
	)
	if codeA != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout=%s\nstderr=%s", codeA, outA, errA)
	}
	mA := mustUnmarshalVerifyResult(t, outA)
	trA := traceFromResult(t, mA)
	assertTraceGovernance(
		t,
		trA,
		dstFixDemo,
		[]string{
			"snapshot.json",
			string(filepath.Separator) + "claims" + string(filepath.Separator),
		},
	)

	// Case B: prefer-data=true -> data first
	codeB, outB, errB := runCLI(t, bin,
		"verify",
		"--ref", "demo",
		"--data", dataRoot,
		"--fixture-root", fixtureRoot,
		"--prefer-data",
		"--json",
	)
	if codeB != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout=%s\nstderr=%s", codeB, outB, errB)
	}
	mB := mustUnmarshalVerifyResult(t, outB)
	trB := traceFromResult(t, mB)
	assertTraceGovernance(
		t,
		trB,
		dstDataDemo,
		[]string{
			"snapshot.json",
			string(filepath.Separator) + "claims" + string(filepath.Separator),
		},
	)
}

func TestVerify_ResolverNotFound_IsDeterministic(t *testing.T) {
	bin := buildCLIBinary(t)
	tmp := t.TempDir()
	// no snapshots/<ref> exists in either
	dataRoot := filepath.Join(tmp, "data")
	fixtureRoot := filepath.Join(tmp, "fixtures")

	code, out, _ := runCLI(t, bin,
		"verify",
		"--ref", "does-not-exist",
		"--data", dataRoot,
		"--fixture-root", fixtureRoot,
		"--json",
	)

	// We don't hardcode the exact code here if governance evolves,
	// but it MUST be non-zero and JSON parseable.
	if code == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout=%s", out)
	}
	_ = mustUnmarshalVerifyResult(t, out)
}
