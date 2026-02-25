package verify

import "github.com/BrunoBaumgartner78/digiemu-core/pkg/snapshot"

// Result describes the outcome of a verification.
type Result struct {
	OK      bool         `json:"ok"`
	Message string       `json:"message,omitempty"`
	Ref     snapshot.Ref `json:"ref"`
}

// Verifier is the stable interface for snapshot verification.
// Implementations live in internal/*.
type Verifier interface {
	Verify(ref snapshot.Ref) (Result, error)
}
