package meaning

// NodeID is a stable identifier for a meaning node.
// v1.0: callers must treat IDs as stable inputs (no random generation in deterministic paths).
type NodeID string

// Node is the minimal public representation of a meaning unit.
type Node struct {
	ID   NodeID `json:"id"`
	Text string `json:"text"`
}

// Validate performs minimal invariants for v1.0.
func (n Node) Validate() error {
	if n.ID == "" {
		return ErrInvalidID
	}
	return nil
}
