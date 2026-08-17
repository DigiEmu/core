package admission

import (
	"sort"
)

// Engine evaluates P0 architectural Admission without depending on filesystem,
// YAML, CLI, HTTP, or Core state identity.
type Engine struct {
	config Config
}

// NewEngine returns an Admission Engine bound to the provided typed Config.
func NewEngine(c Config) *Engine {
	return &Engine{config: c}
}

// Evaluate runs the ordered Admission rules against the Intent.
// It returns a Result that satisfies schemas/admission_result_v0.2.schema.json.
func (e *Engine) Evaluate(in Intent) (Result, error) {
	var transitionRef string
	ruleRefs := make([]string, 0, len(e.config.Rules))
	reasonCodes := make([]string, 0, 1)

	// P0.ADMISSION.INTENT_REQUIRED_FIELDS must be satisfied before any
	// command/capability/aggregate lookup. The rule is also present in the
	// ordered rule list so that valid intents record it as evaluated.
	if !e.intentRequiredFields(in) {
		return e.buildResult(in, "REJECT", "", []string{"P0.ADMISSION.INTENT_REQUIRED_FIELDS"}, []string{"MISSING_REQUIRED_FIELD"})
	}

	for _, rule := range e.config.Rules {
		passed, tref := e.checkRule(rule.ID, in)
		ruleRefs = append(ruleRefs, rule.ID)
		if !passed {
			reasonCodes = append(reasonCodes, rule.FailureReasonCode)
			return e.buildResult(in, "REJECT", "", ruleRefs, reasonCodes)
		}
		if tref != "" {
			transitionRef = tref
		}
	}

	return e.buildResult(in, "ADMIT", transitionRef, ruleRefs, reasonCodes)
}

func (e *Engine) buildResult(in Intent, decision, transitionRef string, ruleRefs, reasonCodes []string) (Result, error) {
	intentDigest, err := ComputeIntentDigest(in)
	if err != nil {
		return Result{}, err
	}
	admissionID, err := ComputeAdmissionID(intentDigest, in, e.config.ArchitectureRevision, decision, transitionRef, ruleRefs, reasonCodes)
	if err != nil {
		return Result{}, err
	}

	rr := make([]string, len(ruleRefs))
	copy(rr, ruleRefs)
	sort.Strings(rr)
	rc := make([]string, len(reasonCodes))
	copy(rc, reasonCodes)
	sort.Strings(rc)

	r := Result{
		SchemaVersion:        "v0.2",
		ArchitectureRevision: e.config.ArchitectureRevision,
		AdmissionID:          admissionID,
		Decision:             decision,
		CapabilityRef:        in.CapabilityRef,
		AggregateRef:         in.AggregateRef,
		CommandRef:           in.CommandRef,
		RuleRefs:             rr,
		ReasonCodes:          rc,
	}
	if decision == "ADMIT" {
		r.TransitionRef = transitionRef
	}
	return r, nil
}

func (e *Engine) intentRequiredFields(in Intent) bool {
	if in.SchemaVersion == "" {
		return false
	}
	if in.ArchitectureRevision == "" {
		return false
	}
	if in.CapabilityRef == "" {
		return false
	}
	if in.AggregateRef == "" {
		return false
	}
	if in.CommandRef == "" {
		return false
	}
	if in.Payload == nil {
		return false
	}
	return true
}

func (e *Engine) checkRule(id string, in Intent) (pass bool, transitionRef string) {
	switch id {
	case "P0.ADMISSION.ARCHITECTURE_REVISION":
		return in.ArchitectureRevision == e.config.ArchitectureRevision, ""
	case "P0.ADMISSION.CAPABILITY_EXISTS":
		_, ok := e.config.Capabilities[in.CapabilityRef]
		return ok, ""
	case "P0.ADMISSION.CAPABILITY_MUTATES":
		cap, ok := e.config.Capabilities[in.CapabilityRef]
		if !ok {
			return false, ""
		}
		return cap.Mutation, ""
	case "P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY":
		caps, ok := e.config.Ownership[in.AggregateRef]
		if !ok {
			return false, ""
		}
		for _, c := range caps {
			if c == in.CapabilityRef {
				return true, ""
			}
		}
		return false, ""
	case "P0.ADMISSION.COMMAND_EXISTS":
		_, ok := e.config.Commands[in.CommandRef]
		if !ok {
			return false, ""
		}
		return true, ""
	case "P0.ADMISSION.COMMAND_CAPABILITY_MATCH":
		cmd, ok := e.config.Commands[in.CommandRef]
		if !ok {
			return false, ""
		}
		return cmd.CapabilityID == in.CapabilityRef, ""
	case "P0.ADMISSION.COMMAND_AGGREGATE_MATCH":
		cmd, ok := e.config.Commands[in.CommandRef]
		if !ok {
			return false, ""
		}
		return cmd.AggregateID == in.AggregateRef, ""
	case "P0.ADMISSION.COMMAND_TRANSITION_DEFINED":
		cmd, ok := e.config.Commands[in.CommandRef]
		if !ok {
			return false, ""
		}
		return cmd.TransitionID != "", cmd.TransitionID
	case "P0.ADMISSION.INTENT_REQUIRED_FIELDS":
		return true, ""
	}
	return false, ""
}
