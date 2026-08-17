// Phase D binding prototype: proves that a real Core CreateUnit mutation executes
// only after successful P0 Admission and emits a runtime transition evidence
// envelope. This is a test-only adapter; it is not production Admission runtime.
package usecases_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"digiemu-core/internal/admission"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

// repoRoot reaches the repository root from internal/kernel/usecases.
func repoRoot() string {
	return filepath.Join("..", "..", "..")
}

func loadSchema(t *testing.T, rel string) *jsonschema.Schema {
	t.Helper()
	p := filepath.Join(repoRoot(), rel)
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("abs schema %s: %v", rel, err)
	}
	compiler := jsonschema.NewCompiler()
	schemaURL := "file:///" + filepath.ToSlash(abs)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile %s: %v", rel, err)
	}
	return sch
}

func stringFromMap(m map[string]any, k string) string {
	v, ok := m[k]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func toAny(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func newUnitCreateIntent(key, title, description, actor string) admission.Intent {
	return admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-" + key,
		CapabilityRef:        "core.unit.create",
		AggregateRef:         "unit",
		CommandRef:           "unit.create",
		Payload: map[string]any{
			"key":         key,
			"title":       title,
			"description": description,
			"actor_id":    actor,
		},
	}
}

func TestUnitCreate_Binding_And_EventEvidence(t *testing.T) {
	intent := newUnitCreateIntent(
		"phase-d-binding-01",
		"Phase D Binding Unit",
		"created via P0 Phase D binding test",
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
	if adm.TransitionRef != "unit:created" {
		t.Fatalf("expected transition unit:created, got %s", adm.TransitionRef)
	}
	t.Logf("admission_id=%s", adm.AdmissionID)

	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}
	uc := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}

	req := ports.CreateUnitRequest{
		Key:         stringFromMap(intent.Payload, "key"),
		Title:       stringFromMap(intent.Payload, "title"),
		Description: stringFromMap(intent.Payload, "description"),
		ActorID:     stringFromMap(intent.Payload, "actor_id"),
	}
	resp, err := uc.CreateUnit(req)
	if err != nil {
		t.Fatalf("CreateUnit: %v", err)
	}
	if resp.Key != req.Key || resp.Title != req.Title {
		t.Fatalf("response mismatch: got %v", resp)
	}

	u, ok, err := repo.FindUnitByKey(resp.Key)
	if err != nil {
		t.Fatalf("FindUnitByKey: %v", err)
	}
	if !ok {
		t.Fatalf("unit not created for key %s", resp.Key)
	}
	if u.Key != req.Key || u.Title != req.Title {
		t.Fatalf("created unit mismatch: got %+v", u)
	}

	reader := memory.NewAuditByUnitReader(audit)
	events, err := reader.ListByUnitID(u.ID)
	if err != nil {
		t.Fatalf("ListByUnitID: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	if events[0].Type != "unit.created" {
		t.Fatalf("expected audit type unit.created, got %s", events[0].Type)
	}

	evidence := map[string]any{
		"unit_id": u.ID,
		"key":     resp.Key,
		"title":   resp.Title,
	}
	eventEnv := map[string]any{
		"schema_version":        "v0.1",
		"architecture_revision": "0.3",
		"event_id":              events[0].ID,
		"command_ref":           "unit.create",
		"transition_ref":        "unit:created",
		"runtime_event_type":    "unit.created",
		"evidence":              evidence,
	}
	eventSchema := loadSchema(t, "schemas/event-envelope.schema.json")
	if err := eventSchema.Validate(eventEnv); err != nil {
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

func TestUnitCreate_UnknownCapability_NoMutation(t *testing.T) {
	intent := admission.Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-reject-01",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "unit.create",
		Payload:              map[string]any{},
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

	repo := memory.NewUnitRepo()
	audit := memory.NewAuditLog()
	clock := memory.FakeClock{Now: 1700000000}
	_ = usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}

	units, err := repo.ListUnits()
	if err != nil {
		t.Fatalf("ListUnits: %v", err)
	}
	if len(units) != 0 {
		t.Fatalf("expected no units, got %d", len(units))
	}
	if len(audit.Events) != 0 {
		t.Fatalf("expected no audit events, got %d", len(audit.Events))
	}
	if adm.TransitionRef != "" {
		t.Fatalf("expected no transition_ref for reject, got %s", adm.TransitionRef)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
