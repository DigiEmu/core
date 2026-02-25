package snapshot

// Hash is a stable hex string (e.g., sha256) used as a public reference.
type Hash string

// Ref identifies a snapshot by its hash.
type Ref struct {
	Hash Hash `json:"hash"`
}

// Manifest is a minimal public snapshot manifest (v1.0).
// Internal engines may carry richer state; this is the public stable shape.
type Manifest struct {
	Version string   `json:"version"`
	Hash    Hash     `json:"hash"`
	Inputs  []string `json:"inputs"`
}

func (r Ref) Validate() error {
	if r.Hash == "" {
		return ErrInvalidHash
	}
	return nil
}
