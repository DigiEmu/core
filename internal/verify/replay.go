package verify

import (
	"encoding/json"
	"fmt"
)

// ReplayV1 deterministically reconstructs the replayed state from a loaded BundleV1.
//
// Normative behavior for verification v1:
// - Replay MUST NOT compute the hash.
// - expected_hash_v1 MUST be excluded from the replayed state's snapshot scope.
// - Replay carries forward the input trace (caller-provided).
func ReplayV1(b BundleV1, trace []string) (ReconstructedStateV1, error) {
	// Normalize snapshot scope for hashing: remove expected_hash_v1 to prevent self-reference.
	var snapObj map[string]any
	if err := json.Unmarshal(b.Snapshot, &snapObj); err != nil {
		return ReconstructedStateV1{}, fmt.Errorf("decode snapshot.json object: %v", err)
	}
	delete(snapObj, "expected_hash_v1")

	cleanSnap, err := json.Marshal(snapObj)
	if err != nil {
		return ReconstructedStateV1{}, fmt.Errorf("re-encode snapshot.json: %v", err)
	}
	b.Snapshot = cleanSnap

	state := StateV1FromBundle(b)
	return ReconstructedStateV1{StateV1: state, Trace: trace}, nil
}
