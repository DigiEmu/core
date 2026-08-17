// Phase D binding proof for core.version.create:
// A real CreateVersion mutation executes only after P0 ADMIT and emits a
// runtime transition evidence envelope. Test-only; not production enforcement.
package usecases_test

import (
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func newVersionCreateIntent(unitKey, label, content, actor string) admission.Intent {
	return admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-version-" + unitKey,
		CapabilityRef:        "core.version.create",
		AggregateRef:         "unit",
		CommandRef:           "version.create",
		Payload: map[string]any{
			"unit_key":  unitKey,
			"label":     label,
			"content":   content,
			"actor_id":  actor,
		},
	}
}

func TestVersionCreate_Binding_And_EventEvidence(t *testing.T) {
	// Prerequisite: an existing Unit. This setup is not itself routed through
	// Admission; it is a test fixture for the mutation being proven.
	unitKey := "phase-d-version-unit-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Version Unit",
		Description: "created via P0 Phase D version binding test",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	// Admission intent for version.create
	intent := newVersionCreateIntent(
		unitKey,
		"v1",
		"Phase D version content",
		"phase-d-tester",
	)

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
	if adm.TransitionRef != "version:created" {
		t.Fatalf("expected transition version:created, got %s", adm.TransitionRef)
	}
	t.Logf("admission_id=%s", adm.AdmissionID)

	// Real CreateVersion, same as production runtime, gated by admission result.
	versionAudit := memory.NewAuditLog()
	uc := usecases.CreateVersion{Repo: repo, Audit: versionAudit, Clock: clock}
	req := ports.CreateVersionRequest{
		UnitKey:       stringFromMap(intent.Payload, "unit_key"),
		Label:         stringFromMap(intent.Payload, "label"),
		Content:       stringFromMap(intent.Payload, "content"),
		ActorID:       stringFromMap(intent.Payload, "actor_id"),
		BaseVersionID: "",
	}
	resp, err := uc.CreateVersion(req)
	if err != nil {
		t.Fatalf("CreateVersion: %v", err)
	}
	if resp.UnitID != unitResp.UnitID {
		t.Fatalf("version unit id mismatch: got %s want %s", resp.UnitID, unitResp.UnitID)
	}
	if resp.Label != req.Label {
		t.Fatalf("version label mismatch: got %s want %s", resp.Label, req.Label)
	}

	// 1. Version actually persisted.
	v, ok, err := repo.FindVersionByID(resp.VersionID)
	if err != nil {
		t.Fatalf("FindVersionByID: %v", err)
	}
	if !ok {
		t.Fatalf("version not created for id %s", resp.VersionID)
	}
	if v.UnitID != unitResp.UnitID {
		t.Fatalf("created version unit id mismatch: got %s want %s", v.UnitID, unitResp.UnitID)
	}
	if v.Label != req.Label || v.Content != req.Content {
		t.Fatalf("created version mismatch: got %+v", v)
	}

	// 2. Version belongs to the expected Unit.
	vs, err := repo.ListVersionsByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListVersionsByUnitID: %v", err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}

	// 3. Unit head version updated.
	u, ok, _ := repo.FindUnitByID(unitResp.UnitID)
	if !ok {
		t.Fatalf("prerequisite unit not found")
	}
	if u.HeadVersionID != resp.VersionID {
		t.Fatalf("expected head version %s, got %s", resp.VersionID, u.HeadVersionID)
	}

	// 4. Runtime AuditEvent emitted.
	reader := memory.NewAuditByUnitReader(versionAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 version audit event, got %d", len(events))
	}

	// 5. Exact runtime event string, un-normalized.
	if events[0].Type != "version.created" {
		t.Fatalf("expected audit type version.created, got %s", events[0].Type)
	}

	// 6. Event Envelope evidence from the real runtime event.
	evidence := map[string]any{
		"unit_id":         u.ID,
		"version_id":      resp.VersionID,
		"label":           v.Label,
		"content_hash":    v.ContentHash,
		"prev_version_id": v.PrevVersionID,
	}
	eventEnv := map[string]any{
		"schema_version":        "v0.1",
		"architecture_revision": "0.3",
		"event_id":              events[0].ID,
		"command_ref":           "version.create",
		"transition_ref":        "version:created",
		"runtime_event_type":    "version.created",
		"evidence":              evidence,
	}
	eventSchema := loadSchema(t, "schemas/event-envelope.schema.json")
	if err := eventSchema.Validate(toAny(t, eventEnv)); err != nil {
		t.Fatalf("event envelope schema validation: %v", err)
	}

	if events[0].UnitID != u.ID {
		t.Fatalf("audit event unit id mismatch")
	}
	if events[0].ActorID != req.ActorID {
		t.Fatalf("audit event actor mismatch: got %s, want %s", events[0].ActorID, req.ActorID)
	}
	_ = adm.AdmissionID
}

func TestVersionCreate_UnknownCapability_NoMutation(t *testing.T) {
	// Prerequisite Unit so we can prove no version is added on REJECT.
	unitKey := "phase-d-version-reject-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase D Version Reject Unit",
		Description: "fixture for version rejection proof",
		ActorID:     "phase-d-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	// Reject path: unknown capability on an otherwise valid version.create shape.
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-version-reject",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "version.create",
		Payload: map[string]any{
			"unit_key":  unitKey,
			"label":     "v1",
			"content":   "should not be created",
			"actor_id":  "phase-d-tester",
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

	// CreateVersion is NOT invoked. Only the prerequisite Unit exists.
	vs, err := repo.ListVersionsByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListVersionsByUnitID: %v", err)
	}
	if len(vs) != 0 {
		t.Fatalf("expected no versions after reject, got %d", len(vs))
	}

	versionAudit := memory.NewAuditLog()
	reader := memory.NewAuditByUnitReader(versionAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no version audit events, got %d", len(events))
	}
}
