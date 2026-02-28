package main

import (
	"testing"

	derr "digiemu-core/pkg/digiemu"
)

func TestExitCodeForError_Contract(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: 0},
		{name: "hash mismatch", err: derr.HashMismatch("verify.hash", "", "a", "b", nil), want: 1},
		{name: "schema invalid", err: derr.SchemaInvalid("verify.schema", "", "x", nil), want: 1},
		{name: "file missing", err: derr.FileMissing("read.file", "x", nil), want: 1},
		{name: "bundle invalid", err: derr.BundleInvalid("verify.bundle", "x", nil), want: 1},
		{name: "verify failed", err: derr.VerifyFailed("verify", "x", nil), want: 1},
		{name: "usage", err: derr.UsageError("cli", "", nil), want: 2},
		{name: "io", err: derr.IOError("read", "x", nil), want: 3},
		{name: "internal", err: derr.InternalError("boom", nil), want: 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := exitCodeForError(tc.err)
			if got != tc.want {
				t.Fatalf("want %d, got %d", tc.want, got)
			}
		})
	}
}
