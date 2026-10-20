package verify

import (
	"strings"

	derr "digiemu-core/pkg/digiemu"
)

// Canonical Verify Result version identifier emitted in v2 output.
// "2.0-draft" is a compatibility alias accepted on input; the canonical
// emitted identifier is "v2-draft".
const VerifyResultV2Version = "v2-draft"

// ResultV2 is the opt-in draft v2 verify result contract.
//
// Field names follow schemas/verify_result_v2.schema.json. All fields listed
// in the schema's required array are always emitted; optional fields are
// emitted only when a corresponding v1 value exists.
type ResultV2 struct {
	VerifyResultVersion       string   `json:"verify_result_version"`
	Result                    string   `json:"result"`
	ReasonCode                string   `json:"reason_code"`
	CanonicalizationProfile   string   `json:"canonicalization_profile"`
	HashAlgorithm             string   `json:"hash_algorithm"`
	ReconstructionProfile     string   `json:"reconstruction_profile"`
	VerifierVersion           string   `json:"verifier_version"`
	SnapshotHash              string   `json:"snapshot_hash,omitempty"`
	SnapshotHashExpected      string   `json:"snapshot_hash_expected,omitempty"`
	SnapshotHashActual        string   `json:"snapshot_hash_actual,omitempty"`
	FailurePath               string   `json:"failure_path,omitempty"`
	EvidenceReferences        []string `json:"evidence_references,omitempty"`
	SignatureStatus           string   `json:"signature_status,omitempty"`
	CryptoProfile             string   `json:"crypto_profile,omitempty"`
	OutsideHashMetadataStatus string   `json:"outside_hash_metadata_status,omitempty"`
}

// ParseVerifyResultVersion normalizes a --verify-result-version flag value
// to a canonical version selector. Empty or "v1" selects the default v1
// contract; "v2", "v2-draft", and the compatibility alias "2.0-draft" select
// the opt-in v2-draft contract.
func ParseVerifyResultVersion(s string) (string, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "", "v1":
		return "v1", nil
	case "v2", "v2-draft", "2.0-draft":
		return VerifyResultV2Version, nil
	default:
		return "", derr.UsageError("verify.resultVersion", "--verify-result-version", nil)
	}
}

// isSHA256Hex reports whether s is a 64-character lowercase/uppercase hex digest.
func isSHA256Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

// taggedHash returns "sha256:<hex>" for 64-hex digests and the raw value
// otherwise (placeholders are reported honestly, not disguised as digests).
func taggedHash(s string) string {
	if isSHA256Hex(s) {
		return "sha256:" + strings.ToLower(s)
	}
	return s
}

// reasonCodeForError maps an operational verify error to a v2 reason code.
func reasonCodeForError(opErr error) string {
	code, ok := derr.CodeOf(opErr)
	if !ok {
		return "INTERNAL_ERROR"
	}
	switch code {
	case derr.CodeSchemaInvalid, derr.CodeBundleInvalid:
		return "INVALID_SNAPSHOT_SCHEMA"
	case derr.CodeFileMissing:
		return "MISSING_REQUIRED_FIELD"
	default:
		return "INTERNAL_ERROR"
	}
}

// reasonCodeForResultErrors maps a ResultV1 error trail to a v2 reason code.
func reasonCodeForResultErrors(res ResultV1) (code string, failurePath string) {
	msg := res.Message + "\n" + strings.Join(res.Errors, "\n")
	switch {
	case strings.Contains(msg, "missing expected_hash_v1"):
		return "MISSING_REQUIRED_FIELD", "/expected_hash_v1"
	case strings.Contains(msg, "decode snapshot.json"):
		return "INVALID_SNAPSHOT_SCHEMA", "/snapshot"
	default:
		return "INTERNAL_ERROR", ""
	}
}

// ResultV2FromResultV1 projects a v1 ResultV1 (plus any operational error)
// onto the draft v2 contract. The projection is lossy in one direction only:
// v1 trace/errors/message fields have no v2 counterpart and are dropped.
// Hash state is never fabricated: hash fields are emitted only when the
// corresponding v1 value is present.
func ResultV2FromResultV1(res ResultV1, opErr error, verifierVersion string) ResultV2 {
	v2 := ResultV2{
		VerifyResultVersion:     VerifyResultV2Version,
		CanonicalizationProfile: "digiemu-canonical-json-v1",
		HashAlgorithm:           "sha256",
		ReconstructionProfile:   "digiemu-reconstruct-v1",
		VerifierVersion:         verifierVersion,
		SignatureStatus:         "none",
	}

	switch {
	case opErr != nil:
		v2.Result = "ERROR"
		v2.ReasonCode = reasonCodeForError(opErr)
	case res.OK:
		v2.Result = "PASS"
		v2.ReasonCode = "STATE_RECONSTRUCTED"
		v2.SnapshotHash = taggedHash(res.Got)
	case res.Got != "" && res.Expected != "" && len(res.Errors) == 0:
		v2.Result = "FAIL"
		v2.ReasonCode = "HASH_MISMATCH"
		v2.SnapshotHashExpected = taggedHash(res.Expected)
		v2.SnapshotHashActual = taggedHash(res.Got)
		v2.FailurePath = "/expected_hash_v1"
	default:
		v2.Result = "ERROR"
		v2.ReasonCode, v2.FailurePath = reasonCodeForResultErrors(res)
	}

	return v2
}
