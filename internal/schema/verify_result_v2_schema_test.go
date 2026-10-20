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
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "missing_required_field_error", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "malformed_json_error", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "unknown_reason_code_error", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "unsupported_hash_algorithm_error", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "wrong_profile_fail", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "inside_payload_mutation_detected", "expected_verify_result.json"),
		filepath.Join(repoRoot, "testdata", "core_2_conformance", "outside_metadata_ignored", "expected_verify_result.json"),
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

func compileVerifyResultV2Schema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	repoRoot := filepath.Join("..", "..")
	schemaPath := filepath.Join(repoRoot, "schemas", "verify_result_v2.schema.json")
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatalf("failed to resolve absolute path for schema %s: %v", schemaPath, err)
	}
	compiler := jsonschema.NewCompiler()
	schemaURL := "file:///" + filepath.ToSlash(absSchemaPath)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("failed to compile schema %s: %v", schemaURL, err)
	}
	return sch
}

func v2BaseDoc() map[string]interface{} {
	return map[string]interface{}{
		"verify_result_version":    "v2-draft",
		"result":                   "PASS",
		"reason_code":              "STATE_RECONSTRUCTED",
		"canonicalization_profile": "digiemu-canonical-json-v1",
		"hash_algorithm":           "sha256",
		"reconstruction_profile":   "digiemu-reconstruct-v1",
		"verifier_version":         "0.0.0-draft",
	}
}

func TestVerifyResultV2_PASS_RequiresSnapshotHash(t *testing.T) {
	sch := compileVerifyResultV2Schema(t)

	// Positive: PASS with snapshot_hash validates.
	pass := v2BaseDoc()
	pass["snapshot_hash"] = "sha256:" + "a"
	if err := sch.Validate(pass); err != nil {
		t.Fatalf("PASS with snapshot_hash must validate: %v", err)
	}

	// Negative: PASS without snapshot_hash must be rejected.
	missing := v2BaseDoc()
	if err := sch.Validate(missing); err == nil {
		t.Fatalf("PASS without snapshot_hash must be rejected")
	}
}

func TestVerifyResultV2_FAILAndERROR_DoNotRequireSnapshotHash(t *testing.T) {
	sch := compileVerifyResultV2Schema(t)

	fail := v2BaseDoc()
	fail["result"] = "FAIL"
	fail["reason_code"] = "HASH_MISMATCH"
	fail["snapshot_hash_expected"] = "sha256:e"
	fail["snapshot_hash_actual"] = "sha256:a"
	if err := sch.Validate(fail); err != nil {
		t.Fatalf("FAIL without snapshot_hash must remain valid: %v", err)
	}

	errDoc := v2BaseDoc()
	errDoc["result"] = "ERROR"
	errDoc["reason_code"] = "INVALID_SNAPSHOT_SCHEMA"
	if err := sch.Validate(errDoc); err != nil {
		t.Fatalf("ERROR without snapshot_hash must remain valid: %v", err)
	}
}
