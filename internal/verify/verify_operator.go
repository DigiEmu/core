package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	bun "digiemu-core/internal/bundle"
	snaps "digiemu-core/pkg/snapshot"
	pkgverify "digiemu-core/pkg/verify"
)

// newResultV1 constructs a ResultV1 with write-policy defaults.
//
// Defaults (required by contract):
// - wrote_expected=false
// - write_blocked=false
// - write_reason="none"
func newResultV1(base pkgverify.Result) ResultV1 {
	r := ResultV1{Result: base}
	r.WroteExpected = false
	r.WriteBlocked = false
	r.WriteReason = WriteReasonNone
	return r
}

// Operator verifies snapshot bundles and optionally writes expected hashes according to policy.
type Operator struct {
	DataDir     string
	FixtureRoot string
	PreferData  bool
	Options     Options
}

// Verify runs verification and, when enabled, can write expected_hash_v1 to snapshot.json
// ONLY when a placeholder is present and the computed hash is valid.
func (op *Operator) Verify(ref snaps.Ref) (ResultV1, error) {
	refStr := string(ref.Hash)

	fixtureRoot := op.FixtureRoot
	if fixtureRoot == "" {
		fixtureRoot = "data/test-fixtures"
	}

	base := pkgverify.Result{
		OK:             false,
		Ref:            refStr,
		HashAlg:        "sha256(canonical_json_v1)",
		CanonicalScope: "canonical_json_v1_excluding_expected_hash_v1",
		Trace:          []string{},
		Errors:         []string{},
	}

	result := newResultV1(base)

	// locate bundle root
	root, rootsTried, err := FindBundleRoot(refStr, fixtureRoot, op.DataDir, op.PreferData)
	result.Trace = append(result.Trace, rootsTried...)
	if err != nil {
		result.Message = err.Error()
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	bundle, trace, err := LoadBundleRootV1(root)
	result.Trace = append(result.Trace, trace...)
	if err != nil {
		result.Message = err.Error()
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	// read expected_hash_v1 from snapshot.json
	var meta snapshotMeta
	if err := json.Unmarshal(bundle.Snapshot, &meta); err != nil {
		result.Message = fmt.Sprintf("decode snapshot.json meta: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}
	if meta.ExpectedHashV1 == "" {
		result.Message = "snapshot.json missing expected_hash_v1"
		result.Errors = append(result.Errors, "missing expected_hash_v1")
		return result, nil
	}
	result.Expected = meta.ExpectedHashV1

	// remove expected_hash_v1 from snapshot content BEFORE hashing (prevents self-reference)
	var snapObj map[string]any
	if err := json.Unmarshal(bundle.Snapshot, &snapObj); err != nil {
		result.Message = fmt.Sprintf("decode snapshot.json object: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}
	delete(snapObj, "expected_hash_v1")

	cleanSnap, err := json.Marshal(snapObj)
	if err != nil {
		result.Message = fmt.Sprintf("re-encode snapshot.json: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}
	bundle.Snapshot = cleanSnap

	// compute hash (no changes to hash computation)
	state := StateV1FromBundle(bundle)
	hv1, err := snaps.HashV1FromState(state)
	if err != nil {
		result.Message = fmt.Sprintf("hash v1: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	computedHash := string(hv1)
	computedHashLower := strings.ToLower(computedHash)
	result.Got = computedHash

	if computedHash != result.Expected {
		result.Message = fmt.Sprintf("hash mismatch expected=%s got=%s", result.Expected, computedHash)
	} else {
		result.OK = true
		result.Message = "ok"
	}

	// --- write-expected flow ---
	snapshotPath := filepath.Join(root, "snapshot.json")
	expectedIsPlaceholder := IsPlaceholderExpected(result.Expected)

	if !op.Options.WriteExpected {
		// Keep default write_reason="none" unless a placeholder is present.
		// Only then report that the flag was not set.
		if expectedIsPlaceholder {
			// Ensure the computed hash looks like a valid digest before claiming write would be possible.
			if IsValidSHA256HexLower(computedHashLower) {
				result.WriteReason = WriteReasonFlagNotSet
			}
		}
		return result, nil
	}

	// WriteExpected enabled
	if !IsValidSHA256HexLower(computedHashLower) {
		result.WriteBlocked = true
		result.WriteReason = WriteReasonInvalidHash
		return result, fmt.Errorf("%w", bun.ErrInvalidNewHash)
	}

	if err := bun.WriteExpectedHashV1(snapshotPath, computedHash); err != nil {
		result.WriteBlocked = true
		switch {
		case errors.Is(err, bun.ErrExpectedAlreadySet):
			result.WriteReason = WriteReasonExistingExpected
		case errors.Is(err, bun.ErrSnapshotNotFound):
			result.WriteReason = WriteReasonSnapshotNotFound
		case errors.Is(err, bun.ErrSnapshotInvalidJSON):
			result.WriteReason = WriteReasonSnapshotInvalidJSON
		case errors.Is(err, bun.ErrInvalidNewHash):
			result.WriteReason = WriteReasonInvalidHash
		default:
			result.WriteReason = WriteReasonIOError
		}
		return result, err
	}

	result.WroteExpected = true
	result.WriteReason = WriteReasonPlaceholder
	return result, nil
}
