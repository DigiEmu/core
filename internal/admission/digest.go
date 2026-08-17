package admission

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"digiemu-core/internal/canonicaljson"
)

// intentDigestPreimage is the canonical input for P0.ADMISSION.INTENT.v0.1.
// Field order is fixed and matches docs/ADMISSION_ID_SPEC_v0.1.md §3.2.
type intentDigestPreimage struct {
	IntentDigestProfile  string         `json:"intent_digest_profile"`
	SchemaVersion        string         `json:"schema_version"`
	ArchitectureRevision string         `json:"architecture_revision"`
	CapabilityRef        string         `json:"capability_ref"`
	AggregateRef         string         `json:"aggregate_ref"`
	CommandRef           string         `json:"command_ref"`
	Payload              map[string]any `json:"payload"`
}

// ComputeIntentDigest computes the P0.ADMISSION.INTENT.v0.1 digest.
// It excludes intent_id, sorts object keys recursively within payload,
// preserves array order and scalar types, and produces
// "p0-intent:sha256:<lowercase-hex>".
func ComputeIntentDigest(in Intent) (string, error) {
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
		return "", err
	}
	sum := sha256.Sum256(b)
	return "p0-intent:sha256:" + hex.EncodeToString(sum[:]), nil
}

// admissionIDPreimage is the canonical input for P0.ADMISSION.ID.v0.1.
// Field order is fixed and matches docs/ADMISSION_ID_SPEC_v0.1.md §4.2.
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

// ComputeAdmissionID computes the P0.ADMISSION.ID.v0.1 identifier.
// architectureRevision is the architecture baseline revision recorded in the
// result, which may differ from the intent's declared revision (e.g., mismatch).
// rule_refs and reason_codes are sorted lexicographically before hashing.
// For unresolved REJECT, transitionRef "" is represented as JSON null.
// The output is "admission:sha256:<lowercase-hex>".
func ComputeAdmissionID(intentDigest string, in Intent, architectureRevision, decision, transitionRef string, ruleRefs, reasonCodes []string) (string, error) {
	rr := make([]string, len(ruleRefs))
	copy(rr, ruleRefs)
	sort.Strings(rr)
	rc := make([]string, len(reasonCodes))
	copy(rc, reasonCodes)
	sort.Strings(rc)

	var t *string
	if transitionRef != "" {
		t = &transitionRef
	}

	p := admissionIDPreimage{
		AdmissionIDProfile:   "P0.ADMISSION.ID.v0.1",
		SchemaVersion:        "v0.1",
		ArchitectureRevision: architectureRevision,
		IntentDigest:         intentDigest,
		CapabilityRef:        in.CapabilityRef,
		AggregateRef:         in.AggregateRef,
		CommandRef:           in.CommandRef,
		TransitionRef:        t,
		Decision:             decision,
		RuleRefs:             rr,
		ReasonCodes:          rc,
	}
	b, err := canonicaljson.Marshal(p)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "admission:sha256:" + hex.EncodeToString(sum[:]), nil
}
