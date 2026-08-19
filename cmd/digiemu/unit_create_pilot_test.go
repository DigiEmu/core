package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"digiemu-core/internal/admission"
	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

func TestBuildCreateUnitIntent(t *testing.T) {
	intent := buildCreateUnitIntent("my-key", "My Title", "My Desc")

	if intent.SchemaVersion != "v0.1" {
		t.Fatalf("schema_version=%q, want v0.1", intent.SchemaVersion)
	}
	if intent.ArchitectureRevision != "0.3" {
		t.Fatalf("architecture_revision=%q, want 0.3", intent.ArchitectureRevision)
	}
	if intent.IntentID != "p0-pilot-unresolved-intent-id" {
		t.Fatalf("intent_id=%q, want non-normative placeholder", intent.IntentID)
	}
	if intent.CapabilityRef != "core.unit.create" {
		t.Fatalf("capability_ref=%q", intent.CapabilityRef)
	}
	if intent.AggregateRef != "unit" {
		t.Fatalf("aggregate_ref=%q", intent.AggregateRef)
	}
	if intent.CommandRef != "unit.create" {
		t.Fatalf("command_ref=%q", intent.CommandRef)
	}
	if intent.Payload["key"] != "my-key" {
		t.Fatalf("payload.key=%v", intent.Payload["key"])
	}
	if intent.Payload["title"] != "My Title" {
		t.Fatalf("payload.title=%v", intent.Payload["title"])
	}
	if intent.Payload["description"] != "My Desc" {
		t.Fatalf("payload.description=%v", intent.Payload["description"])
	}
}

func TestAcquirePilotLock(t *testing.T) {
	dir := t.TempDir()

	release, err := acquirePilotLock(dir)
	if err != nil {
		t.Fatalf("acquire lock: %v", err)
	}
	lockPath := filepath.Join(dir, ".p0-admission-create-unit.lock")
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("lock file not created: %v", err)
	}

	// Second acquisition must fail closed.
	if _, err := acquirePilotLock(dir); err == nil {
		t.Fatalf("expected lock conflict, got nil")
	}

	// After release, acquisition succeeds again.
	release()
	release2, err := acquirePilotLock(dir)
	if err != nil {
		t.Fatalf("re-acquire after release: %v", err)
	}
	release2()
}

func TestIsPrePersistenceCreateUnitError(t *testing.T) {
	pre := []error{
		domain.ErrInvalidUnitKey,
		domain.ErrInvalidUnitTitle,
		domain.ErrUnitAlreadyExists,
		domain.ErrAuditNotConfigured,
		domain.ErrClockNotConfigured,
	}
	for _, e := range pre {
		if !isPrePersistenceCreateUnitError(e) {
			t.Fatalf("%v should be classified as pre-persistence", e)
		}
	}
	if isPrePersistenceCreateUnitError(fmt.Errorf("some repo error")) {
		t.Fatalf("arbitrary repo error should not be pre-persistence")
	}
}

func TestInspectCreateUnitState(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	// No Unit present.
	msg := inspectCreateUnitState(repo, "missing", "Title", "")
	if msg != "no coherent Unit currently observable" {
		t.Fatalf("absent: %q", msg)
	}

	// Create a Unit through the existing use case.
	audit := fsrepo.NewAuditLog(dir)
	clock := memory.FakeClock{Now: 1700000000}
	uc := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	in := ports.CreateUnitRequest{
		Key:         "my-key",
		Title:       "My Title",
		Description: "My Desc",
		ActorID:     "pilot",
	}
	if _, err := uc.CreateUnit(in); err != nil {
		t.Fatalf("create unit: %v", err)
	}

	msg = inspectCreateUnitState(repo, "my-key", "My Title", "My Desc")
	if msg != "coherent Unit currently observable" {
		t.Fatalf("present: %q", msg)
	}

	// A different Unit is observable for this key: cannot attribute causation.
	msg = inspectCreateUnitState(repo, "my-key", "Other Title", "My Desc")
	if msg != "Unit observable; causation unresolved" {
		t.Fatalf("unexpected: %q", msg)
	}
}

