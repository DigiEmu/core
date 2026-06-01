package snapshot

import "testing"

func TestHashChangesOnInsideMutation(t *testing.T) {
	state1 := map[string]any{"a": 1, "b": map[string]any{"x": 10}}
	state2 := map[string]any{"a": 1, "b": map[string]any{"x": 11}} // inside-hash change

	h1, err := HashV1FromState(state1)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	h2, err := HashV1FromState(state2)
	if err != nil {
		t.Fatalf("HashV1FromState err: %v", err)
	}
	if h1 == h2 {
		t.Fatalf("expected different hashes for inside-hash mutation: %s == %s", h1, h2)
	}
}
