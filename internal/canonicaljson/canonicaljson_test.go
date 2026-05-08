package canonicaljson

import (
	"encoding/json"
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

func TestNestedMapKeyOrdering(t *testing.T) {
	m := map[string]any{
		"z": 3,
		"a": map[string]any{
			"b": 2,
			"a": 1,
		},
		"m": []any{
			map[string]any{
				"y": 2,
				"x": 1,
			},
		},
	}

	b, err := Marshal(m)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	got := string(b)
	want := "{\"a\":{\"a\":1,\"b\":2},\"m\":[{\"x\":1,\"y\":2}],\"z\":3}"

	if got != want {
		t.Fatalf("unexpected nested canonical output:\n got: %s\nwant: %s", got, want)
	}
}

func TestEquivalentNestedMapsWithDifferentInputOrderProduceSameCanonicalJSON(t *testing.T) {
	a := map[string]any{
		"outer_b": map[string]any{
			"beta":  2,
			"alpha": 1,
		},
		"outer_a": []any{
			map[string]any{
				"delta": 4,
				"gamma": 3,
			},
		},
	}

	b := map[string]any{
		"outer_a": []any{
			map[string]any{
				"gamma": 3,
				"delta": 4,
			},
		},
		"outer_b": map[string]any{
			"alpha": 1,
			"beta":  2,
		},
	}

	canonA, err := Marshal(a)
	if err != nil {
		t.Fatalf("Marshal A error: %v", err)
	}

	canonB, err := Marshal(b)
	if err != nil {
		t.Fatalf("Marshal B error: %v", err)
	}

	if string(canonA) != string(canonB) {
		t.Fatalf("equivalent nested maps must produce identical canonical JSON:\n A: %s\n B: %s", canonA, canonB)
	}
}

func TestCanonicalJSONIgnoresInputWhitespace(t *testing.T) {
	compactJSON := []byte(`{"a":1,"b":{"x":2,"y":[3,4]}}`)

	formattedJSON := []byte(`{
		"b": {
			"y": [
				3,
				4
			],
			"x": 2
		},
		"a": 1
	}`)

	var compact any
	if err := json.Unmarshal(compactJSON, &compact); err != nil {
		t.Fatalf("unmarshal compact JSON: %v", err)
	}

	var formatted any
	if err := json.Unmarshal(formattedJSON, &formatted); err != nil {
		t.Fatalf("unmarshal formatted JSON: %v", err)
	}

	canonCompact, err := Marshal(compact)
	if err != nil {
		t.Fatalf("Marshal compact JSON: %v", err)
	}

	canonFormatted, err := Marshal(formatted)
	if err != nil {
		t.Fatalf("Marshal formatted JSON: %v", err)
	}

	if string(canonCompact) != string(canonFormatted) {
		t.Fatalf("canonical JSON must ignore input whitespace:\n compact:   %s\n formatted: %s", canonCompact, canonFormatted)
	}

	want := `{"a":1,"b":{"x":2,"y":[3,4]}}`
	if string(canonCompact) != want {
		t.Fatalf("unexpected canonical output:\n got:  %s\nwant: %s", canonCompact, want)
	}
}
