package admission

import (
	"reflect"
	"sort"
	"testing"
)

func mustAdmit(t *testing.T, in Intent) Result {
	t.Helper()
	eng := NewEngine(V01Registry())
	res, err := eng.Evaluate(in)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	return res
}

func requireEqualSet(t *testing.T, got, want []string, label string) {
	t.Helper()
	g := make([]string, len(got))
	copy(g, got)
	sort.Strings(g)
	w := make([]string, len(want))
	copy(w, want)
	sort.Strings(w)
	if !reflect.DeepEqual(g, w) {
		t.Errorf("%s: got %v, want %v", label, got, want)
	}
}

func baseIntent() Intent {
	return Intent{
		SchemaVersion:        "v0.1",
		ArchitectureRevision: "0.3",
		IntentID:             "test-intent",
		CapabilityRef:        "core.unit.create",
		AggregateRef:         "unit",
		CommandRef:           "unit.create",
		Payload:              map[string]any{},
	}
}

func TestEngine_ValidAdmit(t *testing.T) {
	in := baseIntent()
	in.Payload = map[string]any{
		"key":         "example-unit",
		"title":       "Example Unit",
		"description": "An example unit.",
	}
	res := mustAdmit(t, in)
	if res.Decision != "ADMIT" {
		t.Fatalf("expected ADMIT, got %s", res.Decision)
	}
	if res.TransitionRef != "unit:created" {
		t.Fatalf("expected transition unit:created, got %s", res.TransitionRef)
	}
	if len(res.ReasonCodes) != 0 {
		t.Fatalf("expected no reason codes, got %v", res.ReasonCodes)
	}
	if res.AdmissionID == "" {
		t.Fatalf("admission_id must not be empty")
	}
}

func TestEngine_ArchitectureMismatch(t *testing.T) {
	in := baseIntent()
	in.ArchitectureRevision = "0.4"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"ARCHITECTURE_REVISION_MISMATCH"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION"}, "rule_refs")
	if res.TransitionRef != "" {
		t.Fatalf("expected no transition ref for REJECT")
	}
	if res.AdmissionID == "" {
		t.Fatalf("admission_id must not be empty")
	}
}

func TestEngine_UnknownCapability(t *testing.T) {
	in := baseIntent()
	in.CapabilityRef = "core.unknown.capability"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"UNKNOWN_CAPABILITY"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION", "P0.ADMISSION.CAPABILITY_EXISTS"}, "rule_refs")
}

func TestEngine_CapabilityNotMutating(t *testing.T) {
	in := baseIntent()
	in.CapabilityRef = "core.verify"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"CAPABILITY_NOT_MUTATING"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION", "P0.ADMISSION.CAPABILITY_EXISTS", "P0.ADMISSION.CAPABILITY_MUTATES"}, "rule_refs")
}

func TestEngine_OwnershipMismatch(t *testing.T) {
	in := baseIntent()
	in.AggregateRef = "unknown-aggregate"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"OWNERSHIP_MISMATCH"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION", "P0.ADMISSION.CAPABILITY_EXISTS", "P0.ADMISSION.CAPABILITY_MUTATES", "P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY"}, "rule_refs")
}

func TestEngine_UnknownCommand(t *testing.T) {
	in := baseIntent()
	in.CommandRef = "unit.missing"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"UNKNOWN_COMMAND"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION", "P0.ADMISSION.CAPABILITY_EXISTS", "P0.ADMISSION.CAPABILITY_MUTATES", "P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY", "P0.ADMISSION.COMMAND_EXISTS"}, "rule_refs")
}

func TestEngine_CommandCapabilityMismatch(t *testing.T) {
	in := baseIntent()
	in.CapabilityRef = "core.version.create"
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"COMMAND_CAPABILITY_MISMATCH"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.ARCHITECTURE_REVISION", "P0.ADMISSION.CAPABILITY_EXISTS", "P0.ADMISSION.CAPABILITY_MUTATES", "P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY", "P0.ADMISSION.COMMAND_EXISTS", "P0.ADMISSION.COMMAND_CAPABILITY_MATCH"}, "rule_refs")
}

func TestEngine_CommandAggregateMismatch(t *testing.T) {
	cfg := V01Registry()
	// Add an aggregate 'other' that also owns core.unit.create,
	// and a command other.create with that capability but aggregate 'other'.
	cfg.Ownership["other"] = []string{"core.unit.create"}
	cfg.Commands["other.create"] = Command{
		ID:           "other.create",
		CapabilityID: "core.unit.create",
		AggregateID:  "other",
		TransitionID: "other:created",
	}
	eng := NewEngine(cfg)
	in := baseIntent()
	in.AggregateRef = "unit"
	in.CommandRef = "other.create"
	res, err := eng.Evaluate(in)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"COMMAND_AGGREGATE_MISMATCH"}, "reason_codes")
}

