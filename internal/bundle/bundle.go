package bundle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// WriteExpectedHashV1 updates the expected_hash_v1 field in snapshot.json at snapshotPath.
//
// Behavior matches the existing verify writer:
// - accepts empty/known placeholder current values
// - refuses overwriting a non-placeholder existing expected
// - is tolerant of UTF-8 BOM
//
// Error typing: returned errors wrap sentinel errors from errors.go.
func WriteExpectedHashV1(snapshotPath string, newHash string) error {
	b, err := os.ReadFile(snapshotPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %w", ErrSnapshotNotFound, err)
		}
		return fmt.Errorf("read snapshot: %w", err)
	}

	b = stripUTF8BOM(b)

	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		return fmt.Errorf("%w: %w", ErrSnapshotInvalidJSON, err)
	}

	cur, _ := obj["expected_hash_v1"].(string)
	if !isPlaceholderExpected(cur) {
		return fmt.Errorf("%w: %s", ErrExpectedAlreadySet, cur)
	}

	// Note: newHash validation is performed by the operator; writer preserves previous behavior.
	obj["expected_hash_v1"] = newHash

	nb, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("encode snapshot: %w", err)
	}
	if err := os.WriteFile(snapshotPath, nb, 0644); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w: %w", ErrSnapshotNotFound, err)
		}
		return fmt.Errorf("write snapshot: %w", err)
	}
	return nil
}

func isPlaceholderExpected(s string) bool {
	if s == "" {
		return true
	}
	switch s {
	case "REPLACE_AFTER_ASSEMBLE", "REPLACE_AFTER_COMPUTE":
		return true
	default:
		return false
	}
}

func stripUTF8BOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}
