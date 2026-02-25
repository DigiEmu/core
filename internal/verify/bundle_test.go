package verify

import (
	"path/filepath"
	"runtime"
	"testing"
)

// repoRoot returns <repo>/ by walking up from this test file location.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file = <repo>/internal/verify/bundle_test.go
	// go up: verify -> internal -> repo
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestLoadBundleRootV1_BOMSafeSnapshot(t *testing.T) {
	root := repoRoot(t)

	fixtureRoot := filepath.Join(root, "data", "test-fixtures")
	dataRoot := filepath.Join(root, "data")

	bundleRoot, _, err := FindBundleRoot("demo", fixtureRoot, dataRoot, false)
	if err != nil {
		t.Fatalf("FindBundleRoot: %v", err)
	}

	_, _, err = LoadBundleRootV1(bundleRoot)
	if err != nil {
		t.Fatalf("LoadBundleRootV1: %v", err)
	}
}