func TestEngine_UndefinedTransition(t *testing.T) {
	cfg := V01Registry()
	cfg.Commands["unit.notransition"] = Command{
		ID:           "unit.notransition",
		CapabilityID: "core.unit.create",
		AggregateID:  "unit",
		TransitionID: "",
	}
	eng := NewEngine(cfg)
	in := baseIntent()
	in.CommandRef = "unit.notransition"
	res, err := eng.Evaluate(in)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"UNDEFINED_TRANSITION"}, "reason_codes")
}

func TestEngine_MissingRequiredField(t *testing.T) {
	in := baseIntent()
	in.CommandRef = ""
	res := mustAdmit(t, in)
	if res.Decision != "REJECT" {
		t.Fatalf("expected REJECT, got %s", res.Decision)
	}
	requireEqualSet(t, res.ReasonCodes, []string{"MISSING_REQUIRED_FIELD"}, "reason_codes")
	requireEqualSet(t, res.RuleRefs, []string{"P0.ADMISSION.INTENT_REQUIRED_FIELDS"}, "rule_refs")
	if res.TransitionRef != "" {
		t.Fatalf("expected no transition ref")
	}
}

func TestEngine_NoInvalidIntentReasonCode(t *testing.T) {
	// Malformed intent is handled by the schema check wrapper in callers, not the Engine.
	// Inside the Engine only normative reason codes are returned.
	in := baseIntent()
	in.CommandRef = ""
	res := mustAdmit(t, in)
	for _, c := range res.ReasonCodes {
		if c == "INVALID_INTENT" {
			t.Fatalf("INVALID_INTENT must not appear in Engine result, got %v", res.ReasonCodes)
		}
	}
}

func TestEngine_Determinism_SameInputSameID(t *testing.T) {
	in := baseIntent()
	res1 := mustAdmit(t, in)
	res2 := mustAdmit(t, in)
	if res1.AdmissionID != res2.AdmissionID {
		t.Fatalf("same input must produce same admission_id: %s vs %s", res1.AdmissionID, res2.AdmissionID)
	}
	d1, _ := ComputeIntentDigest(in)
	d2, _ := ComputeIntentDigest(in)
	if d1 != d2 {
		t.Fatalf("same input must produce same intent_digest: %s vs %s", d1, d2)
	}
}

func TestEngine_Determinism_PayloadAlpha(t *testing.T) {
	in := baseIntent()
	in.Payload = map[string]any{"extra": "x", "key": "alpha"}
	res := mustAdmit(t, in)
	gotDigest, _ := ComputeIntentDigest(in)
	wantDigest := "p0-intent:sha256:a5bc1e21afb34ed1b358229935ef1e4ca4c99cdf5f3638688a72e108853e0294"
	wantID := "admission:sha256:7f87382ff543c0b40539a05c70a2fad82fcd046a5989215296e8f8353f41ad4f"
	if gotDigest != wantDigest {
		t.Fatalf("intent digest mismatch\n got: %s\nwant: %s", gotDigest, wantDigest)
	}
	if res.AdmissionID != wantID {
		t.Fatalf("admission_id mismatch\n got: %s\nwant: %s", res.AdmissionID, wantID)
	}
}

func TestEngine_Determinism_PayloadBeta(t *testing.T) {
	in := baseIntent()
	in.Payload = map[string]any{"extra": "x", "key": "beta"}
	res := mustAdmit(t, in)
	gotDigest, _ := ComputeIntentDigest(in)
	wantDigest := "p0-intent:sha256:2aaefddc9882b0894854c30abc0bff11ed7bd1008de0c5018b0fcbc2057112a0"
	wantID := "admission:sha256:4ac7d1f1344b4621d6372a91ade7e81e46077ba6559a0f141e94e6a206bbd33c"
	if gotDigest != wantDigest {
		t.Fatalf("intent digest mismatch\n got: %s\nwant: %s", gotDigest, wantDigest)
	}
	if res.AdmissionID != wantID {
		t.Fatalf("admission_id mismatch\n got: %s\nwant: %s", res.AdmissionID, wantID)
	}
}

func TestEngine_Determinism_PayloadKeyOrder(t *testing.T) {
	alpha1 := baseIntent()
	alpha1.Payload = map[string]any{"extra": "x", "key": "alpha"}
	alpha2 := baseIntent()
	alpha2.Payload = map[string]any{"key": "alpha", "extra": "x"}
	d1, _ := ComputeIntentDigest(alpha1)
	d2, _ := ComputeIntentDigest(alpha2)
	res1 := mustAdmit(t, alpha1)
	res2 := mustAdmit(t, alpha2)
	if d1 != d2 {
		t.Fatalf("payload key order must not affect digest: %s vs %s", d1, d2)
	}
	if res1.AdmissionID != res2.AdmissionID {
		t.Fatalf("payload key order must not affect admission_id: %s vs %s", res1.AdmissionID, res2.AdmissionID)
	}
	want := "admission:sha256:7f87382ff543c0b40539a05c70a2fad82fcd046a5989215296e8f8353f41ad4f"
	if res1.AdmissionID != want {
		t.Fatalf("admission_id mismatch\n got: %s\nwant: %s", res1.AdmissionID, want)
	}
}

