package verify

import (
	"path/filepath"
	"testing"

	snapshotpkg "digiemu-core/pkg/snapshot"
)

func TestGoldenSnapshotHashV1(t *testing.T) {
	const goldenHashV1 = "273910B657210EECE942EFC038F6E06994711DA6399422E57225CE2394299E55"

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

	replayed, err := ReplayV1(bundle, []string{"test-golden-snapshot-hash-v1"})
	if err != nil {
		t.Fatalf("ReplayV1: %v", err)
	}

	hash, err := snapshotpkg.HashV1FromState(replayed.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState: %v", err)
	}

	if string(hash) != goldenHashV1 {
		t.Fatalf("golden snapshot hash changed:\n got:  %s\nwant: %s", hash, goldenHashV1)
	}
}
