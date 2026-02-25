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
	if err := ref.Validate(); err != nil {
		return pkgverify.Result{
			OK:      false,
			Message: fmt.Sprintf("invalid ref: %v", err),
			Ref:     ref,
		}, err
	}

	return pkgverify.Result{
		OK:      true,
		Message: "",
		Ref:     ref,
	}, nil
}
