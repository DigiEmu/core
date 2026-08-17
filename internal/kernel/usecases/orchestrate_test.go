// Phase E test-only orchestration increments.
// Proves the semantic chain from Admission decision to Core execution,
// without becoming production orchestration.
package usecases_test

import (
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

// admissionGateProbe is a test-only helper that evaluates a P0 Admission Intent
// and records whether the post-ADMIT execution closure was invoked.
// It performs no state inspection, no audit inspection, no dispatch routing,
// and no outcome classification beyond the admission decision itself.
type admissionGateProbe struct {
	handlerCalled bool
}

// evaluate runs the Admission Engine and, only when the decision is ADMIT,
// records that the closure was called and invokes it.
// It does not route by transition, construct requests, or inspect Core state.
// The closure returns any handler error so the caller can distinguish
// ADMIT from successful execution.
func (p *admissionGateProbe) evaluate(eng *admission.Engine, intent admission.Intent, onAdmit func() error) (admission.Result, error) {
	adm, err := eng.Evaluate(intent)
	if err != nil {
		return adm, err
	}
	if adm.Decision != "ADMIT" {
		return adm, nil
	}
	p.handlerCalled = true
	if onAdmit != nil {
		if err := onAdmit(); err != nil {
			return adm, err
		}
	}
	return adm, nil
}

// TestPhaseE_Reject_DoesNotInvokeHandler proves that a normative REJECT from
// the reusable Admission Engine prevents the execution closure from running.
func TestPhaseE_Reject_DoesNotInvokeHandler(t *testing.T) {
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-e-reject-01",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "unit.create",
		Payload:              map[string]any{},
	}

	eng := admission.NewEngine(admission.V01Registry())
	probe := &admissionGateProbe{}

	onAdmit := func() error {
		// If the gate helper incorrectly invokes the closure on REJECT, fail
		// immediately so the test cannot pass.
		t.Fatalf("execution closure was invoked after %s", "REJECT")
		return nil
	}

	adm, err := probe.evaluate(eng, intent, onAdmit)
	if err != nil {
		t.Fatalf("admission evaluation: %v", err)
	}
	if adm.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", adm.Decision)
	}

	// Verify the expected normative reason code from the rule registry.
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

	// The primary proof: the execution closure must not have been invoked.
	if probe.handlerCalled {
		t.Fatalf("admissionGateProbe recorded handlerCalled=true after REJECT")
	}
}

// TestPhaseE_Admit_CreateVersion_ExecutesAndProducesCoherentStateAndAudit proves
// the positive path: ADMIT -> real CreateVersion -> coherent state + audit.
func TestPhaseE_Admit_CreateVersion_ExecutesAndProducesCoherentStateAndAudit(t *testing.T) {
	unitKey := "phase-e-version-unit-01"
	repo := memory.NewUnitRepo()
	setupAudit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}

	// Prerequisite: a Unit created with the existing real use case.
	createUnit := usecases.CreateUnit{Repo: repo, Audit: setupAudit, Clock: clock}
	unitResp, err := createUnit.CreateUnit(ports.CreateUnitRequest{
		Key:         unitKey,
		Title:       "Phase E Version Unit",
		Description: "fixture for Phase E CreateVersion admission test",
		ActorID:     "phase-e-tester",
	})
	if err != nil {
		t.Fatalf("prerequisite CreateUnit: %v", err)
	}

	label := "v1"
	content := "Phase E version content"
	actor := "phase-e-tester"

	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-e-version-01",
		CapabilityRef:        "core.version.create",
		AggregateRef:         "unit",
		CommandRef:           "version.create",
		Payload: map[string]any{
			"unit_key": unitKey,
			"label":    label,
			"content":  content,
			"actor_id": actor,
		},
	}

	eng := admission.NewEngine(admission.V01Registry())
	probe := &admissionGateProbe{}

	// Audit log isolated to the CreateVersion operation under test.
	versionAudit := memory.NewAuditLog()
	var versionResp ports.CreateVersionResponse

	// Verify the prerequisite Unit has no head version before CreateVersion.
	beforeUnit, ok, _ := repo.FindUnitByID(unitResp.UnitID)
	if !ok {
		t.Fatalf("prerequisite unit not found")
	}
	if beforeUnit.HeadVersionID != "" {
		t.Fatalf("prerequisite unit already has a head version: %s", beforeUnit.HeadVersionID)
	}

	onAdmit := func() error {
		createVersion := usecases.CreateVersion{Repo: repo, Audit: versionAudit, Clock: clock}
		req := ports.CreateVersionRequest{
			UnitKey:       unitKey,
			Label:         label,
			Content:       content,
			ActorID:       actor,
			BaseVersionID: "",
		}
		resp, err := createVersion.CreateVersion(req)
		versionResp = resp
		return err
	}

	adm, err := probe.evaluate(eng, intent, onAdmit)
	if err != nil {
		t.Fatalf("admission evaluation or CreateVersion: %v", err)
	}
	if adm.Decision != "ADMIT" {
		t.Fatalf("expected ADMIT, got %s", adm.Decision)
	}
	if adm.TransitionRef != "version:created" {
		t.Fatalf("expected transition version:created, got %s", adm.TransitionRef)
	}
	if !probe.handlerCalled {
		t.Fatalf("admissionGateProbe recorded handlerCalled=false after ADMIT")
	}

	// STATE PROOF: the created Version exists and is coherent.
	v, ok, err := repo.FindVersionByID(versionResp.VersionID)
	if err != nil {
		t.Fatalf("FindVersionByID: %v", err)
	}
	if !ok {
		t.Fatalf("version not created for id %s", versionResp.VersionID)
	}
	if v.ID != versionResp.VersionID {
		t.Fatalf("version id mismatch: got %s want %s", v.ID, versionResp.VersionID)
	}
	if v.UnitID != unitResp.UnitID {
		t.Fatalf("version unit id mismatch: got %s want %s", v.UnitID, unitResp.UnitID)
	}
	if v.PrevVersionID != "" {
		t.Fatalf("first version PrevVersionID should be empty, got %s", v.PrevVersionID)
	}
	if v.Label != label || v.Content != content {
		t.Fatalf("version content mismatch: got %+v", v)
	}

	// The Unit lists the new Version.
	vs, err := repo.ListVersionsByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListVersionsByUnitID: %v", err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
	if vs[0].ID != versionResp.VersionID {
		t.Fatalf("listed version id mismatch: got %s want %s", vs[0].ID, versionResp.VersionID)
	}

	// Head points to the new Version.
	u, ok, _ := repo.FindUnitByID(unitResp.UnitID)
	if !ok {
		t.Fatalf("prerequisite unit not found")
	}
	if u.HeadVersionID != versionResp.VersionID {
		t.Fatalf("expected head version %s, got %s", versionResp.VersionID, u.HeadVersionID)
	}

	// AUDIT PROOF: one version.created event for the operation.
	reader := memory.NewAuditByUnitReader(versionAudit)
	events, err := reader.ListByUnitID(unitResp.UnitID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	if events[0].Type != "version.created" {
		t.Fatalf("expected audit type version.created, got %s", events[0].Type)
	}
	if events[0].UnitID != unitResp.UnitID {
		t.Fatalf("audit event unit id mismatch: got %s want %s", events[0].UnitID, unitResp.UnitID)
	}
	if events[0].VersionID != versionResp.VersionID {
		t.Fatalf("audit event version id mismatch: got %s want %s", events[0].VersionID, versionResp.VersionID)
	}
	if events[0].ActorID != actor {
		t.Fatalf("audit event actor mismatch: got %s want %s", events[0].ActorID, actor)
	}

	// The successful path does not prove crash durability, filesystem atomicity,
	// or production concurrency. It proves only adapter-level ADMIT, execution,
	// coherent state, and observable audit in the memory test fixture.
}

