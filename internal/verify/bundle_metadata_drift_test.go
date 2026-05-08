package verify

import (
	"path/filepath"
	"testing"

	snapshotpkg "digiemu-core/pkg/snapshot"
)

func TestReplayHashIgnoresExpectedHashV1Metadata(t *testing.T) {
	root := repoRoot(t)

	fixtureRoot := filepath.Join(root, "data", "test-fixtures")
	dataRoot := filepath.Join(root, "data")

	bundleRoot, _, err := FindBundleRoot("demo", fixtureRoot, dataRoot, false)
	if err != nil {
		t.Fatalf("FindBundleRoot: %v", err)
	}

	originalBundle, _, err := LoadBundleRootV1(bundleRoot)
	if err != nil {
		t.Fatalf("LoadBundleRootV1 original: %v", err)
	}

	modifiedBundle, _, err := LoadBundleRootV1(bundleRoot)
	if err != nil {
		t.Fatalf("LoadBundleRootV1 modified: %v", err)
	}

	modifiedBundle.Snapshot = []byte(`{
		"ref": "demo",
		"expected_hash_v1": "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"state": {
			"hello": "world",
			"n": 1
		}
	}`)

	originalReplayed, err := ReplayV1(originalBundle, []string{"test-original"})
	if err != nil {
		t.Fatalf("ReplayV1 original: %v", err)
	}

	modifiedReplayed, err := ReplayV1(modifiedBundle, []string{"test-modified"})
	if err != nil {
		t.Fatalf("ReplayV1 modified: %v", err)
	}

	originalHash, err := snapshotpkg.HashV1FromState(originalReplayed.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState original: %v", err)
	}

	modifiedHash, err := snapshotpkg.HashV1FromState(modifiedReplayed.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState modified: %v", err)
	}

	if originalHash != modifiedHash {
		t.Fatalf("expected_hash_v1 metadata must be excluded from replay hash scope:\n original: %s\n modified: %s", originalHash, modifiedHash)
	}
}