func TestRunCreateUnitPilot_LockConflict(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, ".p0-admission-create-unit.lock")
	if err := os.WriteFile(lockPath, nil, 0o644); err != nil {
		t.Fatalf("create lock file: %v", err)
	}

	_, err := runCreateUnitPilot(dir, "key", "Title", "", nil, nil, nil)
	if err == nil {
		t.Fatalf("expected lock conflict error")
	}
	if !strings.Contains(err.Error(), "pilot lock already held") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCreateUnitPilot_EngineError(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	evalCalls := 0
	createCalls := 0

	evaluate := func(intent admission.Intent) (admission.Result, error) {
		evalCalls++
		return admission.Result{}, fmt.Errorf("registry unavailable")
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		createCalls++
		return ports.CreateUnitResponse{}, nil
	}

	res, err := runCreateUnitPilot(dir, "key", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if evalCalls != 1 {
		t.Fatalf("evaluate calls=%d, want 1", evalCalls)
	}
	if createCalls != 0 {
		t.Fatalf("CreateUnit calls=%d, want 0", createCalls)
	}
	if res.decision != "" {
		t.Fatalf("decision should be empty for engine error, got %q", res.decision)
	}
	if !strings.Contains(res.message, "admission evaluation failed") {
		t.Fatalf("message=%q", res.message)
	}
}

func TestRunCreateUnitPilot_Reject(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	evalCalls := 0
	createCalls := 0

	evaluate := func(intent admission.Intent) (admission.Result, error) {
		evalCalls++
		return admission.Result{
			AdmissionID: "admission:test",
			Decision:    "REJECT",
			ReasonCodes: []string{"MISSING_REQUIRED_FIELD"},
		}, nil
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		createCalls++
		return ports.CreateUnitResponse{}, nil
	}

	res, err := runCreateUnitPilot(dir, "key", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if evalCalls != 1 {
		t.Fatalf("evaluate calls=%d, want 1", evalCalls)
	}
	if createCalls != 0 {
		t.Fatalf("CreateUnit calls=%d, want 0", createCalls)
	}
	if res.decision != "REJECT" {
		t.Fatalf("decision=%s, want REJECT", res.decision)
	}
	if res.admissionID != "admission:test" {
		t.Fatalf("admission_id=%s", res.admissionID)
	}
	if len(res.reasonCodes) != 1 || res.reasonCodes[0] != "MISSING_REQUIRED_FIELD" {
		t.Fatalf("reason_codes=%v", res.reasonCodes)
	}
}

func TestRunCreateUnitPilot_AdmitSuccess(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	evalCalls := 0
	createCalls := 0

	evaluate := func(intent admission.Intent) (admission.Result, error) {
		evalCalls++
		return admission.Result{
			AdmissionID:   "admission:admit",
			Decision:      "ADMIT",
			TransitionRef: "unit:created",
		}, nil
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		createCalls++
		return ports.CreateUnitResponse{UnitID: "unit-123", Key: "my-key"}, nil
	}

	res, err := runCreateUnitPilot(dir, "my-key", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if evalCalls != 1 {
		t.Fatalf("evaluate calls=%d, want 1", evalCalls)
	}
	if createCalls != 1 {
		t.Fatalf("CreateUnit calls=%d, want 1", createCalls)
	}
	if res.decision != "ADMIT" {
		t.Fatalf("decision=%s, want ADMIT", res.decision)
	}
	if res.transitionRef != "unit:created" {
		t.Fatalf("transition_ref=%s", res.transitionRef)
	}
	if res.unitID != "unit-123" {
		t.Fatalf("unitID=%s", res.unitID)
	}
}

func TestRunCreateUnitPilot_AdmitPrePersistenceDomainError(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	evaluate := func(intent admission.Intent) (admission.Result, error) {
		return admission.Result{
			AdmissionID:   "admission:admit",
			Decision:      "ADMIT",
			TransitionRef: "unit:created",
		}, nil
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		return ports.CreateUnitResponse{}, domain.ErrInvalidUnitKey
	}

	res, err := runCreateUnitPilot(dir, "xx", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if res.decision != "ADMIT" {
		t.Fatalf("decision=%s, want ADMIT", res.decision)
	}
	if !strings.Contains(res.message, "execution/domain error") {
		t.Fatalf("message=%q", res.message)
	}
	if res.stateMsg != "persistence not attempted by this call" {
		t.Fatalf("stateMsg=%q, want 'persistence not attempted by this call'", res.stateMsg)
	}
}

func TestRunCreateUnitPilot_UncertainExecutionError_InspectionAbsent(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	evaluate := func(intent admission.Intent) (admission.Result, error) {
		return admission.Result{
			AdmissionID:   "admission:admit",
			Decision:      "ADMIT",
			TransitionRef: "unit:created",
		}, nil
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		return ports.CreateUnitResponse{}, fmt.Errorf("audit failed")
	}

	res, err := runCreateUnitPilot(dir, "my-key", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if res.decision != "ADMIT" {
		t.Fatalf("decision=%s, want ADMIT", res.decision)
	}
	if !strings.Contains(res.message, "execution error") {
		t.Fatalf("message=%q", res.message)
	}
	if res.stateMsg != "no coherent Unit currently observable" {
		t.Fatalf("stateMsg=%q, want 'no coherent Unit currently observable'", res.stateMsg)
	}
}

func TestRunCreateUnitPilot_UncertainExecutionError_InspectionError(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	// Place a corrupt .json file in the units directory to force FindUnitByKey
	// to return an observation error.
	if err := os.WriteFile(filepath.Join(dir, "units", "corrupt.json"), []byte("not-json"), 0o644); err != nil {
		t.Fatalf("write corrupt unit: %v", err)
	}

	evaluate := func(intent admission.Intent) (admission.Result, error) {
		return admission.Result{
			AdmissionID:   "admission:admit",
			Decision:      "ADMIT",
			TransitionRef: "unit:created",
		}, nil
	}
	create := func(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
		return ports.CreateUnitResponse{}, fmt.Errorf("save failed")
	}

	res, err := runCreateUnitPilot(dir, "my-key", "Title", "", evaluate, create, repo)
	if err != nil {
		t.Fatalf("unexpected pilot error: %v", err)
	}
	if res.decision != "ADMIT" {
		t.Fatalf("decision=%s, want ADMIT", res.decision)
	}
	if res.stateMsg != "repository observation failed; mutation state unresolved" {
		t.Fatalf("stateMsg=%q", res.stateMsg)
	}
}
