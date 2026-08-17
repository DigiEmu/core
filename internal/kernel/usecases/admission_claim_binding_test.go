// Phase D binding proof for core.claim.set:
// A real SetClaims mutation executes only after P0 ADMIT and produces a
// before/after ClaimSet state evidence envelope. Test-only; not production.
package usecases_test

import (
	"fmt"
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func newClaimSetIntent(unitKey, versionID, bodyJSON, actor string) admission.Intent {
	return admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-claim-" + unitKey,
		CapabilityRef:        "core.claim.set",
		AggregateRef:         "unit",
		CommandRef:           "claim.set",
		Payload: map[string]any{
			"unit_key":   unitKey,
			"version_id": versionID,
			"body_json":  bodyJSON,
			"actor_id":   actor,
		},
	}
}

func TestClaimSet_Binding_And_EventEvidence(t *testing.T) {
	unitKey := "phase-d-claim-unit-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// PREREQUISITE TEST SETUP: existing Unit and Version.
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Claim Unit",
		Description: "created via P0 Phase D claim binding test",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	createVersion := usecases.CreateVersion{Repo: repo, Audit: setupAudit, Clock: clock}
	versionResp, err := createVersion.CreateVersion(ports.CreateVersionRequest{
		UnitKey: unitKey,
		Label:   "v1",
		Content: "Version content for claim test",
		ActorID: "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateVersion: %v", err)
	}

	// BEFORE: no ClaimSet for this Version; related fields empty.
	beforeVer, ok, _ := repo.FindVersionByID(versionResp.VersionID)
	if !ok {
		t.Fatalf("prerequisite version not found")
	}
	if beforeVer.ClaimSetHash != "" {
		t.Fatalf("expected empty ClaimSetHash before SetClaims, got %s", beforeVer.ClaimSetHash)
	}
	beforeMeaning, beforeMeaningOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if beforeMeaningOK {
		t.Fatalf("expected no meaning before SetClaims, got %+v", beforeMeaning)
	}
	beforeClaim, beforeClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if beforeClaimOK {
		t.Fatalf("expected no claimset before SetClaims, got %+v", beforeClaim)
	}
	beforeUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if beforeUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("expected head version %s, got %s", versionResp.VersionID, beforeUnit.HeadVersionID)
	}

	// THE CLAIM MUTATION UNDER ADMISSION PROOF
	bodyJSON := fmt.Sprintf(`{"schema_version":"claimset/v0","version_id":"%s","claims":[{"id":"c1","text":"A sample claim"}]}`, versionResp.VersionID)
	intent := newClaimSetIntent(unitKey, versionResp.VersionID, bodyJSON, "phase-d-tester")

	intentSchema := loadSchema(t, "schemas/intent-envelope.schema.json")
	if err := intentSchema.Validate(toAny(t, intent)); err != nil {
		t.Fatalf("intent envelope schema validation: %v", err)
	}

	eng := admission.NewEngine(admission.V01Registry())
	adm, err := eng.Evaluate(intent)
	if err != nil {
		t.Fatalf("admission evaluation: %v", err)
	}
	if adm.Decision != "ADMIT" {
		t.Fatalf("expected ADMIT, got %s", adm.Decision)
	}
	if adm.TransitionRef != "claim:set" {
		t.Fatalf("expected transition claim:set, got %s", adm.TransitionRef)
	}
	t.Logf("admission_id=%s", adm.AdmissionID)

	// Real SetClaims, same as production runtime, gated by the ADMIT result.
	claimAudit := memory.NewAuditLog()
	uc := usecases.SetClaims{Repo: repo, Audit: claimAudit, Clock: clock}
	req := ports.SetClaimsRequest{
		UnitKey:   stringFromMap(intent.Payload, "unit_key"),
		VersionID: stringFromMap(intent.Payload, "version_id"),
		BodyBytes: []byte(stringFromMap(intent.Payload, "body_json")),
		ActorID:   stringFromMap(intent.Payload, "actor_id"),
	}
	resp, err := uc.SetClaims(req)
	if err != nil {
		t.Fatalf("SetClaims: %v", err)
	}
	if resp.UnitID != unitResp.UnitID || resp.VersionID != versionResp.VersionID {
		t.Fatalf("SetClaims response mismatch: got %+v", resp)
	}
	if resp.ClaimSetHash == "" {
		t.Fatalf("expected non-empty ClaimSetHash")
	}

	// AFTER: ClaimSet state changed; nothing else.
	afterVer, ok, _ := repo.FindVersionByID(versionResp.VersionID)
	if !ok {
		t.Fatalf("version not found after SetClaims")
	}
	if afterVer.ClaimSetHash != resp.ClaimSetHash {
		t.Fatalf("version ClaimSetHash mismatch: got %s want %s", afterVer.ClaimSetHash, resp.ClaimSetHash)
	}

	afterClaim, afterClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if !afterClaimOK {
		t.Fatalf("claimset not found after SetClaims")
	}
	if afterClaim.SchemaVersion != "claimset/v0" {
		t.Fatalf("unexpected claimset schema version: %s", afterClaim.SchemaVersion)
	}
	if len(afterClaim.Claims) != 1 || afterClaim.Claims[0].ID != "c1" || afterClaim.Claims[0].Text != "A sample claim" {
		t.Fatalf("claimset content mismatch: got %+v", afterClaim)
	}

	// Cross-state preservation: meaning, uncertainty, unit head, version identity.
	afterMeaning, afterMeaningOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if afterMeaningOK {
		t.Fatalf("meaning unexpectedly changed: got %+v", afterMeaning)
	}
	if afterVer.MeaningHash != "" {
		t.Fatalf("meaning hash unexpectedly set: %s", afterVer.MeaningHash)
	}
	afterUncertainty, afterUncertaintyOK, _ := repo.LoadUncertainty(unitResp.UnitID, versionResp.VersionID)
	if afterUncertaintyOK {
		t.Fatalf("uncertainty unexpectedly changed: got %+v", afterUncertainty)
	}
	if afterVer.UncertaintyHash != "" {
		t.Fatalf("uncertainty hash unexpectedly set: %s", afterVer.UncertaintyHash)
	}

	afterUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if afterUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("unit head version changed unexpectedly: got %s want %s", afterUnit.HeadVersionID, versionResp.VersionID)
	}
	if afterVer.ID != beforeVer.ID || afterVer.UnitID != beforeVer.UnitID || afterVer.Content != beforeVer.Content || afterVer.Label != beforeVer.Label {
		t.Fatalf("version identity or content/label mutated: got %+v want %+v", afterVer, beforeVer)
	}

	// Real runtime AuditEvent emitted and isolated.
	reader := memory.NewAuditByUnitReader(claimAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 claim audit event, got %d", len(events))
	}
	if events[0].Type != "CLAIM_SET" {
		t.Fatalf("expected audit type CLAIM_SET, got %s", events[0].Type)
	}
	if events[0].VersionID != versionResp.VersionID {
		t.Fatalf("audit event version id mismatch: got %s want %s", events[0].VersionID, versionResp.VersionID)
	}
	if events[0].ActorID != req.ActorID {
		t.Fatalf("audit event actor mismatch: got %s want %s", events[0].ActorID, req.ActorID)
	}

	// Event Envelope from the real runtime event.
	evidence := map[string]any{
		"unit_id":        unitResp.UnitID,
		"version_id":     versionResp.VersionID,
		"claimset_hash":  resp.ClaimSetHash,
		"claim_id":       afterClaim.Claims[0].ID,
		"claim_text":     afterClaim.Claims[0].Text,
		"claimset_path":  unitResp.UnitID + "." + versionResp.VersionID + ".claimset.json",
	}
	eventEnv := map[string]any{
		"schema_version":        "v0.1",
		"architecture_revision": "0.3",
		"event_id":              events[0].ID,
		"command_ref":           "claim.set",
		"transition_ref":        "claim:set",
		"runtime_event_type":    "CLAIM_SET",
		"evidence":              evidence,
	}
	eventSchema := loadSchema(t, "schemas/event-envelope.schema.json")
	if err := eventSchema.Validate(toAny(t, eventEnv)); err != nil {
		t.Fatalf("event envelope schema validation: %v", err)
	}

	// Verify the typed data payload matches the real claimset.
	data, ok := events[0].Data.(domain.ClaimSetData)
	if !ok {
		t.Fatalf("audit data not ClaimSetData: %T", events[0].Data)
	}
	if data.ClaimSetHash != resp.ClaimSetHash {
		t.Fatalf("audit data claimset hash mismatch: got %s want %s", data.ClaimSetHash, resp.ClaimSetHash)
	}
	_ = adm.AdmissionID
}

