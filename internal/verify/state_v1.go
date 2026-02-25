package verify

import "encoding/json"

// StateV1 is the deterministic shape used for hashing assembled bundle state.
type StateV1 struct {
	Snapshot    json.RawMessage   `json:"snapshot"`
	Units       []json.RawMessage `json:"units,omitempty"`
	Versions    []json.RawMessage `json:"versions,omitempty"`
	Claims      []json.RawMessage `json:"claims,omitempty"`
	Meaning     []json.RawMessage `json:"meaning,omitempty"`
	Uncertainty []json.RawMessage `json:"uncertainty,omitempty"`
}

// ReconstructedStateV1 is the deterministic replay output shape.
// It matches the StateV1 content and may additionally carry a loader trace.
type ReconstructedStateV1 struct {
	StateV1
	Trace []string `json:"trace,omitempty"`
}

// StateV1FromBundle assembles the StateV1 from a loaded BundleV1.
func StateV1FromBundle(b BundleV1) StateV1 {
	return StateV1{
		Snapshot:    b.Snapshot,
		Units:       b.Units,
		Versions:    b.Versions,
		Claims:      b.Claims,
		Meaning:     b.Meaning,
		Uncertainty: b.Uncertainty,
	}
}
