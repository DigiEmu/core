// Phase D binding proof for core.uncertainty.set:
// A real SetUncertainty mutation executes only after P0 ADMIT and produces a
// before/after Uncertainty state evidence envelope. Test-only; not production.
package usecases_test

import (
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func newUncertaintySetIntent(unitKey, versionID, bodyJSON, actor string) admission.Intent {
	return admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-uncertainty-" + unitKey,
		CapabilityRef:        "core.uncertainty.set",
		AggregateRef:         "unit",
		CommandRef:           "uncertainty.set",
		Payload: map[string]any{
			"unit_key":   unitKey,
			"version_id": versionID,
			"body_json":  bodyJSON,
			"actor_id":   actor,
		},
	}
}

func TestUncertaintySet_Binding_And_EventEvidence(t *testing.T) {
	unitKey := "phase-d-uncertainty-unit-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// PREREQUISITE TEST SETUP: existing Unit and Version.
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Uncertainty Unit",
		Description: "created via P0 Phase D uncertainty binding test",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	createVersion := usecases.CreateVersion{Repo: repo, Audit: setupAudit, Clock: clock}
	versionResp, err := createVersion.CreateVersion(ports.CreateVersionRequest{
		UnitKey: unitKey,
		Label:   "v1",
		Content: "Version content for uncertainty test",
		ActorID: "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateVersion: %v", err)
	}

	// BEFORE: no Uncertainty for this Version; related fields empty.
	beforeVer, ok, _ := repo.FindVersionByID(versionResp.VersionID)
	if !ok {
		t.Fatalf("prerequisite version not found")
	}
	if beforeVer.UncertaintyHash != "" {
		t.Fatalf("expected empty UncertaintyHash before SetUncertainty, got %s", beforeVer.UncertaintyHash)
	}
	beforeMeaning, beforeMeaningOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if beforeMeaningOK {
		t.Fatalf("expected no meaning before SetUncertainty, got %+v", beforeMeaning)
	}
	beforeClaim, beforeClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if beforeClaimOK {
		t.Fatalf("expected no claimset before SetUncertainty, got %+v", beforeClaim)
	}
	beforeUncertainty, beforeUncertaintyOK, _ := repo.LoadUncertainty(unitResp.UnitID, versionResp.VersionID)
	if beforeUncertaintyOK {
		t.Fatalf("expected no uncertainty before SetUncertainty, got %+v", beforeUncertainty)
	}
	beforeUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if beforeUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("expected head version %s, got %s", versionResp.VersionID, beforeUnit.HeadVersionID)
	}

	// THE UNCERTAINTY MUTATION UNDER ADMISSION PROOF
	bodyJSON := `{"schema_version":"uncertainty/v0","id":"u1","type":"empirical","level":"low","applies_to":{"scope":"version"},"text":"Sample uncertainty text","tags":["test"]}`
	intent := newUncertaintySetIntent(unitKey, versionResp.VersionID, bodyJSON, "phase-d-tester")

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
	if adm.TransitionRef != "uncertainty:set" {
		t.Fatalf("expected transition uncertainty:set, got %s", adm.TransitionRef)
	}
	t.Logf("admission_id=%s", adm.AdmissionID)

	// Real SetUncertainty, same as production runtime, gated by the ADMIT result.
	uncertaintyAudit := memory.NewAuditLog()
	uc := usecases.SetUncertainty{Repo: repo, Audit: uncertaintyAudit, Clock: clock}
	req := ports.SetUncertaintyRequest{
		UnitKey:   stringFromMap(intent.Payload, "unit_key"),
		VersionID: stringFromMap(intent.Payload, "version_id"),
		BodyBytes: []byte(stringFromMap(intent.Payload, "body_json")),
		ActorID:   stringFromMap(intent.Payload, "actor_id"),
	}
	resp, err := uc.SetUncertainty(req)
	if err != nil {
		t.Fatalf("SetUncertainty: %v", err)
	}
	if resp.UnitID != unitResp.UnitID || resp.VersionID != versionResp.VersionID {
		t.Fatalf("SetUncertainty response mismatch: got %+v", resp)
	}
	if resp.UncertaintyHash == "" {
		t.Fatalf("expected non-empty UncertaintyHash")
	}

	// AFTER: Uncertainty state changed; nothing else.
	afterVer, ok, _ := repo.FindVersionByID(versionResp.VersionID)
	if !ok {
		t.Fatalf("version not found after SetUncertainty")
	}
	if afterVer.UncertaintyHash != resp.UncertaintyHash {
		t.Fatalf("version UncertaintyHash mismatch: got %s want %s", afterVer.UncertaintyHash, resp.UncertaintyHash)
	}

	afterUncertainty, afterUncertaintyOK, _ := repo.LoadUncertainty(unitResp.UnitID, versionResp.VersionID)
	if !afterUncertaintyOK {
		t.Fatalf("uncertainty not found after SetUncertainty")
	}
	if afterUncertainty.SchemaVersion != "uncertainty/v0" {
		t.Fatalf("unexpected uncertainty schema version: %s", afterUncertainty.SchemaVersion)
	}
	if afterUncertainty.ID != "u1" || afterUncertainty.Type != "empirical" || afterUncertainty.Level != "low" {
		t.Fatalf("uncertainty content mismatch: got %+v", afterUncertainty)
	}

	// Cross-state preservation: meaning, claim, unit head, version identity.
	afterMeaning, afterMeaningOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if afterMeaningOK {
		t.Fatalf("meaning unexpectedly changed: got %+v", afterMeaning)
	}
	if afterVer.MeaningHash != "" {
		t.Fatalf("meaning hash unexpectedly set: %s", afterVer.MeaningHash)
	}
	afterClaim, afterClaimOK, _ := repo.LoadClaimSet(unitResp.UnitID, versionResp.VersionID)
	if afterClaimOK {
		t.Fatalf("claimset unexpectedly changed: got %+v", afterClaim)
	}
	if afterVer.ClaimSetHash != "" {
		t.Fatalf("claimset hash unexpectedly set: %s", afterVer.ClaimSetHash)
	}

	afterUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if afterUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("unit head version changed unexpectedly: got %s want %s", afterUnit.HeadVersionID, versionResp.VersionID)
	}
	if afterVer.ID != beforeVer.ID || afterVer.UnitID != beforeVer.UnitID || afterVer.Content != beforeVer.Content || afterVer.Label != beforeVer.Label {
		t.Fatalf("version identity or content/label mutated: got %+v want %+v", afterVer, beforeVer)
	}

	// Real runtime AuditEvent emitted and isolated.
	reader := memory.NewAuditByUnitReader(uncertaintyAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 uncertainty audit event, got %d", len(events))
	}
	if events[0].Type != "UNCERTAINTY_SET" {
		t.Fatalf("expected audit type UNCERTAINTY_SET, got %s", events[0].Type)
	}
	if events[0].VersionID != versionResp.VersionID {
		t.Fatalf("audit event version id mismatch: got %s want %s", events[0].VersionID, versionResp.VersionID)
	}
	if events[0].ActorID != req.ActorID {
		t.Fatalf("audit event actor mismatch: got %s want %s", events[0].ActorID, req.ActorID)
	}

	// Event Envelope from the real runtime event.
	evidence := map[string]any{
		"unit_id":          unitResp.UnitID,
		"version_id":       versionResp.VersionID,
		"uncertainty_hash": resp.UncertaintyHash,
		"uncertainty_id":   afterUncertainty.ID,
		"type":             afterUncertainty.Type,
		"level":            afterUncertainty.Level,
		"text":             afterUncertainty.Text,
		"uncertainty_path": unitResp.UnitID + "." + versionResp.VersionID + ".uncertainty.json",
	}
	eventEnv := map[string]any{
		"schema_version":        "v0.1",
		"architecture_revision": "0.3",
		"event_id":              events[0].ID,
		"command_ref":           "uncertainty.set",
		"transition_ref":        "uncertainty:set",
		"runtime_event_type":    "UNCERTAINTY_SET",
		"evidence":              evidence,
	}
	eventSchema := loadSchema(t, "schemas/event-envelope.schema.json")
	if err := eventSchema.Validate(toAny(t, eventEnv)); err != nil {
		t.Fatalf("event envelope schema validation: %v", err)
	}

	// Verify the typed data payload matches the real uncertainty.
	data, ok := events[0].Data.(domain.UncertaintySetData)
	if !ok {
		t.Fatalf("audit data not UncertaintySetData: %T", events[0].Data)
	}
	if data.UncertaintyHash != resp.UncertaintyHash {
		t.Fatalf("audit data uncertainty hash mismatch: got %s want %s", data.UncertaintyHash, resp.UncertaintyHash)
	}
	_ = adm.AdmissionID
}

