package domain

// Meaning represents structured context attached to a Unit/Version.
// For deterministic export, replay and hashing, fields are serialized explicitly.
// Empty values remain visible instead of silently disappearing.
type Meaning struct {
	SchemaVersion string             `json:"schema_version"`
	Title         string             `json:"title"`
	Purpose       string             `json:"purpose"`
	Scope         *MeaningScope      `json:"scope"`
	Claims        []MeaningClaim     `json:"claims"`
	Sources       []MeaningSource    `json:"sources"`
	Provenance    *MeaningProvenance `json:"provenance"`
	Integrity     *MeaningIntegrity  `json:"integrity"`
}

type MeaningScope struct {
	Audience     []string          `json:"audience"`
	Jurisdiction []string          `json:"jurisdiction"`
	Locale       []string          `json:"locale"`
	Timeframe    *MeaningTimeframe `json:"timeframe"`
}

type MeaningTimeframe struct {
	ValidFrom  string `json:"valid_from"`
	ValidUntil string `json:"valid_until"`
}

type MeaningClaim struct {
	Text     string   `json:"text"`
	Strength string   `json:"strength"`
	Tags     []string `json:"tags"`
}

type MeaningSource struct {
	ID    string              `json:"id"`
	Type  string              `json:"type"`
	Ref   string              `json:"ref"`
	Quote *MeaningSourceQuote `json:"quote"`
}

type MeaningSourceQuote struct {
	Snippet string `json:"snippet"`
	Locator string `json:"locator"`
}

type MeaningProvenance struct {
	Author    string `json:"author"`
	Org       string `json:"org"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MeaningIntegrity struct {
	NarrativeID   string   `json:"narrative_id"`
	Supersedes    string   `json:"supersedes"`
	ConflictsWith []string `json:"conflicts_with"`
}

// MeaningRef is used in export manifests to reference a meaning document.
type MeaningRef struct {
	MeaningHash   string `json:"meaning_hash"`
	MeaningPath   string `json:"meaning_path"`
	SchemaVersion string `json:"schema_version"`
}

// MeaningHash is the hex-encoded sha256 of the canonicalized meaning.json
type MeaningHash string
