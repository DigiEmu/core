package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// BundlePath returns the expected location of a snapshot bundle.
// MVP layout:
//
//	<data>/snapshots/<ref>/snapshot.json
func BundlePath(dataDir, ref string) string {
	return filepath.Join(dataDir, "snapshots", ref, "snapshot.json")
}

// SnapshotBundleV1 is the MVP on-disk bundle format.
// NOTE: This is deliberately minimal to make `verify` usable for demos.
type SnapshotBundleV1 struct {
	Ref            string      `json:"ref"`
	ExpectedHashV1 string      `json:"expected_hash_v1"`
	State          interface{} `json:"state"`
}

func LoadBundleV1(dataDir, ref string) (SnapshotBundleV1, error) {
	p := BundlePath(dataDir, ref)
	return loadBundleFromPath(p, ref)
}

// loadBundleFromPath reads and decodes a bundle from a specific path.
func loadBundleFromPath(path, ref string) (SnapshotBundleV1, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return SnapshotBundleV1{}, fmt.Errorf("read bundle %s: %w", path, err)
	}

	// strip UTF-8 BOM if present to be tolerant of Windows editors
	b = stripUTF8BOM(b)

	var sb SnapshotBundleV1
	if err := json.Unmarshal(b, &sb); err != nil {
		// If BOM was detected originally, mention it to aid debugging
		if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
			return SnapshotBundleV1{}, fmt.Errorf("decode bundle %s: %w (note: UTF-8 BOM detected and stripped)", path, err)
		}
		return SnapshotBundleV1{}, fmt.Errorf("decode bundle %s: %w", path, err)
	}

	if sb.Ref == "" {
		sb.Ref = ref
	}

	if sb.ExpectedHashV1 == "" {
		return SnapshotBundleV1{}, fmt.Errorf("bundle missing expected_hash_v1")
	}
	if sb.State == nil {
		return SnapshotBundleV1{}, fmt.Errorf("bundle missing state")
	}
	return sb, nil
}

// stripUTF8BOM removes a UTF-8 BOM prefix if present.
func stripUTF8BOM(b []byte) []byte {
	if len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		return b[3:]
	}
	return b
}

// FindBundleV1 attempts to locate a bundle file in fixtureRoot and dataRoot.
// If preferData is true, dataRoot is preferred when both exist.
// Returns the loaded bundle, the chosen path, list of attempted paths (in order checked), and error.
func FindBundleV1(fixtureRoot, dataRoot, ref string, preferData bool) (SnapshotBundleV1, string, []string, error) {
	fixturePath := filepath.Join(fixtureRoot, "snapshots", ref, "snapshot.json")
	dataPath := filepath.Join(dataRoot, "snapshots", ref, "snapshot.json")

	var attempts []string
	tryOrder := []string{fixturePath, dataPath}
	if preferData {
		tryOrder = []string{dataPath, fixturePath}
	}

	var lastErr error
	for _, p := range tryOrder {
		attempts = append(attempts, p)
		if _, err := os.Stat(p); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// continue to next
				lastErr = err
				continue
			}
			lastErr = err
			continue
		}
		sb, err := loadBundleFromPath(p, ref)
		if err != nil {
			return SnapshotBundleV1{}, p, attempts, err
		}
		return sb, p, attempts, nil
	}
	return SnapshotBundleV1{}, "", attempts, fmt.Errorf("no bundle found (attempted %d paths); last error: %w", len(attempts), lastErr)
}
