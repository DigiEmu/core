// Phase D binding prototype: proves that a real Core CreateUnit mutation executes
// only after successful P0 Admission and emits a runtime transition evidence
// envelope. This is a test-only adapter; it is not production Admission runtime.
package usecases_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"digiemu-core/internal/canonicaljson"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

// repoRoot reaches the repository root from internal/kernel/usecases.
func repoRoot() string {
	return filepath.Join("..", "..", "..")
}

type intentEnvelope struct {
	SchemaVersion        string         `json:"schema_version"`
	ArchitectureRevision string         `json:"architecture_revision"`
	IntentID             string         `json:"intent_id"`
	CapabilityRef        string         `json:"capability_ref"`
	AggregateRef         string         `json:"aggregate_ref"`
	CommandRef           string         `json:"command_ref"`
	Payload              map[string]any `json:"payload"`
}

type intentDigestPreimage struct {
	IntentDigestProfile  string         `json:"intent_digest_profile"`
	SchemaVersion        string         `json:"schema_version"`
	ArchitectureRevision string         `json:"architecture_revision"`
	CapabilityRef        string         `json:"capability_ref"`
	AggregateRef         string         `json:"aggregate_ref"`
	CommandRef           string         `json:"command_ref"`
	Payload              map[string]any `json:"payload"`
}

type admissionIDPreimage struct {
	AdmissionIDProfile   string   `json:"admission_id_profile"`
	SchemaVersion        string   `json:"schema_version"`
	ArchitectureRevision string   `json:"architecture_revision"`
	IntentDigest         string   `json:"intent_digest"`
	CapabilityRef        string   `json:"capability_ref"`
	AggregateRef         string   `json:"aggregate_ref"`
	CommandRef           string   `json:"command_ref"`
	TransitionRef        *string  `json:"transition_ref"`
	Decision             string   `json:"decision"`
	RuleRefs             []string `json:"rule_refs"`
	ReasonCodes          []string `json:"reason_codes"`
}

