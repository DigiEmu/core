package verify

import pkgverify "digiemu-core/pkg/verify"

// WriteReason describes why a write-expected operation happened or was blocked.
//
// String values are part of the ResultV1 JSON contract.
type WriteReason string

const (
	WriteReasonNone                WriteReason = "none"
	WriteReasonFlagNotSet          WriteReason = "flag_not_set"
	WriteReasonPlaceholder         WriteReason = "placeholder"
	WriteReasonExistingExpected    WriteReason = "existing_expected_present"
	WriteReasonInvalidHash         WriteReason = "invalid_hash"
	WriteReasonSnapshotNotFound    WriteReason = "snapshot_not_found"
	WriteReasonSnapshotInvalidJSON WriteReason = "snapshot_invalid_json"
	WriteReasonIOError             WriteReason = "io_error"
)

// ResultV1 is the internal verification result contract for the write-expected workflow.
//
// It inlines the stable public pkg/verify.Result fields and adds write-policy fields.
// Existing fields are preserved without breaking changes.
type ResultV1 struct {
	pkgverify.Result

	WroteExpected bool        `json:"wrote_expected"`
	WriteBlocked  bool        `json:"write_blocked"`
	WriteReason   WriteReason `json:"write_reason"`
}
