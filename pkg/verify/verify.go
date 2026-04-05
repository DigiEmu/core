package verify

import "digiemu-core/pkg/snapshot"

// Result describes the outcome of a verification.
// JSON shape is stable for CI and audit consumers.
type Result struct {
	OK             bool     `json:"ok"`
	Ref            string   `json:"ref"`
	Expected       string   `json:"expected"`
	Got            string   `json:"got"`
	HashAlg        string   `json:"hash_alg"`
	CanonicalScope string   `json:"canonical_scope"`
	Trace          []string `json:"trace"`
	Errors         []string `json:"errors,omitempty"`
	Message        string   `json:"message"`
}

// Verifier is the stable interface for snapshot verification.
// Implementations live in internal/*.
type Verifier interface {
	Verify(ref snapshot.Ref) (Result, error)
}
