package verify

import (
	"encoding/json"
	"testing"

	snapshotpkg "digiemu-core/pkg/snapshot"
)

func TestExpectedHashV1IsExcludedByReplay(t *testing.T) {
	// snapshot with expected_hash_v1 present
	snapWith := map[string]any{"expected_hash_v1": "X", "value": 5}
	bWith := BundleV1{}
	bs, _ := json.Marshal(snapWith)
	bWith.Snapshot = json.RawMessage(bs)

	rWith, err := ReplayV1(bWith, nil)
	if err != nil {
		t.Fatalf("ReplayV1 err: %v", err)
	}

	// snapshot without expected_hash_v1
	snapWithout := map[string]any{"value": 5}
	bWithout := BundleV1{}
	bs2, _ := json.Marshal(snapWithout)
	bWithout.Snapshot = json.RawMessage(bs2)

	rWithout, err := ReplayV1(bWithout, nil)
	if err != nil {
		t.Fatalf("ReplayV1 err: %v", err)
	}

	hWith, err := snapshotpkg.HashV1FromState(rWith.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	hWithout, err := snapshotpkg.HashV1FromState(rWithout.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	if hWith != hWithout {
		t.Fatalf("expected ReplayV1 to exclude expected_hash_v1 from hashed scope; hashes differ: %s vs %s", hWith, hWithout)
	}
}

func TestOutsideMetadataAffectsHashIfNotExcluded(t *testing.T) {
	// Demonstrate that fields other than expected_hash_v1 are currently included in the hashed snapshot
	snapA := map[string]any{"value": 5, "created_at": "t1"}
	snapB := map[string]any{"value": 5, "created_at": "t2"}

	bA := BundleV1{}
	ba, _ := json.Marshal(snapA)
	bA.Snapshot = json.RawMessage(ba)
	rA, err := ReplayV1(bA, nil)
	if err != nil {
		t.Fatalf("ReplayV1 err: %v", err)
	}

	bB := BundleV1{}
	bb, _ := json.Marshal(snapB)
	bB.Snapshot = json.RawMessage(bb)
	rB, err := ReplayV1(bB, nil)
	if err != nil {
		t.Fatalf("ReplayV1 err: %v", err)
	}

	hA, err := snapshotpkg.HashV1FromState(rA.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	hB, err := snapshotpkg.HashV1FromState(rB.StateV1)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	if hA == hB {
		t.Fatalf("expected outside metadata to affect hash if not excluded: hashes equal %s", hA)
	}
}
