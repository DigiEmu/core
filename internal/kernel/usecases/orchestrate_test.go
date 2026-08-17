// Phase E Increment 1 — test-only admission gate probe.
// Proves that REJECT prevents the supplied execution closure from running.
// This is a test-only helper; it is not production orchestration.
package usecases_test

import (
	"testing"

	"digiemu-core/internal/admission"
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
func (p *admissionGateProbe) evaluate(eng *admission.Engine, intent admission.Intent, onAdmit func()) (admission.Result, error) {
	adm, err := eng.Evaluate(intent)
	if err != nil {
		return adm, err
	}
	if adm.Decision != "ADMIT" {
		return adm, nil
	}
	p.handlerCalled = true
	if onAdmit != nil {
		onAdmit()
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

	onAdmit := func() {
		// If the gate helper incorrectly invokes the closure on REJECT, fail
		// immediately so the test cannot pass.
		t.Fatalf("execution closure was invoked after Admission %s", "REJECT")
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
