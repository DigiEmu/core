package domain

// AuditEvent is append-only journal event (NDJSON friendly).
type AuditEvent struct {
	Schema  string `json:"schema"`
	ID      string `json:"id"`
	Type    string `json:"type"`
	AtUnix  int64  `json:"atUnix"`
	ActorID string `json:"actorId"`

	UnitID    string `json:"unitId"`
	VersionID string `json:"versionId"`

	Data any `json:"data"`
}

type UnitCreatedData struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type VersionCreatedData struct {
	PrevVersionID string `json:"prevVersionId"`
	ContentHash   string `json:"contentHash"`
	Label         string `json:"label"`
}

type MeaningInlinePreview struct {
	Title   string `json:"title"`
	Purpose string `json:"purpose"`
}

type MeaningSetData struct {
	MeaningHash   string                `json:"meaning_hash"`
	MeaningPath   string                `json:"meaning_path"`
	SchemaVersion string                `json:"schema_version"`
	InlinePreview *MeaningInlinePreview `json:"inline_preview"`
}

type ClaimSetData struct {
	UnitID       string `json:"unit_id"`
	VersionID    string `json:"version_id"`
	ClaimSetHash string `json:"claimset_hash"`
	ClaimSetPath string `json:"claimset_path"`
}

type ClaimRelationSetData struct {
	UnitID      string `json:"unit_id"`
	VersionID   string `json:"version_id"`
	Type        string `json:"type"`
	FromClaimID string `json:"from_claim_id"`
	ToClaimID   string `json:"to_claim_id"`
}

type UncertaintySetData struct {
	UnitID          string `json:"unit_id"`
	VersionID       string `json:"version_id"`
	UncertaintyHash string `json:"uncertainty_hash"`
	UncertaintyPath string `json:"uncertainty_path"`
}
