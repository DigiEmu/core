package main

import (
	"bytes"
	"testing"
)

func TestWriteCanonicalJSON_ByteStableAndSorted(t *testing.T) {
	obj := map[string]any{
		"b": 1,
		"a": map[string]any{
			"d": 2,
			"c": 3,
		},
	}

	var buf1 bytes.Buffer
	if err := writeCanonicalJSON(&buf1, obj); err != nil {
		t.Fatalf("writeCanonicalJSON: %v", err)
	}

	var buf2 bytes.Buffer
	if err := writeCanonicalJSON(&buf2, obj); err != nil {
		t.Fatalf("writeCanonicalJSON (second): %v", err)
	}

	if buf1.String() != buf2.String() {
		t.Fatalf("expected identical bytes across runs")
	}

	want := "{\"a\":{\"c\":3,\"d\":2},\"b\":1}\n"
	if buf1.String() != want {
		t.Fatalf("unexpected canonical JSON.\nwant: %q\n got: %q", want, buf1.String())
	}
}
