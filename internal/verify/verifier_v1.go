package verify

import (
	"encoding/json"
	"fmt"

	snaps "digiemu-core/pkg/snapshot"
	pkgverify "digiemu-core/pkg/verify"
)

// VerifierV1 implements pkg/verify.Verifier for assembled snapshot bundles.
// It verifies expected_hash_v1 against HashV1FromState(assembled_state),
// where expected_hash_v1 is EXCLUDED from the hashed state (otherwise self-referential).
type VerifierV1 struct {
	DataDir     string
	FixtureRoot string
	PreferData  bool
}

type snapshotMeta struct {
	ExpectedHashV1 string `json:"expected_hash_v1"`
}

func (v *VerifierV1) Verify(ref snaps.Ref) (pkgverify.Result, error) {
	refStr := string(ref.Hash)

	fixtureRoot := v.FixtureRoot
	if fixtureRoot == "" {
		fixtureRoot = "data/test-fixtures"
	}

	result := pkgverify.Result{
		OK:             false,
		Ref:            refStr,
		HashAlg:        "sha256(canonical_json_v1)",
		CanonicalScope: "canonical_json_v1_excluding_expected_hash_v1",
		Trace:          []string{},
		Errors:         []string{},
	}

	// locate bundle root
	root, rootsTried, err := FindBundleRoot(refStr, fixtureRoot, v.DataDir, v.PreferData)
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

	// 1) Read expected_hash_v1 from snapshot.json (BOM-safe already)
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

	// 2) Remove expected_hash_v1 from snapshot content BEFORE hashing (prevents self-reference)
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

	// 3) Assemble StateV1 and compute hash over canonical_json_v1
	state := StateV1FromBundle(bundle)
	hv1, err := snaps.HashV1FromState(state)
	if err != nil {
		result.Message = fmt.Sprintf("hash v1: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	got := string(hv1)
	result.Got = got

	if got != result.Expected {
		result.Message = fmt.Sprintf("hash mismatch expected=%s got=%s", result.Expected, got)
		return result, nil
	}

	result.OK = true
	result.Message = "ok"
	return result, nil
}