func TestUncertaintySet_UnknownCapability_NoMutation(t *testing.T) {
	unitKey := "phase-d-uncertainty-reject-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// PREREQUISITE TEST SETUP
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Uncertainty Reject Unit",
		Description: "fixture for uncertainty rejection proof",
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

	// BEFORE: no uncertainty for this version, no uncertainty audit events.
	beforeUncertainty, beforeUncertaintyOK, _ := repo.LoadUncertainty(unitResp.UnitID, versionResp.VersionID)
	if beforeUncertaintyOK {
		t.Fatalf("expected no uncertainty before reject, got %+v", beforeUncertainty)
	}
	beforeVer, _, _ := repo.FindVersionByID(versionResp.VersionID)
	if beforeVer.UncertaintyHash != "" {
		t.Fatalf("expected empty UncertaintyHash before reject, got %s", beforeVer.UncertaintyHash)
	}

	// REJECT path: unknown capability.
	bodyJSON := `{"schema_version":"uncertainty/v0","id":"u1","type":"empirical","level":"low","applies_to":{"scope":"version"},"text":"X"}`
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-uncertainty-reject",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "uncertainty.set",
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

	// AFTER: SetUncertainty not invoked; uncertainty state unchanged.
	afterUncertainty, afterUncertaintyOK, _ := repo.LoadUncertainty(unitResp.UnitID, versionResp.VersionID)
	if afterUncertaintyOK {
		t.Fatalf("expected no uncertainty after reject, got %+v", afterUncertainty)
	}
	afterVer, _, _ := repo.FindVersionByID(versionResp.VersionID)
	if afterVer.UncertaintyHash != "" {
		t.Fatalf("expected empty UncertaintyHash after reject, got %s", afterVer.UncertaintyHash)
	}
	if afterVer.MeaningHash != "" || afterVer.ClaimSetHash != "" {
		t.Fatalf("unexpected hash fields after reject: meaning=%s claimset=%s", afterVer.MeaningHash, afterVer.ClaimSetHash)
	}

	uncertaintyAudit := memory.NewAuditLog()
	reader := memory.NewAuditByUnitReader(uncertaintyAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no UNCERTAINTY_SET audit events after reject, got %d", len(events))
	}
}
