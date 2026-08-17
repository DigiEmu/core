package admission

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestV01Registry_Parity verifies that V01Registry() stays in sync with the
// normative committed YAML registries. The YAML files are the source of truth;
// V01Registry() is the compiled Go snapshot used by the Engine.
//
// Compared fields:
//
//	architecture-baseline.yaml:     baseline.revision
//	core-capability-registry.yaml:  capability id and mutation flag
//	aggregate-ownership-registry.yaml: aggregate id and owned capability ids
//	command-event-catalogue.yaml:   command id, capability id, aggregate id, transition id
//	admission-rule-registry.yaml:   executable rule id, order, and failure reason code
//
// Intentionally excluded (not represented in admission.Config):
//
//	descriptions, titles, status text, source comments, handler metadata,
//	event ids, invariant refs, deterministic_inputs metadata.
func TestV01Registry_Parity(t *testing.T) {
	root := repoRoot()
	reg := V01Registry()

	assertArchitectureBaseline(t, root, reg)
	assertCapabilities(t, root, reg)
	assertOwnership(t, root, reg)
	assertCommands(t, root, reg)
	assertRules(t, root, reg)
}

func repoRoot() string {
	return filepath.Join("..", "..")
}

func loadYAML(t *testing.T, rel string, v any) {
	t.Helper()
	p := filepath.Join(repoRoot(), rel)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if err := yaml.Unmarshal(b, v); err != nil {
		t.Fatalf("unmarshal %s: %v", rel, err)
	}
}

// architecture-baseline.yaml parity.
func assertArchitectureBaseline(t *testing.T, root string, reg Config) {
	var doc struct {
		Baseline struct {
			Revision string `yaml:"revision"`
		} `yaml:"baseline"`
	}
	loadYAML(t, "architecture-baseline.yaml", &doc)
	if reg.ArchitectureRevision != doc.Baseline.Revision {
		t.Fatalf("architecture revision mismatch: yaml=%q go=%q", doc.Baseline.Revision, reg.ArchitectureRevision)
	}
}

// core-capability-registry.yaml parity.
func assertCapabilities(t *testing.T, root string, reg Config) {
	var doc struct {
		Capabilities []struct {
			ID       string `yaml:"id"`
			Mutation bool   `yaml:"mutation"`
		} `yaml:"capabilities"`
	}
	loadYAML(t, "core-capability-registry.yaml", &doc)

	yamlCaps := make(map[string]bool, len(doc.Capabilities))
	for _, c := range doc.Capabilities {
		yamlCaps[c.ID] = c.Mutation
	}

	if len(reg.Capabilities) != len(yamlCaps) {
		t.Fatalf("capability count mismatch: yaml=%d go=%d", len(yamlCaps), len(reg.Capabilities))
	}
	for id, mut := range yamlCaps {
		goCap, ok := reg.Capabilities[id]
		if !ok {
			t.Fatalf("capability %s in YAML but missing from V01Registry", id)
		}
		if goCap.Mutation != mut {
			t.Fatalf("capability %s mutation mismatch: yaml=%t go=%t", id, mut, goCap.Mutation)
		}
	}
	for id := range reg.Capabilities {
		if _, ok := yamlCaps[id]; !ok {
			t.Fatalf("capability %s in V01Registry but missing from YAML", id)
		}
	}
}

