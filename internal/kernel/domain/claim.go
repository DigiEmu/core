package domain

import (
	"fmt"
	"sort"
)

const (
	ClaimSchemaV0    = "claim/v0"
	ClaimSetSchemaV0 = "claimset/v0"
)

type Claim struct {
	ID   string   `json:"id"`
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

type RelationType string

const (
	RelationContradicts RelationType = "CONTRADICTS"
)

type ClaimRelation struct {
	Type        RelationType `json:"type"`
	FromClaimID string       `json:"from_claim_id"`
	ToClaimID   string       `json:"to_claim_id"`
}

type ClaimSet struct {
	SchemaVersion string          `json:"schema_version"`
	VersionID     string          `json:"version_id"`
	Claims        []Claim         `json:"claims"`
	Relations     []ClaimRelation `json:"relations"`
}

// ValidateMinimal checks basic schema and referential integrity rules.
func (cs *ClaimSet) ValidateMinimal() error {
	if cs == nil {
		return fmt.Errorf("claimset is nil")
	}
	if cs.SchemaVersion != ClaimSetSchemaV0 {
		return fmt.Errorf("invalid schema_version: want %s got %s", ClaimSetSchemaV0, cs.SchemaVersion)
	}
	if cs.VersionID == "" {
		return fmt.Errorf("version_id is required")
	}

	idSeen := make(map[string]struct{}, len(cs.Claims))
	claimIDs := make([]string, 0, len(cs.Claims))

	for i := 0; i < len(cs.Claims); i++ {
		c := cs.Claims[i]
		if c.ID == "" {
			return fmt.Errorf("claim[%d]: id is required", i)
		}
		if c.Text == "" {
			return fmt.Errorf("claim[%d]: text is required for id=%s", i, c.ID)
		}
		if _, ok := idSeen[c.ID]; ok {
			return fmt.Errorf("duplicate claim id: %s", c.ID)
		}
		idSeen[c.ID] = struct{}{}
		claimIDs = append(claimIDs, c.ID)
	}

	sort.Strings(claimIDs)

	for i := 0; i < len(cs.Relations); i++ {
		r := cs.Relations[i]
		if r.Type != RelationContradicts {
			return fmt.Errorf("relation[%d]: unsupported relation type: %s", i, r.Type)
		}
		if r.FromClaimID == "" || r.ToClaimID == "" {
			return fmt.Errorf("relation[%d]: from_claim_id and to_claim_id are required", i)
		}
		if _, ok := idSeen[r.FromClaimID]; !ok {
			return fmt.Errorf("relation[%d]: from_claim_id references unknown claim id: %s", i, r.FromClaimID)
		}
		if _, ok := idSeen[r.ToClaimID]; !ok {
			return fmt.Errorf("relation[%d]: to_claim_id references unknown claim id: %s", i, r.ToClaimID)
		}
	}

	return nil
}
