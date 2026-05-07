package verify

import (
	"path/filepath"
	"testing"

	snapshotpkg "digiemu-core/pkg/snapshot"
)

func TestReplayV1Idempotent(t *testing.T) {
	root := repoRoot(t)

	fixtureRoot := filepath.Join(root, "data", "test-fixtures")
	dataRoot := filepath.Join(root, "data")

	bundleRoot, _, err := FindBundleRoot("demo", fixtureRoot, dataRoot, false)
	if err != nil {
		t.Fatalf("FindBundleRoot: %v", err)
	}

	bundle, _, err := LoadBundleRootV1(bundleRoot)
	if err != nil {
		t.Fatalf("LoadBundleRootV1: %v", err)
	}

	var first snapshotpkg.HashV1

	for i := 0; i < 10; i++ {
		replayed, err := ReplayV1(bundle, []string{"test-replay-idempotent"})
		if err != nil {
			t.Fatalf("ReplayV1 iteration %d: %v", i, err)
		}

		hash, err := snapshotpkg.HashV1FromState(replayed.StateV1)
		if err != nil {
			t.Fatalf("HashV1FromState iteration %d: %v", i, err)
		}

		if i == 0 {
			first = hash
			continue
		}

		if hash != first {
			t.Fatalf("ReplayV1 must be idempotent: iteration 0 hash=%s iteration %d hash=%s", first, i, hash)
		}
	}
}
