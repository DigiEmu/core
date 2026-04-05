package domain

// Uncertainty schema v0 - minimal, auditable uncertainty metadata
const UncertaintySchemaV0 = "uncertainty/v0"

type AppliesToScope string

const (
	ScopeVersion AppliesToScope = "version"
	ScopeClaim   AppliesToScope = "claim"
)

type Uncertainty struct {
	SchemaVersion string               `json:"schema_version"`
	ID            string               `json:"id"`
	Type          string               `json:"type"`
	Level         string               `json:"level"`
	Text          string               `json:"text"`
	Tags          []string             `json:"tags"`
	AppliesTo     UncertaintyAppliesTo `json:"applies_to"`
}

type UncertaintyAppliesTo struct {
	Scope   AppliesToScope `json:"scope"`
	ClaimID string         `json:"claim_id"`
}

// ValidateMinimal enforces the minimal invariants described in the spec.
func (u Uncertainty) ValidateMinimal() error {
	if u.SchemaVersion != UncertaintySchemaV0 {
		return ErrInvalidSchemaVersion
	}
	if u.ID == "" {
		return ErrInvalidUncertaintyID
	}
	switch u.Type {
	case "empirical", "interpretative", "incomplete":
		// ok
	default:
		return ErrInvalidUncertaintyType
	}
	switch u.Level {
	case "low", "medium", "high":
		// ok
	default:
		return ErrInvalidUncertaintyLevel
	}
	switch u.AppliesTo.Scope {
	case ScopeVersion:
		// ok
	case ScopeClaim:
		if u.AppliesTo.ClaimID == "" {
			return ErrMissingClaimIDForUncertainty
		}
	default:
		return ErrInvalidAppliesToScope
	}
	return nil
}