func TestEngine_Determinism_RuleRefsOrder(t *testing.T) {
	// ComputeAdmissionID must sort rule_refs before hashing.
	in := baseIntent()
	intentDigest, err := ComputeIntentDigest(in)
	if err != nil {
		t.Fatalf("ComputeIntentDigest: %v", err)
	}
	id1, err := ComputeAdmissionID(intentDigest, in, in.ArchitectureRevision, "ADMIT", "unit:created", []string{"z", "a"}, []string{})
	if err != nil {
		t.Fatalf("ComputeAdmissionID: %v", err)
	}
	id2, err := ComputeAdmissionID(intentDigest, in, in.ArchitectureRevision, "ADMIT", "unit:created", []string{"a", "z"}, []string{})
	if err != nil {
		t.Fatalf("ComputeAdmissionID: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("rule_refs order must not affect admission_id: %s vs %s", id1, id2)
	}
}

func TestEngine_Determinism_ReasonCodesOrder(t *testing.T) {
	in := baseIntent()
	in.ArchitectureRevision = "0.4"
	intentDigest, err := ComputeIntentDigest(in)
	if err != nil {
		t.Fatalf("ComputeIntentDigest: %v", err)
	}
	id1, err := ComputeAdmissionID(intentDigest, in, in.ArchitectureRevision, "REJECT", "", []string{"P0.ADMISSION.ARCHITECTURE_REVISION"}, []string{"z", "a"})
	if err != nil {
		t.Fatalf("ComputeAdmissionID: %v", err)
	}
	id2, err := ComputeAdmissionID(intentDigest, in, in.ArchitectureRevision, "REJECT", "", []string{"P0.ADMISSION.ARCHITECTURE_REVISION"}, []string{"a", "z"})
	if err != nil {
		t.Fatalf("ComputeAdmissionID: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("reason_codes order must not affect admission_id: %s vs %s", id1, id2)
	}
}

func TestEngine_Determinism_RejectStableNoTransition(t *testing.T) {
	in := baseIntent()
	in.CommandRef = "unit.missing"
	res := mustAdmit(t, in)
	if res.TransitionRef != "" {
		t.Fatalf("REJECT before transition resolution must not have transition_ref")
	}
	if res.AdmissionID == "" {
		t.Fatalf("admission_id must be present for REJECT")
	}
	// Recompute to prove stability.
	res2 := mustAdmit(t, in)
	if res.AdmissionID != res2.AdmissionID {
		t.Fatalf("REJECT admission_id must be stable: %s vs %s", res.AdmissionID, res2.AdmissionID)
	}
}

func TestEngine_AdmitResolvesTransition(t *testing.T) {
	in := baseIntent()
	in.Payload = map[string]any{"key": "k", "title": "t"}
	res := mustAdmit(t, in)
	if res.TransitionRef != "unit:created" {
		t.Fatalf("ADMIT must resolve transition_ref, got %s", res.TransitionRef)
	}
}

func TestEngine_ResultContainsOnlyNormativeRefs(t *testing.T) {
	in := baseIntent()
	in.Payload = map[string]any{"key": "k", "title": "t"}
	res := mustAdmit(t, in)
	normative := map[string]bool{
		"P0.ADMISSION.ARCHITECTURE_REVISION":      true,
		"P0.ADMISSION.CAPABILITY_EXISTS":          true,
		"P0.ADMISSION.CAPABILITY_MUTATES":         true,
		"P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY":  true,
		"P0.ADMISSION.COMMAND_EXISTS":             true,
		"P0.ADMISSION.COMMAND_CAPABILITY_MATCH":   true,
		"P0.ADMISSION.COMMAND_AGGREGATE_MATCH":    true,
		"P0.ADMISSION.COMMAND_TRANSITION_DEFINED": true,
		"P0.ADMISSION.INTENT_REQUIRED_FIELDS":     true,
	}
	for _, r := range res.RuleRefs {
		if !normative[r] {
			t.Fatalf("non-normative rule_ref %q", r)
		}
	}
	normativeReason := map[string]bool{
		"ARCHITECTURE_REVISION_MISMATCH": true,
		"UNKNOWN_CAPABILITY":             true,
		"CAPABILITY_NOT_MUTATING":        true,
		"OWNERSHIP_MISMATCH":             true,
		"UNKNOWN_COMMAND":                true,
		"COMMAND_CAPABILITY_MISMATCH":    true,
		"COMMAND_AGGREGATE_MISMATCH":     true,
		"UNDEFINED_TRANSITION":           true,
		"MISSING_REQUIRED_FIELD":         true,
	}
	for _, c := range res.ReasonCodes {
		if !normativeReason[c] {
			t.Fatalf("non-normative reason_code %q", c)
		}
	}
}
