package usecases

import (
	"fmt"
	"sort"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// VerifyAudit verifies that each Unit and Version has a corresponding audit event.
// It detects:
// - missing unit.created for units
// - missing version.created for versions
// - duplicates (multiple events for same unit/version)
// - optional content hash mismatch (StrictHash)
type VerifyAudit struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLogReader
}

func sortedStringKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (uc VerifyAudit) VerifyAudit(in ports.VerifyAuditRequest) (ports.VerifyAuditResponse, error) {
	if uc.Repo == nil {
		return ports.VerifyAuditResponse{}, fmt.Errorf("repo not configured")
	}
	if uc.Audit == nil {
		return ports.VerifyAuditResponse{}, fmt.Errorf("audit reader not configured")
	}

	var units []domain.Unit
	if in.UnitKey != "" {
		u, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
		if err != nil {
			return ports.VerifyAuditResponse{}, err
		}
		if !ok {
			return ports.VerifyAuditResponse{}, domain.ErrUnitNotFound
		}
		units = []domain.Unit{u}
	} else {
		all, err := uc.Repo.ListUnits()
		if err != nil {
			return ports.VerifyAuditResponse{}, err
		}
		units = all
	}

	// Build expectations
	expectedUnitCreated := make(map[string]struct{}, len(units)) // unitID -> exists
	expectedVersions := make(map[string]domain.Version)          // versionID -> version
	versionToUnit := make(map[string]string)                     // versionID -> unitID

	totalVersions := 0
	for _, u := range units {
		expectedUnitCreated[u.ID] = struct{}{}

		vs, err := uc.Repo.ListVersionsByUnitID(u.ID)
		if err != nil {
			return ports.VerifyAuditResponse{}, err
		}
		totalVersions += len(vs)
		for _, v := range vs {
			expectedVersions[v.ID] = v
			versionToUnit[v.ID] = u.ID
		}
	}

	// Track found events and duplicates
	foundUnitCreated := make(map[string]int)    // unitID -> count
	foundVersionCreated := make(map[string]int) // versionID -> count
	foundVersionHash := make(map[string]string) // versionID -> contentHash from event
	foundMeaningEvent := make(map[string]int)
	foundMeaningHash := make(map[string]string)
	foundClaimEvent := make(map[string]int)
	foundClaimHash := make(map[string]string)
	foundUncertaintyEvent := make(map[string]int)
	foundUncertaintyHash := make(map[string]string)

	// Scan audit log
	if err := uc.Audit.Scan(func(ev domain.AuditEvent) error {
		switch ev.Type {
		case "unit.created":
			if ev.UnitID != "" {
				if _, ok := expectedUnitCreated[ev.UnitID]; ok {
					foundUnitCreated[ev.UnitID]++
				}
			}
		case "version.created":
			if ev.VersionID != "" {
				if _, ok := expectedVersions[ev.VersionID]; ok {
					foundVersionCreated[ev.VersionID]++
					switch d := ev.Data.(type) {
					case map[string]any:
						if h, ok := d["contentHash"].(string); ok && h != "" {
							foundVersionHash[ev.VersionID] = h
						}
					case domain.VersionCreatedData:
						if d.ContentHash != "" {
							foundVersionHash[ev.VersionID] = d.ContentHash
						}
					}
				}
			}
		case "MEANING_SET":
			if ev.VersionID != "" {
				if _, ok := expectedVersions[ev.VersionID]; ok {
					foundMeaningEvent[ev.VersionID]++
					switch d := ev.Data.(type) {
					case map[string]any:
						if h, ok := d["meaning_hash"].(string); ok && h != "" {
							foundMeaningHash[ev.VersionID] = h
						}
					case domain.MeaningSetData:
						if d.MeaningHash != "" {
							foundMeaningHash[ev.VersionID] = d.MeaningHash
						}
					}
				}
			}
		case "CLAIM_SET":
			if ev.VersionID != "" {
				if _, ok := expectedVersions[ev.VersionID]; ok {
					foundClaimEvent[ev.VersionID]++
					switch d := ev.Data.(type) {
					case map[string]any:
						if h, ok := d["claimset_hash"].(string); ok && h != "" {
							foundClaimHash[ev.VersionID] = h
						}
					case domain.ClaimSetData:
						if d.ClaimSetHash != "" {
							foundClaimHash[ev.VersionID] = d.ClaimSetHash
						}
					}
				}
			}
		case "CLAIM_RELATION_SET":
			// presence is noted but handled later when claimset exists
		case "UNCERTAINTY_SET":
			if ev.VersionID != "" {
				if _, ok := expectedVersions[ev.VersionID]; ok {
					foundUncertaintyEvent[ev.VersionID]++
					switch d := ev.Data.(type) {
					case map[string]any:
						if h, ok := d["uncertainty_hash"].(string); ok && h != "" {
							foundUncertaintyHash[ev.VersionID] = h
						}
					case domain.UncertaintySetData:
						if d.UncertaintyHash != "" {
							foundUncertaintyHash[ev.VersionID] = d.UncertaintyHash
						}
					}
				}
			}
		}
		return nil
	}); err != nil {
		return ports.VerifyAuditResponse{}, err
	}

	out := ports.VerifyAuditResponse{
		TotalUnits:     len(units),
		TotalVersions:  totalVersions,
		Missing:        []ports.MissingAudit{},
		Duplicates:     []ports.DuplicateAudit{},
		HashMismatches: []ports.HashMismatch{},
	}

	// Missing or duplicate unit.created
	for _, unitID := range sortedStringKeys(expectedUnitCreated) {
		n := foundUnitCreated[unitID]
		if n == 0 {
			out.Missing = append(out.Missing, ports.MissingAudit{
				UnitID: unitID, VersionID: "", EventType: "unit.created",
			})
		} else if n > 1 {
			out.Duplicates = append(out.Duplicates, ports.DuplicateAudit{
				EventType: "unit.created", TargetID: unitID,
			})
		}
	}

	// Missing or duplicate version.created
	for _, verID := range sortedStringKeys(expectedVersions) {
		n := foundVersionCreated[verID]
		if n == 0 {
			out.Missing = append(out.Missing, ports.MissingAudit{
				UnitID: versionToUnit[verID], VersionID: verID, EventType: "version.created",
			})
		} else if n > 1 {
			out.Duplicates = append(out.Duplicates, ports.DuplicateAudit{
				EventType: "version.created", TargetID: verID,
			})
		}
	}

	// Missing, duplicate, or mismatched MEANING_SET events when repo versions expect a meaning
	for _, verID := range sortedStringKeys(expectedVersions) {
		v := expectedVersions[verID]
		if v.MeaningHash == "" {
			continue
		}
		n := foundMeaningEvent[verID]
		if n == 0 {
			out.Missing = append(out.Missing, ports.MissingAudit{
				UnitID: versionToUnit[verID], VersionID: verID, EventType: "MEANING_SET",
			})
		} else if n > 1 {
			out.Duplicates = append(out.Duplicates, ports.DuplicateAudit{
				EventType: "MEANING_SET", TargetID: verID,
			})
		}
		eh, ok := foundMeaningHash[verID]
		if !ok || eh == "" || eh != v.MeaningHash {
			out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
				UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.MeaningHash, EventHash: eh,
			})
		}
		if in.StrictHash {
			m, ok, err := uc.Repo.LoadMeaning(versionToUnit[verID], verID)
			if err != nil {
				return ports.VerifyAuditResponse{}, err
			}
			if !ok {
				out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
					UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.MeaningHash, EventHash: "(missing sidecar)",
				})
			} else {
				mh, err := ComputeMeaningHash(m)
				if err != nil || mh != v.MeaningHash {
					out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
						UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.MeaningHash, EventHash: mh,
					})
				}
			}
		}
	}

	// Missing, duplicate, or mismatched CLAIM_SET events when repo versions expect a claimset
	for _, verID := range sortedStringKeys(expectedVersions) {
		v := expectedVersions[verID]
		if v.ClaimSetHash == "" {
			continue
		}
		n := foundClaimEvent[verID]
		if n == 0 {
			out.Missing = append(out.Missing, ports.MissingAudit{
				UnitID: versionToUnit[verID], VersionID: verID, EventType: "CLAIM_SET",
			})
		} else if n > 1 {
			out.Duplicates = append(out.Duplicates, ports.DuplicateAudit{EventType: "CLAIM_SET", TargetID: verID})
		}
		eh, ok := foundClaimHash[verID]
		if !ok || eh == "" || eh != v.ClaimSetHash {
			out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
				UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.ClaimSetHash, EventHash: eh,
			})
		}
		if in.StrictHash {
			cs, ok, err := uc.Repo.LoadClaimSet(versionToUnit[verID], verID)
			if err != nil {
				return ports.VerifyAuditResponse{}, err
			}
			if !ok {
				out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
					UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.ClaimSetHash, EventHash: "(missing sidecar)",
				})
			} else {
				ch, err := ComputeClaimSetHashFromStruct(cs)
				if err != nil || ch != v.ClaimSetHash {
					out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
						UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.ClaimSetHash, EventHash: ch,
					})
				}
			}
		}
	}

	// Missing, duplicate, or mismatched UNCERTAINTY_SET events when repo versions expect an uncertainty
	for _, verID := range sortedStringKeys(expectedVersions) {
		v := expectedVersions[verID]
		if v.UncertaintyHash == "" {
			continue
		}
		n := foundUncertaintyEvent[verID]
		if n == 0 {
			out.Missing = append(out.Missing, ports.MissingAudit{
				UnitID: versionToUnit[verID], VersionID: verID, EventType: "UNCERTAINTY_SET",
			})
		} else if n > 1 {
			out.Duplicates = append(out.Duplicates, ports.DuplicateAudit{EventType: "UNCERTAINTY_SET", TargetID: verID})
		}
		eh, ok := foundUncertaintyHash[verID]
		if !ok || eh == "" || eh != v.UncertaintyHash {
			out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
				UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.UncertaintyHash, EventHash: eh,
			})
		}
		if in.StrictHash {
			u, ok, err := uc.Repo.LoadUncertainty(versionToUnit[verID], verID)
			if err != nil {
				return ports.VerifyAuditResponse{}, err
			}
			if !ok {
				out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
					UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.UncertaintyHash, EventHash: "(missing sidecar)",
				})
			} else {
				uHash, err := ComputeUncertaintyHashFromStruct(u)
				if err != nil || uHash != v.UncertaintyHash {
					out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
						UnitID: versionToUnit[verID], VersionID: verID, ExpectedHash: v.UncertaintyHash, EventHash: uHash,
					})
				}
			}
		}
	}

	// Optional strict hash verification
	if in.StrictHash {
		for _, verID := range sortedStringKeys(expectedVersions) {
			v := expectedVersions[verID]
			eh, ok := foundVersionHash[verID]
			if foundVersionCreated[verID] > 0 {
				if !ok || eh == "" || eh != v.ContentHash {
					out.HashMismatches = append(out.HashMismatches, ports.HashMismatch{
						UnitID:       versionToUnit[verID],
						VersionID:    verID,
						ExpectedHash: v.ContentHash,
						EventHash:    eh,
					})
				}
			}
		}
	}

	out.Ok = len(out.Missing) == 0 && len(out.Duplicates) == 0 && len(out.HashMismatches) == 0
	return out, nil
}
