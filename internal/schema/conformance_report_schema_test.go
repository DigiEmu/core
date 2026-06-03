package schema_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

func TestConformanceReportExampleAndCLIValidateSchema(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	schemaPath := filepath.Join(repoRoot, "schemas", "core_2_conformance_report.schema.json")
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatalf("abs schema: %v", err)
	}

	compiler := jsonschema.NewCompiler()
	schemaURL := "file:///" + filepath.ToSlash(absSchemaPath)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema %s: %v", schemaURL, err)
	}

	// Validate example file
	example := filepath.Join(repoRoot, "testdata", "core_2_conformance", "report_expected_basic.json")
	data, err := os.ReadFile(example)
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("invalid example json: %v", err)
	}
	if err := sch.Validate(v); err != nil {
		t.Fatalf("schema validation failed for example: %v", err)
	}

	// Run CLI to produce JSON output and validate
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "./cmd/digiemu", "experimental", "conformance", "run", filepath.Join("testdata", "core_2_conformance"), "--json")
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cli run failed: %v\noutput:\n%s", err, string(out))
	}
	s := string(out)
	// find first JSON object
	i := -1
	for idx, ch := range s {
		if ch == '{' {
			i = idx
			break
		}
	}
	if i < 0 {
		t.Fatalf("no JSON found in CLI output:\n%s", s)
	}
	var rep interface{}
	if err := json.Unmarshal([]byte(s[i:]), &rep); err != nil {
		t.Fatalf("invalid CLI JSON: %v\njson fragment:\n%s", err, s[i:])
	}
	if err := sch.Validate(rep); err != nil {
		t.Fatalf("schema validation failed for CLI output: %v", err)
	}

	// Basic structural assertions
	m, ok := rep.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected CLI JSON shape")
	}
	if m["report_version"] == nil {
		t.Fatalf("missing report_version in CLI JSON")
	}
	if m["status"] != "PASS" {
		t.Fatalf("expected status PASS, got %v", m["status"])
	}
	if int(m["total"].(float64)) != 11 {
		t.Fatalf("expected total 11, got %v", m["total"])
	}
}
