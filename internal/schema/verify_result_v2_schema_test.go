package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

func TestVerifyResultV2ExamplesValidateSchema(t *testing.T) {
	// Resolve paths relative to repository root (two levels up from this package)
	repoRoot := filepath.Join("..", "..")
	// Locate schema file
	schemaPath := filepath.Join(repoRoot, "schemas", "verify_result_v2.schema.json")
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatalf("failed to resolve absolute path for schema %s: %v", schemaPath, err)
	}

	compiler := jsonschema.NewCompiler()
	// Use file URL with triple slash for absolute paths on Windows
	schemaURL := "file:///" + filepath.ToSlash(absSchemaPath)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("failed to compile schema %s: %v", schemaURL, err)
	}

	// Example files to validate
	examples := []string{
		filepath.Join(repoRoot, "testdata", "verify_result_v2", "pass_state_reconstructed.json"),
		filepath.Join(repoRoot, "testdata", "verify_result_v2", "fail_hash_mismatch.json"),
		filepath.Join(repoRoot, "testdata", "verify_result_v2", "error_invalid_snapshot_schema.json"),

		// Conformance expected results
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "basic_pass", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "hash_mismatch_fail", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "invalid_schema_error", "expected_verify_result.json"),
	}

	for _, rel := range examples {
		data, err := os.ReadFile(rel)
		if err != nil {
			t.Fatalf("failed to read example file %s: %v", rel, err)
		}
		var v interface{}
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("invalid JSON in example file %s: %v", rel, err)
		}
		if err := sch.Validate(v); err != nil {
			t.Fatalf("schema validation failed for %s: %v", rel, err)
		}
	}
}
