package verify

import (
	"fmt"

	"digiemu-core/pkg/snapshot"
	pkgverify "digiemu-core/pkg/verify"
)

// StubVerifier is a deterministic placeholder verifier.
// Phase 6 will replace this with a real engine-backed implementation.
type StubVerifier struct{}

func (StubVerifier) Verify(ref snapshot.Ref) (pkgverify.Result, error) {
	// Deterministic behavior:
	// - invalid ref => error
	// - otherwise => OK=true
	r := pkgverify.Result{
		OK:             false,
		Ref:            string(ref.Hash),
		HashAlg:        "sha256(canonical_json_v1)",
		CanonicalScope: "canonical_json_v1",
		Trace:          []string{},
	}

	if err := ref.Validate(); err != nil {
		r.Message = fmt.Sprintf("invalid ref: %v", err)
		r.Errors = append(r.Errors, r.Message)
		return r, err
	}

	r.OK = true
	r.Message = "ok"
	return r, nil
}