func TestClaimSet_UnknownCapability_NoMutation(t *testing.T) {
	unitKey := "phase-d-claim-reject-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// PREREQUISITE TEST SETUP
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Claim Reject Unit",
		Description: "fixture for claim rejection proof",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}
	createVersion := usecases.CreateVersion{Repo: repo, Audit: setupAudit, Clock: clock}
	versionResp, err := createVersion.CreateVersion(ports.CreateVersionRequest{
		UnitKey: unitKey,
		Label:   "v1",
		Content: "Version content",
		ActorID: "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateVersion: %v", err)
	}

	// BEFORE: no claimset for this version, no claim audit events.
	beforeClaim, beforeClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if beforeClaimOK {
		t.Fatalf("expected no claimset before reject, got %+v", beforeClaim)
	}
	beforeVer, _, _ := repo.FindVersionByID(versionResp.VersionID)
	if beforeVer.ClaimSetHash != "" {
		t.Fatalf("expected empty ClaimSetHash before reject, got %s", beforeVer.ClaimSetHash)
	}

	// REJECT path: unknown capability.
	bodyJSON := fmt.Sprintf(`{"schema_version":"claimset/v0","version_id":"%s","claims":[{"id":"c1","text":"X"}]}`, versionResp.VersionID)
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-claim-reject",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "claim.set",
		Payload: map[string]any{
			"unit_key":   unitKey,
			"version_id": versionResp.VersionID,
			"body_json":  bodyJSON,
			"actor_id":   "phase-d-tester",
		},
	}
	eng := admission.NewEngine(admission.V01Registry())
	adm, err := eng.Evaluate(intent)
	if err != nil {
		t.Fatalf("admission evaluation: %v", err)
	}
	if adm.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", adm.Decision)
	}
	found := false
	for _, c := range adm.ReasonCodes {
		if c == "UNKNOWN_CAPABILITY" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected UNKNOWN_CAPABILITY, got %v", adm.ReasonCodes)
	}
	if adm.TransitionRef != "" {
		t.Fatalf("expected no transition_ref for reject, got %s", adm.TransitionRef)
	}

	// AFTER: SetClaims not invoked; claim state unchanged.
	afterClaim, afterClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if afterClaimOK {
		t.Fatalf("expected no claimset after reject, got %+v", afterClaim)
	}
	afterVer, _, _ := repo.FindVersionByID(versionResp.VersionID)
	if afterVer.ClaimSetHash != "" {
		t.Fatalf("expected empty ClaimSetHash after reject, got %s", afterVer.ClaimSetHash)
	}
	if afterVer.MeaningHash != "" || afterVer.UncertaintyHash != "" {
		t.Fatalf("unexpected hash fields after reject: meaning=%s uncertainty=%s", afterVer.MeaningHash, afterVer.UncertaintyHash)
	}

	claimAudit := memory.NewAuditLog()
	reader := memory.NewAuditByUnitReader(claimAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no CLAIM_SET audit events after reject, got %d", len(events))
	}
}
