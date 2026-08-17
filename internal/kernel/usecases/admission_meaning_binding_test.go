// Phase D binding proof for core.meaning.set:
// A real SetMeaning mutation executes only after P0 ADMIT and produces a
// before/after semantic-state evidence envelope. Test-only; not production.
package usecases_test

import (
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func newMeaningSetIntent(unitKey, versionID, meaningJSON, actor string) admission.Intent {
	return admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-meaning-" + unitKey,
		CapabilityRef:        "core.meaning.set",
		AggregateRef:         "unit",
		CommandRef:           "meaning.set",
		Payload: map[string]any{
			"unit_key":   unitKey,
			"version_id": versionID,
			"meaning_json": meaningJSON,
			"actor_id":   actor,
		},
	}
}

func TestMeaningSet_Binding_And_EventEvidence(t *testing.T) {
	unitKey := "phase-d-meaning-unit-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// TEST FIXTURE SETUP: existing Unit and Version.
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Meaning Unit",
		Description: "created via P0 Phase D meaning binding test",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	createVersion := usecases.CreateVersion{Repo: repo, Audit: setupAudit, Clock: clock}
	versionResp, err := createVersion.CreateVersion(ports.CreateVersionRequest{
		UnitKey: unitKey,
		Label:   "v1",
		Content: "Version content for meaning test",
		ActorID: "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateVersion: %v", err)
	}

	// BEFORE: capture semantic state before the Admission-gated mutation.
	beforeVer, ok, err := repo.FindVersionByID(versionResp.VersionID)
	if err != nil || !ok {
		t.Fatalf("prerequisite version not found: err=%v ok=%t", err, ok)
	}
	if beforeVer.MeaningHash != "" {
		t.Fatalf("expected empty MeaningHash before SetMeaning, got %s", beforeVer.MeaningHash)
	}
	beforeMeaning, beforeOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if beforeOK {
		t.Fatalf("expected no meaning stored before SetMeaning, got %+v", beforeMeaning)
	}
	beforeUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if beforeUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("expected head version %s, got %s", versionResp.VersionID, beforeUnit.HeadVersionID)
	}

	// THE MUTATION UNDER ADMISSION PROOF
	meaningJSON := `{"schema_version":"meaning/v1","title":"Meaning Title","purpose":"Meaning Purpose"}`
	intent := newMeaningSetIntent(unitKey, versionResp.VersionID, meaningJSON, "phase-d-tester")

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
	if adm.TransitionRef != "meaning:set" {
		t.Fatalf("expected transition meaning:set, got %s", adm.TransitionRef)
	}
	t.Logf("admission_id=%s", adm.AdmissionID)

	// Real SetMeaning, same as production runtime, gated by the ADMIT result.
	meaningAudit := memory.NewAuditLog()
	uc := usecases.SetMeaning{Repo: repo, Audit: meaningAudit, Clock: clock}
	req := ports.SetMeaningRequest{
		UnitKey:     stringFromMap(intent.Payload, "unit_key"),
		VersionID:   stringFromMap(intent.Payload, "version_id"),
		MeaningJSON: []byte(stringFromMap(intent.Payload, "meaning_json")),
		ActorID:     stringFromMap(intent.Payload, "actor_id"),
	}
	resp, err := uc.SetMeaning(req)
	if err != nil {
		t.Fatalf("SetMeaning: %v", err)
	}
	if resp.UnitID != unitResp.UnitID || resp.VersionID != versionResp.VersionID {
		t.Fatalf("SetMeaning response mismatch: got %+v", resp)
	}
	if resp.MeaningHash == "" {
		t.Fatalf("expected non-empty MeaningHash")
	}

	// AFTER: prove semantic state changed and nothing else.
	afterVer, ok, _ := repo.FindVersionByID(versionResp.VersionID)
	if !ok {
		t.Fatalf("version not found after SetMeaning")
	}
	if afterVer.MeaningHash != resp.MeaningHash {
		t.Fatalf("version MeaningHash mismatch: got %s want %s", afterVer.MeaningHash, resp.MeaningHash)
	}

	afterMeaning, afterOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if !afterOK {
		t.Fatalf("meaning not found after SetMeaning")
	}
	if afterMeaning.Title != "Meaning Title" || afterMeaning.Purpose != "Meaning Purpose" {
		t.Fatalf("meaning content mismatch: got title=%q purpose=%q", afterMeaning.Title, afterMeaning.Purpose)
	}

	// Version identity preserved, Unit head unchanged.
	afterUnit, _, _ := repo.FindUnitByID(unitResp.UnitID)
	if afterUnit.HeadVersionID != versionResp.VersionID {
		t.Fatalf("unit head version changed unexpectedly: got %s want %s", afterUnit.HeadVersionID, versionResp.VersionID)
	}
	if afterVer.ID != beforeVer.ID || afterVer.UnitID != beforeVer.UnitID || afterVer.Content != beforeVer.Content {
		t.Fatalf("version identity or content mutated: got %+v want %+v", afterVer, beforeVer)
	}

	// Real runtime AuditEvent emitted and isolated.
	reader := memory.NewAuditByUnitReader(meaningAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 meaning audit event, got %d", len(events))
	}
	if events[0].Type != "MEANING_SET" {
		t.Fatalf("expected audit type MEANING_SET, got %s", events[0].Type)
	}
	if events[0].VersionID != versionResp.VersionID {
		t.Fatalf("audit event version id mismatch: got %s want %s", events[0].VersionID, versionResp.VersionID)
	}
	if events[0].ActorID != req.ActorID {
		t.Fatalf("audit event actor mismatch: got %s want %s", events[0].ActorID, req.ActorID)
	}

	// Event Envelope from the real runtime event.
	evidence := map[string]any{
		"unit_id":      unitResp.UnitID,
		"version_id":   versionResp.VersionID,
		"meaning_hash": resp.MeaningHash,
		"title":        afterMeaning.Title,
		"purpose":      afterMeaning.Purpose,
	}
	eventEnv := map[string]any{
		"schema_version":        "v0.1",
		"architecture_revision": "0.3",
		"event_id":              events[0].ID,
		"command_ref":           "meaning.set",
		"transition_ref":        "meaning:set",
		"runtime_event_type":    "MEANING_SET",
		"evidence":              evidence,
	}
	eventSchema := loadSchema(t, "schemas/event-envelope.schema.json")
	if err := eventSchema.Validate(toAny(t, eventEnv)); err != nil {
		t.Fatalf("event envelope schema validation: %v", err)
	}

	// Optional: verify the typed data payload matches the real meaning.
	data, ok := events[0].Data.(domain.MeaningSetData)
	if !ok {
		t.Fatalf("audit data not MeaningSetData: %T", events[0].Data)
	}
	if data.MeaningHash != resp.MeaningHash {
		t.Fatalf("audit data meaning hash mismatch: got %s want %s", data.MeaningHash, resp.MeaningHash)
	}
	_ = adm.AdmissionID
}

