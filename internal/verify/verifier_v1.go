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
		CanonicalScope: "canonical_json_v1",
		Trace:          []string{},
		Errors:         []string{},
	}
	if len(attempts) > 0 {
		// include attempted paths as trace information
		for _, p := range attempts {
			result.Trace = append(result.Trace, p)
		}
	}

	if err != nil {
		result.Message = err.Error()
		result.Errors = append(result.Errors, err.Error())
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
	// chosenPath may be useful for debugging; add to trace
	if chosenPath != "" {
		result.Trace = append(result.Trace, fmt.Sprintf("used:%s", chosenPath))
	}

	if got != exp {
		result.Message = fmt.Sprintf("hash mismatch expected=%s got=%s", exp, got)
		return result, nil
	}

	result.OK = true
	result.Message = "ok"
	return result, nil
}
