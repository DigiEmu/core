package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

func TestProfileRegistryValidatesSchema(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	schemaPath := filepath.Join(repoRoot, "schemas", "core_2_profile_registry.schema.json")
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		t.Fatalf("abs schema: %v", err)
	}
	compiler := jsonschema.NewCompiler()
	schemaURL := "file:///" + filepath.ToSlash(absSchemaPath)
	sch, err := compiler.Compile(schemaURL)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}

	regPath := filepath.Join(repoRoot, "testdata", "core_2_profiles", "profile_registry_valid.json")
	data, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("invalid registry json: %v", err)
	}
	if err := sch.Validate(v); err != nil {
		t.Fatalf("registry failed schema validation: %v", err)
	}
}

func TestProfileRegistry_VerifyResultVersions_ContainsCanonicalV2Draft(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	regPath := filepath.Join(repoRoot, "testdata", "core_2_profiles", "profile_registry_valid.json")
	data, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read registry: %v", err)
	}
	var reg map[string]interface{}
	if err := json.Unmarshal(data, &reg); err != nil {
		t.Fatalf("invalid registry json: %v", err)
	}
	vers, ok := reg["verify_result_versions"].(map[string]interface{})
	if !ok {
		t.Fatalf("verify_result_versions missing or wrong type")
	}
	active, ok := vers["active"].([]interface{})
	if !ok {
		t.Fatalf("verify_result_versions.active missing or wrong type")
	}
	found := false
	for _, a := range active {
		if a == "v2-draft" {
			found = true
		}
	}
	if !found {
		t.Fatalf("verify_result_versions.active does not contain canonical v2-draft: %v", active)
	}
}
