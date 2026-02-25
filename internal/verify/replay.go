package verify

import (
	"encoding/json"
	"fmt"
)

// ReplayV1 deterministically reconstructs the state from a loaded BundleV1.
//
// It uses the same assembly logic as verification (StateV1FromBundle) and performs
// no hash computation.
func ReplayV1(b BundleV1) (ReconstructedStateV1, error) {
	state := StateV1FromBundle(b)

	// Validate that all embedded raw JSON values are syntactically valid.
	// (The bundle loader is BOM-tolerant but does not parse JSON.)
	if !json.Valid(state.Snapshot) {
		return ReconstructedStateV1{}, fmt.Errorf("invalid snapshot json")
	}
	for i, m := range state.Units {
		if !json.Valid(m) {
			return ReconstructedStateV1{}, fmt.Errorf("invalid units[%d] json", i)
		}
	}
	for i, m := range state.Versions {
		if !json.Valid(m) {
			return ReconstructedStateV1{}, fmt.Errorf("invalid versions[%d] json", i)
		}
	}
	for i, m := range state.Claims {
		if !json.Valid(m) {
			return ReconstructedStateV1{}, fmt.Errorf("invalid claims[%d] json", i)
		}
	}
	for i, m := range state.Meaning {
		if !json.Valid(m) {
			return ReconstructedStateV1{}, fmt.Errorf("invalid meaning[%d] json", i)
		}
	}
	for i, m := range state.Uncertainty {
		if !json.Valid(m) {
			return ReconstructedStateV1{}, fmt.Errorf("invalid uncertainty[%d] json", i)
		}
	}

	return ReconstructedStateV1{StateV1: state}, nil
}
