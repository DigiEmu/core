package canonicaljson

import (
	"encoding/json"
	"strconv"
	"testing"
)

// These tests document and lock the current canonicalization behavior.

func TestFloatEncodingStable(t *testing.T) {
	v := 1.0
	b, err := Marshal(v)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	// strconv.FormatFloat(1.0, 'g', -1, 64) produces "1"
	want := "1"
	if string(b) != want {
		t.Fatalf("unexpected float encoding: got %s want %s", string(b), want)
	}
}

func TestStringQuotingAndUnicode(t *testing.T) {
	s := "héllo\n\"quote\""
	b, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	// canonical uses strconv.Quote for strings
	if string(b) != strconv.Quote(s) {
		t.Fatalf("string quoting mismatch:\n got:  %s\nwant: %s", string(b), strconv.Quote(s))
	}
}

func TestArrayOrderPreserved(t *testing.T) {
	a := []int{3, 1, 2}
	b, err := Marshal(a)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != "[3,1,2]" {
		t.Fatalf("array order not preserved: %s", string(b))
	}
}

func TestJsonRawMessageEncodedAsBytesArray(t *testing.T) {
	type S struct {
		Raw json.RawMessage `json:"raw"`
	}
	raw := json.RawMessage([]byte(`{"a":1}`))
	s := S{Raw: raw}
	b, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	out := string(b)
	// Current behavior: json.RawMessage is a []byte underlying type and is serialized
	// by canonicaljson as a numeric array of byte values (not embedded JSON string).
	if out == "{\"raw\":{\"a\":1}}" {
		t.Fatalf("json.RawMessage unexpectedly embedded as JSON: %s", out)
	}
	if out == "{}" || out == "" {
		t.Fatalf("unexpected canonical output for json.RawMessage: %s", out)
	}
	// expect numeric array notation starting with '{"raw":['
	if len(out) < 8 || out[:8] != "{\"raw\":[" {
		// allow detected pattern check: must contain '[' indicating byte array
		// if not, fail to make the behavior explicit
		t.Fatalf("json.RawMessage not encoded as byte array as expected: %s", out)
	}
}
