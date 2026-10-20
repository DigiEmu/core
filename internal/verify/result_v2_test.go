package verify

import (
	"errors"
	"fmt"
	"testing"

	pkgverify "digiemu-core/pkg/verify"
	derr "digiemu-core/pkg/digiemu"
)

func TestParseVerifyResultVersion(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"", "v1", false},
		{"v1", "v1", false},
		{"v2", VerifyResultV2Version, false},
		{"v2-draft", VerifyResultV2Version, false},
		{"2.0-draft", VerifyResultV2Version, false},
		{"V2-DRAFT", VerifyResultV2Version, false},
		{"v3", "", true},
		{"bogus", "", true},
	}
	for _, tc := range cases {
		got, err := ParseVerifyResultVersion(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("input %q: expected error, got %q", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("input %q: got %q, want %q", tc.in, got, tc.want)
		}
	}
}

func okResult() ResultV1 {
	return ResultV1{Result: pkgverify.Result{
		OK:       true,
		Ref:      "demo",
		Expected: "A1B2C3D4E5F6A7B8C9D0E1F2A3B4C5D6E7F8A9B0C1D2E3F4A5B6C7D8E9F0A1B2",
		Got:      "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
	}}
}

func TestResultV2_PASS_HasSnapshotHashAndProfiles(t *testing.T) {
	res := okResult()
	v2 := ResultV2FromResultV1(res, nil, "test-version")
	if v2.Result != "PASS" {
		t.Fatalf("result=%q, want PASS", v2.Result)
	}
	if v2.ReasonCode != "STATE_RECONSTRUCTED" {
		t.Fatalf("reason_code=%q", v2.ReasonCode)
	}
	if v2.VerifyResultVersion != VerifyResultV2Version {
		t.Fatalf("verify_result_version=%q", v2.VerifyResultVersion)
	}
	if v2.SnapshotHash == "" {
		t.Fatalf("PASS result missing snapshot_hash")
	}
	if v2.SnapshotHash != "sha256:a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2" {
		t.Fatalf("unexpected snapshot_hash=%q", v2.SnapshotHash)
	}
	if v2.CanonicalizationProfile != "digiemu-canonical-json-v1" ||
		v2.HashAlgorithm != "sha256" ||
		v2.ReconstructionProfile != "digiemu-reconstruct-v1" {
		t.Fatalf("unexpected profile fields: %+v", v2)
	}
	if v2.VerifierVersion != "test-version" {
		t.Fatalf("verifier_version=%q", v2.VerifierVersion)
	}
}

func TestResultV2_FAIL_HashMismatchCarriesExpectedActual(t *testing.T) {
	res := okResult()
	res.OK = false
	res.Expected = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	res.Got = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	res.Message = "hash mismatch"

	v2 := ResultV2FromResultV1(res, nil, "x")
	if v2.Result != "FAIL" || v2.ReasonCode != "HASH_MISMATCH" {
		t.Fatalf("got %s/%s", v2.Result, v2.ReasonCode)
	}
	if v2.SnapshotHashExpected != "sha256:"+res.Expected {
		t.Fatalf("snapshot_hash_expected=%q", v2.SnapshotHashExpected)
	}
	if v2.SnapshotHashActual != "sha256:"+res.Got {
		t.Fatalf("snapshot_hash_actual=%q", v2.SnapshotHashActual)
	}
	if v2.SnapshotHash != "" {
		t.Fatalf("FAIL result should not carry snapshot_hash, got %q", v2.SnapshotHash)
	}
}

func TestResultV2_FAIL_PlaceholderExpectedNotDisguised(t *testing.T) {
	res := okResult()
	res.OK = false
	res.Expected = "REPLACE_AFTER_COMPUTE"
	res.Got = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	res.Message = "hash mismatch"

	v2 := ResultV2FromResultV1(res, nil, "x")
	if v2.Result != "FAIL" {
		t.Fatalf("got %s", v2.Result)
	}
	if v2.SnapshotHashExpected != "REPLACE_AFTER_COMPUTE" {
		t.Fatalf("placeholder expected should be reported raw, got %q", v2.SnapshotHashExpected)
	}
}

func TestResultV2_ERROR_OpErrMapsToReasonCode(t *testing.T) {
	res := okResult()
	res.OK = false
	res.Got = ""
	res.Expected = ""

	v2 := ResultV2FromResultV1(res, fmt.Errorf("boom"), "x")
	if v2.Result != "ERROR" || v2.ReasonCode != "INTERNAL_ERROR" {
		t.Fatalf("got %s/%s", v2.Result, v2.ReasonCode)
	}
}

func TestResultV2_ERROR_MissingExpectedHashMapsReason(t *testing.T) {
	res := okResult()
	res.OK = false
	res.Got = ""
	res.Errors = []string{"missing expected_hash_v1"}
	res.Message = "snapshot.json missing expected_hash_v1"

	v2 := ResultV2FromResultV1(res, nil, "x")
	if v2.Result != "ERROR" || v2.ReasonCode != "MISSING_REQUIRED_FIELD" {
		t.Fatalf("got %s/%s", v2.Result, v2.ReasonCode)
	}
	if v2.FailurePath == "" {
		t.Fatalf("expected failure_path for missing expected_hash_v1")
	}
}

func TestResultV2_ERROR_DecodeFailureMapsSchemaInvalid(t *testing.T) {
	res := okResult()
	res.OK = false
	res.Got = ""
	res.Errors = []string{"decode snapshot.json meta: invalid character"}
	res.Message = "decode snapshot.json meta: invalid character"

	v2 := ResultV2FromResultV1(res, nil, "x")
	if v2.Result != "ERROR" || v2.ReasonCode != "INVALID_SNAPSHOT_SCHEMA" {
		t.Fatalf("got %s/%s", v2.Result, v2.ReasonCode)
	}
}

// TestResultV2_BlockedExistingExpectedPreservesCurrentErrorProjection
// characterizes the current lossy v2 mapping without freezing new governance
// semantics. The CLI exit code and full v1 write_reason remain separate.
func TestResultV2_BlockedExistingExpectedPreservesCurrentErrorProjection(t *testing.T) {
	res := okResult()
	res.WriteBlocked = true
	res.WriteReason = WriteReasonExistingExpected
	err := derr.UsageError("verify.writeExpected", "snapshot.json", errors.New("expected already set"))
	v2 := ResultV2FromResultV1(res, err, "test")
	if v2.Result != "ERROR" || v2.ReasonCode != "INTERNAL_ERROR" {
		t.Fatalf("blocked write projects to %s/%s; want ERROR/INTERNAL_ERROR until policy contract is decided", v2.Result, v2.ReasonCode)
	}
	if v2.SnapshotHash != "" {
		t.Fatalf("ERROR should not manufacture snapshot_hash: %q", v2.SnapshotHash)
	}
}

func TestResultV2_RequiredFieldsNeverEmpty(t *testing.T) {
	// Even a zero-value v1 result must produce the schema-required fields.
	var res ResultV1
	v2 := ResultV2FromResultV1(res, errors.New("x"), "")
	for k, v := range map[string]string{
		"verify_result_version":    v2.VerifyResultVersion,
		"result":                   v2.Result,
		"reason_code":              v2.ReasonCode,
		"canonicalization_profile": v2.CanonicalizationProfile,
		"hash_algorithm":           v2.HashAlgorithm,
		"reconstruction_profile":   v2.ReconstructionProfile,
	} {
		if v == "" {
			t.Fatalf("required field %s is empty", k)
		}
	}
}