// TestPhaseE_Admit_CreateVersion_HandlerError_DoesNotRewriteAdmission proves that
// a real handler execution error is separate from the Admission decision: ADMIT
// remains ADMIT even when the real CreateVersion use case returns a domain error.
func TestPhaseE_Admit_CreateVersion_HandlerError_DoesNotRewriteAdmission(t *testing.T) {
	repo := memory.NewUnitRepo()
	clock := memory.FakeClock{Now: 1700000000}

	missingUnitKey := "phase-e-missing-unit"
	label := "v1"
	content := "some content"
	actor := "phase-e-tester"

	// A valid CreateVersion Intent. Admission is independent of the Unit's
	// existence; that check belongs to the Core handler.
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-e-version-notfound-01",
		CapabilityRef:        "core.version.create",
		AggregateRef:         "unit",
		CommandRef:           "version.create",
		Payload: map[string]any{
			"unit_key": missingUnitKey,
			"label":    label,
			"content":  content,
			"actor_id": actor,
		},
	}

	eng := admission.NewEngine(admission.V01Registry())
	probe := &admissionGateProbe{}

	versionAudit := memory.NewAuditLog()

	onAdmit := func() error {
		createVersion := usecases.CreateVersion{Repo: repo, Audit: versionAudit, Clock: clock}
		req := ports.CreateVersionRequest{
			UnitKey:       missingUnitKey,
			Label:         label,
			Content:       content,
			ActorID:       actor,
			BaseVersionID: "",
		}
		_, err := createVersion.CreateVersion(req)
		return err
	}

	adm, err := probe.evaluate(eng, intent, onAdmit)
	if err != domain.ErrUnitNotFound {
		t.Fatalf("expected ErrUnitNotFound from handler, got %v", err)
	}

	// ADMISSION STABILITY: the decision remains ADMIT regardless of the
	// subsequent handler error.
	if adm.Decision != "ADMIT" {
		t.Fatalf("expected ADMIT after handler error, got %s", adm.Decision)
	}
	if adm.TransitionRef != "version:created" {
		t.Fatalf("expected transition version:created, got %s", adm.TransitionRef)
	}
	if len(adm.ReasonCodes) != 0 {
		t.Fatalf("expected no REJECT reason codes after ADMIT, got %v", adm.ReasonCodes)
	}

	if !probe.handlerCalled {
		t.Fatalf("admissionGateProbe recorded handlerCalled=false after ADMIT")
	}

	// This specific domain error occurs before any Core state change or audit
	// append. This assertion does not generalize to all execution errors,
	// especially not to partial-failure cases.
	units, err := repo.ListUnits()
	if err != nil {
		t.Fatalf("ListUnits: %v", err)
	}
	if len(units) != 0 {
		t.Fatalf("expected 0 units, got %d", len(units))
	}
	if len(versionAudit.Events) != 0 {
		t.Fatalf("expected 0 audit events, got %d", len(versionAudit.Events))
	}
}