func TestMeaningSet_UnknownCapability_NoMutation(t *testing.T) {
	unitKey := "phase-d-meaning-reject-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// TEST FIXTURE SETUP
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Meaning Reject Unit",
		Description: "fixture for meaning rejection proof",
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

	// BEFORE: no meaning for this version, no meaning audit events.
	beforeMeaning, beforeOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if beforeOK {
		t.Fatalf("expected no meaning before reject, got %+v", beforeMeaning)
	}

	// REJECT path: unknown capability.
	meaningJSON := `{"schema_version":"meaning/v1","title":"X","purpose":"Y"}`
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-meaning-reject",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "meaning.set",
		Payload: map[string]any{
			"unit_key":     unitKey,
			"version_id":   versionResp.VersionID,
			"meaning_json": meaningJSON,
			"actor_id":     "phase-d-tester",
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

	// AFTER: SetMeaning not invoked; semantic state unchanged.
	afterMeaning, afterOK, _ := repo.LoadMeaning(unitResp.UnitID, versionResp.VersionID)
	if afterOK {
		t.Fatalf("expected no meaning after reject, got %+v", afterMeaning)
	}
	ver, _, _ := repo.FindVersionByID(versionResp.VersionID)
	if ver.MeaningHash != "" {
		t.Fatalf("expected empty MeaningHash after reject, got %s", ver.MeaningHash)
	}

	meaningAudit := memory.NewAuditLog()
	reader := memory.NewAuditByUnitReader(meaningAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no MEANING_SET audit events after reject, got %d", len(events))
	}
}