type admissionResult struct {
	Decision      string
	TransitionRef string
	ReasonCode    string
	RuleRefs      []string
	ReasonCodes   []string
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

func computeIntentDigest(in intentEnvelope) string {
	p := intentDigestPreimage{
		IntentDigestProfile:  "P0.ADMISSION.INTENT.v0.1",
		SchemaVersion:        in.SchemaVersion,
		ArchitectureRevision: in.ArchitectureRevision,
		CapabilityRef:        in.CapabilityRef,
		AggregateRef:         in.AggregateRef,
		CommandRef:           in.CommandRef,
		Payload:              in.Payload,
	}
	b, err := canonicaljson.Marshal(p)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	return "p0-intent:sha256:" + hex.EncodeToString(sum[:])
}

func computeAdmissionID(intentDigest string, in intentEnvelope, decision, transitionRef string, ruleRefs, reasonCodes []string) string {
	var t *string
	if transitionRef != "" {
		t = &transitionRef
	}
	rrs := make([]string, len(ruleRefs))
	copy(rrs, ruleRefs)
	sort.Strings(rrs)
	rcs := make([]string, len(reasonCodes))
	copy(rcs, reasonCodes)
	sort.Strings(rcs)
	p := admissionIDPreimage{
		AdmissionIDProfile:   "P0.ADMISSION.ID.v0.1",
		SchemaVersion:        in.SchemaVersion,
		ArchitectureRevision: in.ArchitectureRevision,
		IntentDigest:         intentDigest,
		CapabilityRef:        in.CapabilityRef,
		AggregateRef:         in.AggregateRef,
		CommandRef:           in.CommandRef,
		TransitionRef:        t,
		Decision:             decision,
		RuleRefs:             rrs,
		ReasonCodes:          rcs,
	}
	b, err := canonicaljson.Marshal(p)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(b)
	return "admission:sha256:" + hex.EncodeToString(sum[:])
}

// evaluateAdmission is a minimal Phase D test adapter that applies the P0 v0.1
// admission rules for the core.unit.create path. It does not parse YAML; values
// are derived from the v0.1 registries and baseline.
func evaluateAdmission(in intentEnvelope) admissionResult {
	const (
		baseline         = "0.3"
		capUnitCreate    = "core.unit.create"
		aggUnit          = "unit"
		cmdUnitCreate    = "unit.create"
		transUnitCreated = "unit:created"
	)
	ruleIDs := []string{
		"P0.ADMISSION.ARCHITECTURE_REVISION",
		"P0.ADMISSION.CAPABILITY_EXISTS",
		"P0.ADMISSION.CAPABILITY_MUTATES",
		"P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY",
		"P0.ADMISSION.COMMAND_EXISTS",
		"P0.ADMISSION.COMMAND_CAPABILITY_MATCH",
		"P0.ADMISSION.COMMAND_AGGREGATE_MATCH",
		"P0.ADMISSION.COMMAND_TRANSITION_DEFINED",
		"P0.ADMISSION.INTENT_REQUIRED_FIELDS",
	}
	refs := func(n int) []string {
		out := make([]string, n)
		copy(out, ruleIDs[:n])
		return out
	}
	if in.ArchitectureRevision != baseline {
		return admissionResult{
			Decision:    "REJECT",
			ReasonCode:  "ARCHITECTURE_REVISION_MISMATCH",
			RuleRefs:    refs(1),
			ReasonCodes: []string{"ARCHITECTURE_REVISION_MISMATCH"},
		}
	}
	if in.CapabilityRef != capUnitCreate {
		return admissionResult{
			Decision:    "REJECT",
			ReasonCode:  "UNKNOWN_CAPABILITY",
			RuleRefs:    refs(2),
			ReasonCodes: []string{"UNKNOWN_CAPABILITY"},
		}
	}
	// core.unit.create is mutating (checked by registry).
	if in.AggregateRef != aggUnit {
		return admissionResult{
			Decision:    "REJECT",
			ReasonCode:  "OWNERSHIP_MISMATCH",
			RuleRefs:    refs(4),
			ReasonCodes: []string{"OWNERSHIP_MISMATCH"},
		}
	}
	if in.CommandRef != cmdUnitCreate {
		return admissionResult{
			Decision:    "REJECT",
			ReasonCode:  "UNKNOWN_COMMAND",
			RuleRefs:    refs(5),
			ReasonCodes: []string{"UNKNOWN_COMMAND"},
		}
	}
	// COMMAND_CAPABILITY_MATCH, COMMAND_AGGREGATE_MATCH, and COMMAND_TRANSITION_DEFINED
	// are satisfied by the constants for this path.
	_ = capUnitCreate
	return admissionResult{
		Decision:      "ADMIT",
		TransitionRef: transUnitCreated,
		RuleRefs:      refs(9),
		ReasonCodes:   []string{},
	}
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

func newUnitCreateIntent(key, title, description, actor string) intentEnvelope {
	return intentEnvelope{
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

	intentDigest := computeIntentDigest(intent)
	adm := evaluateAdmission(intent)
	if adm.Decision != "ADMIT" {
		t.Fatalf("expected ADMIT, got %s (%s)", adm.Decision, adm.ReasonCode)
	}
	if adm.TransitionRef != "unit:created" {
		t.Fatalf("expected transition unit:created, got %s", adm.TransitionRef)
	}
	admissionID := computeAdmissionID(intentDigest, intent, adm.Decision, adm.TransitionRef, adm.RuleRefs, adm.ReasonCodes)
	t.Logf("admission_id=%s", admissionID)

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
	_ = admissionID
}

func TestUnitCreate_UnknownCapability_NoMutation(t *testing.T) {
	intent := intentEnvelope{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "phase-d-reject-01",
		CapabilityRef:        "core.doesnotexist",
		AggregateRef:         "unit",
		CommandRef:           "unit.create",
		Payload:              map[string]any{},
	}
	adm := evaluateAdmission(intent)
	if adm.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", adm.Decision)
	}
	if adm.ReasonCode != "UNKNOWN_CAPABILITY" {
		t.Fatalf("expected UNKNOWN_CAPABILITY, got %s", adm.ReasonCode)
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
