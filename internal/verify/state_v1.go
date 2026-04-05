package verify

import "encoding/json"

// StateV1 is the deterministic shape used for hashing assembled bundle state.
type StateV1 struct {
	Snapshot    json.RawMessage   `json:"snapshot"`
	Units       []json.RawMessage `json:"units"`
	Versions    []json.RawMessage `json:"versions"`
	Claims      []json.RawMessage `json:"claims"`
	Meaning     []json.RawMessage `json:"meaning"`
	Uncertainty []json.RawMessage `json:"uncertainty"`
}

// ReconstructedStateV1 is the deterministic replay output shape.
// It matches the StateV1 content and may additionally carry a loader trace.
type ReconstructedStateV1 struct {
	StateV1
	Trace []string `json:"trace"`
}

// StateV1FromBundle assembles the StateV1 from a loaded BundleV1.
func StateV1FromBundle(b BundleV1) StateV1 {
	units := b.Units
	if units == nil {
		units = make([]json.RawMessage, 0)
	}

	versions := b.Versions
	if versions == nil {
		versions = make([]json.RawMessage, 0)
	}

	claims := b.Claims
	if claims == nil {
		claims = make([]json.RawMessage, 0)
	}

	meaning := b.Meaning
	if meaning == nil {
		meaning = make([]json.RawMessage, 0)
	}

	uncertainty := b.Uncertainty
	if uncertainty == nil {
		uncertainty = make([]json.RawMessage, 0)
	}

	return StateV1{
		Snapshot:    b.Snapshot,
		Units:       units,
		Versions:    versions,
		Claims:      claims,
		Meaning:     meaning,
		Uncertainty: uncertainty,
	}
}
