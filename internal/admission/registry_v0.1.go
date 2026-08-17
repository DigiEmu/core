package admission

// V01Registry returns a typed Config representing the committed v0.1 P0
// architecture baseline and registries. It does NOT load YAML or access the
// filesystem; the values are compiled from the normative registry files.
func V01Registry() Config {
	return Config{
		ArchitectureRevision: "0.3",
		Capabilities: map[string]Capability{
			"core.verify":           {ID: "core.verify", Mutation: false},
			"core.replay":           {ID: "core.replay", Mutation: false},
			"core.snapshot":         {ID: "core.snapshot", Mutation: false},
			"core.conformance":      {ID: "core.conformance", Mutation: false},
			"core.audit":            {ID: "core.audit", Mutation: false},
			"core.unit.create":      {ID: "core.unit.create", Mutation: true},
			"core.version.create":   {ID: "core.version.create", Mutation: true},
			"core.meaning.set":      {ID: "core.meaning.set", Mutation: true},
			"core.claim.set":        {ID: "core.claim.set", Mutation: true},
			"core.uncertainty.set":  {ID: "core.uncertainty.set", Mutation: true},
			"core.import":           {ID: "core.import", Mutation: false},
			"p0.admission":          {ID: "p0.admission", Mutation: false},
		},
		Ownership: map[string][]string{
			"unit": {
				"core.unit.create",
				"core.version.create",
				"core.meaning.set",
				"core.claim.set",
				"core.uncertainty.set",
			},
		},
		Commands: map[string]Command{
			"unit.create": {
				ID:           "unit.create",
				CapabilityID: "core.unit.create",
				AggregateID:  "unit",
				TransitionID: "unit:created",
			},
			"version.create": {
				ID:           "version.create",
				CapabilityID: "core.version.create",
				AggregateID:  "unit",
				TransitionID: "version:created",
			},
			"meaning.set": {
				ID:           "meaning.set",
				CapabilityID: "core.meaning.set",
				AggregateID:  "unit",
				TransitionID: "meaning:set",
			},
			"claim.set": {
				ID:           "claim.set",
				CapabilityID: "core.claim.set",
				AggregateID:  "unit",
				TransitionID: "claim:set",
			},
			"uncertainty.set": {
				ID:           "uncertainty.set",
				CapabilityID: "core.uncertainty.set",
				AggregateID:  "unit",
				TransitionID: "uncertainty:set",
			},
		},
		Rules: []Rule{
			{ID: "P0.ADMISSION.ARCHITECTURE_REVISION", FailureReasonCode: "ARCHITECTURE_REVISION_MISMATCH"},
			{ID: "P0.ADMISSION.CAPABILITY_EXISTS", FailureReasonCode: "UNKNOWN_CAPABILITY"},
			{ID: "P0.ADMISSION.CAPABILITY_MUTATES", FailureReasonCode: "CAPABILITY_NOT_MUTATING"},
			{ID: "P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY", FailureReasonCode: "OWNERSHIP_MISMATCH"},
			{ID: "P0.ADMISSION.COMMAND_EXISTS", FailureReasonCode: "UNKNOWN_COMMAND"},
			{ID: "P0.ADMISSION.COMMAND_CAPABILITY_MATCH", FailureReasonCode: "COMMAND_CAPABILITY_MISMATCH"},
			{ID: "P0.ADMISSION.COMMAND_AGGREGATE_MATCH", FailureReasonCode: "COMMAND_AGGREGATE_MISMATCH"},
			{ID: "P0.ADMISSION.COMMAND_TRANSITION_DEFINED", FailureReasonCode: "UNDEFINED_TRANSITION"},
			{ID: "P0.ADMISSION.INTENT_REQUIRED_FIELDS", FailureReasonCode: "MISSING_REQUIRED_FIELD"},
		},
	}
}
