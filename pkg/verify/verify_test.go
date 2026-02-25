package verify

import (
	"testing"

	"digiemu-core/pkg/snapshot"
)

type stubVerifier struct{}

func (stubVerifier) Verify(ref snapshot.Ref) (Result, error) {
	return Result{OK: true, Ref: string(ref.Hash), _Ref: ref}, nil
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
}
