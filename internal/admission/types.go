// Package admission implements the P0 Architecture Constitution Admission Engine.
// It is an architectural admissibility layer and does not modify Core domain,
// canonicalization, state identity, or production dispatch.
package admission

// Intent represents the normative fields from schemas/intent-envelope.schema.json.
// intent_id is present for envelope completeness but is excluded from the
// deterministic P0 intent digest.
type Intent struct {
	SchemaVersion        string         `json:"schema_version"`
	ArchitectureRevision string         `json:"architecture_revision"`
	IntentID             string         `json:"intent_id"`
	CapabilityRef        string         `json:"capability_ref"`
	AggregateRef         string         `json:"aggregate_ref"`
	CommandRef           string         `json:"command_ref"`
	Payload              map[string]any `json:"payload"`
}

// Result represents schemas/admission_result_v0.2.schema.json.
type Result struct {
	SchemaVersion        string   `json:"schema_version"`
	ArchitectureRevision string   `json:"architecture_revision"`
	AdmissionID          string   `json:"admission_id"`
	Decision             string   `json:"decision"`
	CapabilityRef        string   `json:"capability_ref"`
	AggregateRef         string   `json:"aggregate_ref"`
	CommandRef           string   `json:"command_ref"`
	TransitionRef        string   `json:"transition_ref,omitempty"`
	RuleRefs             []string `json:"rule_refs"`
	ReasonCodes          []string `json:"reason_codes"`
}

// Capability is a registered Core or P0 capability.
type Capability struct {
	ID       string
	Mutation bool
}

// Command is a catalogued command.
type Command struct {
	ID           string
	CapabilityID string
	AggregateID  string
	TransitionID string
}

// Rule is an admission rule identifier and its normative failure reason code.
type Rule struct {
	ID                string
	FailureReasonCode string
}

// Config is the typed, immutable Admission configuration.
// It contains the architecture revision and lookup tables for capabilities,
// aggregate ownership, commands, and the ordered list of rules to evaluate.
type Config struct {
	ArchitectureRevision string
	Capabilities         map[string]Capability
	Ownership            map[string][]string
	Commands             map[string]Command
	Rules                []Rule
}
