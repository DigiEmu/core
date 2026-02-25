package canonicaljson

import (
	"testing"
)

func TestMapKeyOrdering(t *testing.T) {
	m := map[string]int{"b": 2, "a": 1}
	b, err := Marshal(m)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	s := string(b)
	if s != "{\"a\":1,\"b\":2}" {
		t.Fatalf("unexpected canonical output: %s", s)
	}
}

func TestStructFieldOrder(t *testing.T) {
	type S struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	s := S{A: 1, B: 2}
	b, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != "{\"a\":1,\"b\":2}" {
		t.Fatalf("unexpected struct output: %s", string(b))
	}
}
