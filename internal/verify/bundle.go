package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func stableTracePath(p string) string {
	// Trace is part of a public, byte-stable demo artifact.
	// Normalize to Windows-style separators to keep outputs stable across OSes.
	return strings.ReplaceAll(p, "/", "\\")
}

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
	Ref            string          `json:"ref"`
	ExpectedHashV1 string          `json:"expected_hash_v1"`
	State          json.RawMessage `json:"state"`
}

// BundleV1 represents an assembled on-disk bundle consisting of snapshot.json
// and optional collections of JSON files.
type BundleV1 struct {
	Snapshot    json.RawMessage   `json:"snapshot"`
	Units       []json.RawMessage `json:"units,omitempty"`
	Versions    []json.RawMessage `json:"versions,omitempty"`
	Claims      []json.RawMessage `json:"claims,omitempty"`
	Meaning     []json.RawMessage `json:"meaning,omitempty"`
	Uncertainty []json.RawMessage `json:"uncertainty,omitempty"`
}

// LoadSnapshotBundleV1 loads the minimal (old) bundle format from <dataDir>/snapshots/<ref>/snapshot.json
// and expects `state` + `expected_hash_v1` fields.
func LoadSnapshotBundleV1(dataDir, ref string) (SnapshotBundleV1, error) {
	p := BundlePath(dataDir, ref)
	return loadSnapshotBundleFromPath(p, ref)
}

// FindBundleRoot attempts to locate the bundle root directory given fixture and data roots.
// It returns the root dir used, trace of attempted roots, or error.
func FindBundleRoot(ref, fixtureRoot, dataRoot string, preferData bool) (string, []string, error) {
	fixtureDir := filepath.Join(fixtureRoot, "snapshots", ref)
	dataDir := filepath.Join(dataRoot, "snapshots", ref)

	tryOrder := []string{fixtureDir, dataDir}
	if preferData {
		tryOrder = []string{dataDir, fixtureDir}
	}

	var attempts []string
	var lastErr error
	for _, d := range tryOrder {
		attempts = append(attempts, stableTracePath(d))
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			return d, attempts, nil
		} else if err != nil {
			lastErr = err
			continue
		}
	}
	return "", attempts, fmt.Errorf("no bundle root found (attempted %d roots); last error: %w", len(attempts), lastErr)
}

// ReadJSONFileBOMSafe reads a JSON file and strips UTF-8 BOM if present.
func ReadJSONFileBOMSafe(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return stripUTF8BOM(b), nil
}

// LoadBundleRootV1 reads snapshot.json from root and optional collection directories.
// It returns the assembled BundleV1, a trace of loaded file paths, and error.
func LoadBundleRootV1(root string) (BundleV1, []string, error) {
	var trace []string

	snapshotPath := filepath.Join(root, "snapshot.json")
	b, err := ReadJSONFileBOMSafe(snapshotPath)
	trace = append(trace, stableTracePath(snapshotPath))
	if err != nil {
		return BundleV1{}, trace, fmt.Errorf("read snapshot.json: %w", err)
	}

	// keep snapshot raw (verifier may normalize it before hashing)
	bundle := BundleV1{Snapshot: json.RawMessage(b)}

	// optional dirs to load
	optDirs := []string{"units", "versions", "claims", "meaning", "uncertainty"}
	for _, d := range optDirs {
		dirPath := filepath.Join(root, d)
		fi, err := os.Stat(dirPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return BundleV1{}, trace, fmt.Errorf("stat %s: %w", dirPath, err)
		}
		if !fi.IsDir() {
			continue
		}

		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return BundleV1{}, trace, fmt.Errorf("read dir %s: %w", dirPath, err)
		}

		var names []string
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
				names = append(names, e.Name())
			}
		}

		sort.Strings(names)
		for _, name := range names {
			p := filepath.Join(dirPath, name)
			trace = append(trace, stableTracePath(p))

			rb, err := ReadJSONFileBOMSafe(p)
			if err != nil {
				return BundleV1{}, trace, fmt.Errorf("read file %s: %w", p, err)
			}

			raw := json.RawMessage(rb)
			switch d {
			case "units":
				bundle.Units = append(bundle.Units, raw)
			case "versions":
				bundle.Versions = append(bundle.Versions, raw)
			case "claims":
				bundle.Claims = append(bundle.Claims, raw)
			case "meaning":
				bundle.Meaning = append(bundle.Meaning, raw)
			case "uncertainty":
				bundle.Uncertainty = append(bundle.Uncertainty, raw)
			}
		}
	}

	trace = append(trace, fmt.Sprintf("used:%s", stableTracePath(root)))
	return bundle, trace, nil
}

// --- minimal snapshot bundle helpers (old format) ---

func loadSnapshotBundleFromPath(path, ref string) (SnapshotBundleV1, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return SnapshotBundleV1{}, fmt.Errorf("read bundle %s: %w", path, err)
	}

	// strip UTF-8 BOM if present to be tolerant of Windows editors
	b = stripUTF8BOM(b)

	var sb SnapshotBundleV1
	if err := json.Unmarshal(b, &sb); err != nil {
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
