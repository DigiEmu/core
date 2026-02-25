package snapshot

import "testing"

func TestRefValidate_EmptyHash(t *testing.T) {
	var r Ref
	if err := r.Validate(); err == nil {
		t.Fatalf("expected error for empty hash")
	}
}

func TestRefValidate_NonEmptyHash(t *testing.T) {
	r := Ref{Hash: Hash("abc")}
	if err := r.Validate(); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
