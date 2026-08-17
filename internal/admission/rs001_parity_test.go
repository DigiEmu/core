package admission

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestRS001_Parity(t *testing.T) {
	casesDir := filepath.Join("..", "..", "testdata", "rs_001")
	entries, err := os.ReadDir(casesDir)
	if err != nil {
		t.Fatalf("ReadDir %s: %v", casesDir, err)
	}
	eng := NewEngine(V01Registry())
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "_schema_validation" {
			continue
		}
		caseDir := filepath.Join(casesDir, e.Name())
		inputPath := filepath.Join(caseDir, "input.json")
		expectedPath := filepath.Join(caseDir, "expected.json")
		if _, err := os.Stat(inputPath); err != nil {
			continue
		}
		if _, err := os.Stat(expectedPath); err != nil {
			continue
		}

		inputJSON, err := os.ReadFile(inputPath)
		if err != nil {
			t.Fatalf("%s: read input: %v", e.Name(), err)
		}
		expectedJSON, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Fatalf("%s: read expected: %v", e.Name(), err)
		}

		var in Intent
		if err := json.Unmarshal(inputJSON, &in); err != nil {
			t.Fatalf("%s: unmarshal input: %v", e.Name(), err)
		}

		var expected struct {
			Decision    string   `json:"decision"`
			ReasonCodes []string `json:"reason_codes"`
		}
		if err := json.Unmarshal(expectedJSON, &expected); err != nil {
			t.Fatalf("%s: unmarshal expected: %v", e.Name(), err)
		}

		res, err := eng.Evaluate(in)
		if err != nil {
			t.Fatalf("%s: Evaluate: %v", e.Name(), err)
		}

		if res.Decision != expected.Decision {
			t.Errorf("%s: decision got %s want %s", e.Name(), res.Decision, expected.Decision)
		}
		if !slicesEqual(res.ReasonCodes, expected.ReasonCodes) {
			t.Errorf("%s: reason_codes got %v want %v", e.Name(), res.ReasonCodes, expected.ReasonCodes)
		}
		if expected.Decision == "ADMIT" && res.TransitionRef != "unit:created" {
			t.Errorf("%s: expected transition_ref unit:created, got %s", e.Name(), res.TransitionRef)
		}
		if expected.Decision == "REJECT" && res.TransitionRef != "" {
			t.Errorf("%s: expected empty transition_ref for REJECT, got %s", e.Name(), res.TransitionRef)
		}
		// capability_ref / aggregate_ref / command_ref must be echoed.
		if res.CapabilityRef != in.CapabilityRef {
			t.Errorf("%s: capability_ref mismatch", e.Name())
		}
		if res.AggregateRef != in.AggregateRef {
			t.Errorf("%s: aggregate_ref mismatch", e.Name())
		}
		if res.CommandRef != in.CommandRef {
			t.Errorf("%s: command_ref mismatch", e.Name())
		}
		if res.AdmissionID == "" {
			t.Errorf("%s: admission_id must not be empty", e.Name())
		}
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	aa := make([]string, len(a))
	copy(aa, a)
	bb := make([]string, len(b))
	copy(bb, b)
	sort.Strings(aa)
	sort.Strings(bb)
	for i := range aa {
		if aa[i] != bb[i] {
			return false
		}
	}
	return true
}
