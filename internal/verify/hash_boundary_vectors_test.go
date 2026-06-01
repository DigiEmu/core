package verify

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	snapshotpkg "digiemu-core/pkg/snapshot"
)

func readJSONFile(t *testing.T, rel string) []byte {
	t.Helper()
	// tests run with package directory as cwd (internal/verify),
	// so reference testdata from the repository root
	p := filepath.Join("..", "..", "testdata", "core_2_hash_boundary", rel)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", p, err)
	}
	return b
}

func TestInsideHashMutationChangesHash(t *testing.T) {
	b1 := readJSONFile(t, filepath.Join("inside_mutation", "base.json"))
	b2 := readJSONFile(t, filepath.Join("inside_mutation", "mutated_inside.json"))

	var s1 map[string]any
	var s2 map[string]any
	if err := json.Unmarshal(b1, &s1); err != nil {
		t.Fatalf("unmarshal base: %v", err)
	}
	if err := json.Unmarshal(b2, &s2); err != nil {
		t.Fatalf("unmarshal mutated: %v", err)
	}

	h1, err := snapshotpkg.HashV1FromState(s1)
	if err != nil {
		t.Fatalf("HashV1FromState base err: %v", err)
	}
	h2, err := snapshotpkg.HashV1FromState(s2)
	if err != nil {
		t.Fatalf("HashV1FromState mutated err: %v", err)
	}
	if h1 == h2 {
		t.Fatalf("inside-hash mutation should change HashV1: %s == %s", h1, h2)
	}
}

func TestExpectedHashV1IsExcludedByReplay_FromTestdata(t *testing.T) {
	// base snapshot (no expected_hash_v1)
	bBase := readJSONFile(t, filepath.Join("outside_metadata", "base.json"))
	// variant with expected_hash_v1 present
	bWith := readJSONFile(t, filepath.Join("outside_metadata", "with_outside_metadata.json"))

	// Build bundle envelopes
	baseBundle := BundleV1{}
	baseBundle.Snapshot = json.RawMessage(bBase)
	withBundle := BundleV1{}
	withBundle.Snapshot = json.RawMessage(bWith)

	rBase, err := ReplayV1(baseBundle, nil)
	if err != nil {
		t.Fatalf("ReplayV1 base err: %v", err)
	}
	rWith, err := ReplayV1(withBundle, nil)
	if err != nil {
		t.Fatalf("ReplayV1 with err: %v", err)
	}

	hBase, err := snapshotpkg.HashV1FromState(rBase.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState base err: %v", err)
	}
	hWith, err := snapshotpkg.HashV1FromState(rWith.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState with err: %v", err)
	}
	if hBase != hWith {
		t.Fatalf("expected ReplayV1 to exclude expected_hash_v1 from hashed scope; hashes differ: %s vs %s", hBase, hWith)
	}
}
