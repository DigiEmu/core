package verify

import (
	"testing"

	"digiemu-core/pkg/snapshot"
)

type stubVerifier struct{}

func (stubVerifier) Verify(ref snapshot.Ref) (Result, error) {
	return Result{
		OK:             true,
		Ref:            string(ref.Hash),
		_Ref:           ref,
		HashAlg:        "sha256(canonical_json_v1)",
		CanonicalScope: "canonical_json_v1",
		Trace:          []string{},
		Errors:         []string{},
	}, nil
}

func TestVerifierInterface(t *testing.T) {
	var v Verifier = stubVerifier{}
	ref := snapshot.Ref{Hash: snapshot.Hash("abc")}
	res, err := v.Verify(ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK {
		t.Fatalf("expected OK=true")
	}
	if res._Ref.Hash != ref.Hash {
		t.Fatalf("expected ref to roundtrip")
	}
	if res.CanonicalScope != "canonical_json_v1" {
		t.Fatalf("expected canonical_scope canonical_json_v1, got %s", res.CanonicalScope)
	}
	if res.Trace == nil {
		t.Fatalf("expected trace field present (may be empty)")
	}
	if len(res.Errors) != 0 {
		t.Fatalf("expected no errors in stub verifier")
	}
}
