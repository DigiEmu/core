package verify

import "digiemu-core/pkg/snapshot"

// Result describes the outcome of a verification.
// JSON shape is stable for CI and audit consumers.
type Result struct {
	OK             bool     `json:"ok"`
	Ref            string   `json:"ref"`
	Expected       string   `json:"expected,omitempty"`
	Got            string   `json:"got,omitempty"`
	HashAlg        string   `json:"hash_alg,omitempty"`
	CanonicalScope string   `json:"canonical_scope,omitempty"`
	Errors         []string `json:"errors,omitempty"`
	Message        string   `json:"message,omitempty"`
	// keep a typed Ref for callers in memory (not serialized)
	_Ref snapshot.Ref `json:"-"`
}

// Verifier is the stable interface for snapshot verification.
// Implementations live in internal/*.
type Verifier interface {
	Verify(ref snapshot.Ref) (Result, error)
}
