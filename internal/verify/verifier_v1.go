package verify

import (
	"fmt"

	snaps "digiemu-core/pkg/snapshot"
	pkgverify "digiemu-core/pkg/verify"
)

// VerifierV1 implements pkg/verify.Verifier for MVP snapshot bundles.
// It verifies expected_hash_v1 against HashV1FromState(bundle.state).
type VerifierV1 struct {
	DataDir     string
	FixtureRoot string
	PreferData  bool
}

func (v *VerifierV1) Verify(ref snaps.Ref) (pkgverify.Result, error) {
	refStr := string(ref.Hash)

	fixtureRoot := v.FixtureRoot
	if fixtureRoot == "" {
		fixtureRoot = "data/test-fixtures"
	}

	sb, chosenPath, attempts, err := FindBundleV1(fixtureRoot, v.DataDir, refStr, v.PreferData)
	result := pkgverify.Result{
		OK:             false,
		Ref:            refStr,
		HashAlg:        "sha256(canonical_json_v1)",
		CanonicalScope: "canonical_utf8_without_sha256_comment_line",
		Errors:         []string{},
	}
	if len(attempts) > 0 {
		// include attempted paths in errors if failure
		for _, p := range attempts {
			result.Errors = append(result.Errors, p)
		}
	}

	if err != nil {
		result.Message = err.Error()
		return result, nil
	}

	// compute hash
	hv1, err := snaps.HashV1FromState(sb.State)
	if err != nil {
		result.Message = fmt.Sprintf("hash v1: %v", err)
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}

	got := string(hv1)
	exp := sb.ExpectedHashV1

	result.Expected = exp
	result.Got = got
	// chosenPath may be useful for debugging; add to errors slice as informational entry
	if chosenPath != "" {
		result.Errors = append(result.Errors, fmt.Sprintf("used:%s", chosenPath))
	}

	if got != exp {
		result.Message = fmt.Sprintf("hash mismatch expected=%s got=%s", exp, got)
		return result, nil
	}

	result.OK = true
	result.Message = "ok"
	return result, nil
}