// aggregate-ownership-registry.yaml parity.
func assertOwnership(t *testing.T, root string, reg Config) {
	var doc struct {
		Ownership []struct {
			AggregateID       string   `yaml:"aggregate_id"`
			OwnedCapabilities []string `yaml:"owned_capabilities"`
		} `yaml:"ownership"`
	}
	loadYAML(t, "aggregate-ownership-registry.yaml", &doc)

	yamlOwner := make(map[string]map[string]bool, len(doc.Ownership))
	for _, o := range doc.Ownership {
		caps := make(map[string]bool, len(o.OwnedCapabilities))
		for _, c := range o.OwnedCapabilities {
			caps[c] = true
		}
		yamlOwner[o.AggregateID] = caps
	}

	if len(reg.Ownership) != len(yamlOwner) {
		t.Fatalf("ownership aggregate count mismatch: yaml=%d go=%d", len(yamlOwner), len(reg.Ownership))
	}
	for agg, goCaps := range reg.Ownership {
		yamlCaps, ok := yamlOwner[agg]
		if !ok {
			t.Fatalf("aggregate %s in V01Registry but missing from YAML", agg)
		}
		if len(goCaps) != len(yamlCaps) {
			t.Fatalf("aggregate %s owned capability count mismatch: yaml=%d go=%d", agg, len(yamlCaps), len(goCaps))
		}
		for _, c := range goCaps {
			if !yamlCaps[c] {
				t.Fatalf("aggregate %s capability %s in V01Registry but missing from YAML", agg, c)
			}
		}
		for c := range yamlCaps {
			found := false
			for _, gc := range goCaps {
				if gc == c {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("aggregate %s capability %s in YAML but missing from V01Registry", agg, c)
			}
		}
	}
	for agg := range yamlOwner {
		if _, ok := reg.Ownership[agg]; !ok {
			t.Fatalf("aggregate %s in YAML but missing from V01Registry", agg)
		}
	}
}

// command-event-catalogue.yaml parity.
func assertCommands(t *testing.T, root string, reg Config) {
	var doc struct {
		Commands []struct {
			ID           string `yaml:"command_id"`
			CapabilityID string `yaml:"capability_id"`
			AggregateID  string `yaml:"aggregate_id"`
			TransitionID string `yaml:"transition_id"`
		} `yaml:"commands"`
	}
	loadYAML(t, "command-event-catalogue.yaml", &doc)

	yamlCmds := make(map[string]Command, len(doc.Commands))
	for _, c := range doc.Commands {
		yamlCmds[c.ID] = Command{
			ID:           c.ID,
			CapabilityID: c.CapabilityID,
			AggregateID:  c.AggregateID,
			TransitionID: c.TransitionID,
		}
	}

	if len(reg.Commands) != len(yamlCmds) {
		t.Fatalf("command count mismatch: yaml=%d go=%d", len(yamlCmds), len(reg.Commands))
	}
	for id, yc := range yamlCmds {
		gc, ok := reg.Commands[id]
		if !ok {
			t.Fatalf("command %s in YAML but missing from V01Registry", id)
		}
		if gc.CapabilityID != yc.CapabilityID {
			t.Fatalf("command %s capability mismatch: yaml=%s go=%s", id, yc.CapabilityID, gc.CapabilityID)
		}
		if gc.AggregateID != yc.AggregateID {
			t.Fatalf("command %s aggregate mismatch: yaml=%s go=%s", id, yc.AggregateID, gc.AggregateID)
		}
		if gc.TransitionID != yc.TransitionID {
			t.Fatalf("command %s transition mismatch: yaml=%s go=%s", id, yc.TransitionID, gc.TransitionID)
		}
	}
	for id := range reg.Commands {
		if _, ok := yamlCmds[id]; !ok {
			t.Fatalf("command %s in V01Registry but missing from YAML", id)
		}
	}
}

// admission-rule-registry.yaml parity.
func assertRules(t *testing.T, root string, reg Config) {
	var doc struct {
		Rules []struct {
			RuleID            string `yaml:"rule_id"`
			Executable        bool   `yaml:"executable"`
			FailureReasonCode string `yaml:"failure_reason_code"`
		} `yaml:"admission_rules"`
	}
	loadYAML(t, "admission-rule-registry.yaml", &doc)

	var yamlRules []Rule
	for _, r := range doc.Rules {
		if r.Executable {
			yamlRules = append(yamlRules, Rule{ID: r.RuleID, FailureReasonCode: r.FailureReasonCode})
		}
	}

	if len(yamlRules) != len(reg.Rules) {
		t.Fatalf("executable rule count mismatch: yaml=%d go=%d", len(yamlRules), len(reg.Rules))
	}
	for i := range yamlRules {
		yr := yamlRules[i]
		gr := reg.Rules[i]
		if gr.ID != yr.ID {
			t.Fatalf("rule %d id mismatch: yaml=%s go=%s", i, yr.ID, gr.ID)
		}
		if gr.FailureReasonCode != yr.FailureReasonCode {
			t.Fatalf("rule %s reason code mismatch: yaml=%s go=%s", gr.ID, yr.FailureReasonCode, gr.FailureReasonCode)
		}
	}

	// Confirm the expected normative order is still present (human-readable guard).
	wantOrder := []string{
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
	gotOrder := make([]string, len(reg.Rules))
	for i, r := range reg.Rules {
		gotOrder[i] = r.ID
	}
	if !slicesEqual(gotOrder, wantOrder) {
		t.Fatalf("rule order mismatch:\nyaml/exec order = %v\ngo order        = %v", gotOrder, gotOrder)
	}
}
