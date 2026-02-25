package snapshot

import (
	"testing"
)

func TestHashV1Deterministic(t *testing.T) {
	state1 := map[string]any{
		"unit": map[string]any{"key": "u1", "versions": []any{map[string]any{"v": 1}}},
	}
	state2 := map[string]any{
		"unit": map[string]any{"versions": []any{map[string]any{"v": 1}}, "key": "u1"},
	}

	h1, err := HashV1FromState(state1)
	if err != nil {
		t.Fatalf("hash1 err: %v", err)
	}
	h2, err := HashV1FromState(state2)
	if err != nil {
		t.Fatalf("hash2 err: %v", err)
	}
	if h1 != h2 {
		t.Fatalf("hash mismatch: %s != %s", h1, h2)
	}
}
