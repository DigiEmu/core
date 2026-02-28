package verify

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	bun "digiemu-core/internal/bundle"
	derr "digiemu-core/pkg/digiemu"
	"digiemu-core/pkg/snapshot"
)

func TestOperatorVerify_WriteExpectedAlreadySet_IsCategorized(t *testing.T) {
	tmp := t.TempDir()

	ref := strings.Repeat("b", 64)
	root := filepath.Join(tmp, "snapshots", ref)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	snapshotPath := filepath.Join(root, "snapshot.json")
	expectedAlreadySet := strings.Repeat("a", 64)
	content := "{\n" +
		"  \"expected_hash_v1\": \"" + expectedAlreadySet + "\",\n" +
		"  \"foo\": 1\n" +
		"}\n"
	if err := os.WriteFile(snapshotPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write snapshot.json: %v", err)
	}

	op := &Operator{
		DataDir:     tmp,
		FixtureRoot: tmp,
		PreferData:  false,
		Options: Options{
			WriteExpected: true,
		},
	}

	res, err := op.Verify(snapshot.Ref{Hash: snapshot.Hash(ref)})
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !res.WriteBlocked {
		t.Fatalf("expected write_blocked")
	}
	if res.WriteReason != WriteReasonExistingExpected {
		t.Fatalf("expected write_reason=%q, got %q", WriteReasonExistingExpected, res.WriteReason)
	}

	// Preserve existing sentinel classification used by CLI exit code mapping.
	if !errors.Is(err, bun.ErrExpectedAlreadySet) {
		t.Fatalf("expected errors.Is bun.ErrExpectedAlreadySet")
	}

	// Attach new stable taxonomy classification.
	if !errors.Is(err, derr.ErrUsage) {
		t.Fatalf("expected errors.Is derr.ErrUsage")
	}
	if code, ok := derr.CodeOf(err); !ok || code != derr.CodeUsage {
		t.Fatalf("expected CodeOf=USAGE, got %q (ok=%v)", code, ok)
	}
}
