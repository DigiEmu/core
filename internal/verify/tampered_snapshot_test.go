package verify

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	snaps "digiemu-core/pkg/snapshot"
)

func TestVerifyFailsForTamperedSnapshot(t *testing.T) {
	root := repoRoot(t)

	srcFixtureRoot := filepath.Join(root, "data", "test-fixtures")
	tmpFixtureRoot := filepath.Join(t.TempDir(), "test-fixtures")

	copyDir(t, srcFixtureRoot, tmpFixtureRoot)

	snapshotPath := filepath.Join(tmpFixtureRoot, "snapshots", "demo", "snapshot.json")

	raw, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("read snapshot fixture: %v", err)
	}

	var snap map[string]any
	if err := json.Unmarshal(raw, &snap); err != nil {
		t.Fatalf("decode snapshot fixture: %v", err)
	}

	// Keep expected_hash_v1 unchanged, but tamper with deterministic snapshot scope.
	snap["tampered_field"] = "this must change the reconstructed state hash"

	tampered, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		t.Fatalf("encode tampered snapshot fixture: %v", err)
	}

	if err := os.WriteFile(snapshotPath, tampered, 0o644); err != nil {
		t.Fatalf("write tampered snapshot fixture: %v", err)
	}

	v := VerifierV1{
		FixtureRoot: tmpFixtureRoot,
		DataDir:     filepath.Join(root, "data"),
		PreferData:  false,
	}

	result, err := v.Verify(snaps.Ref{Hash: "demo"})
	if err != nil {
		t.Fatalf("Verify returned unexpected error: %v", err)
	}

	if result.OK {
		t.Fatalf("Verify must fail for a tampered snapshot; got OK=true message=%q", result.Message)
	}

	if !strings.Contains(result.Message, "hash mismatch") {
		t.Fatalf("Verify should report hash mismatch for tampered snapshot; got message=%q errors=%v", result.Message, result.Errors)
	}

	if result.Expected == "" {
		t.Fatal("Verify result should include expected hash")
	}

	if result.Got == "" {
		t.Fatal("Verify result should include computed hash")
	}

	if result.Expected == result.Got {
		t.Fatalf("tampered snapshot should produce different hash: expected=%s got=%s", result.Expected, result.Got)
	}
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()

	if err := filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0o644)
	}); err != nil {
		t.Fatalf("copy fixture dir: %v", err)
	}
}
