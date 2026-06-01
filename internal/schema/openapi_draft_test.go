package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestOpenAPIDraftBasicStructure(t *testing.T) {
	repoRoot := filepath.Join("..", "..")
	openapiPath := filepath.Join(repoRoot, "openapi", "core_2_conformance_api.yaml")
	data, err := os.ReadFile(openapiPath)
	if err != nil {
		t.Fatalf("failed to read openapi file %s: %v", openapiPath, err)
	}

	var v map[string]interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	// basic required fields
	if _, ok := v["openapi"]; !ok {
		t.Fatalf("openapi version field missing")
	}
	info, ok := v["info"].(map[string]interface{})
	if !ok {
		t.Fatalf("info object missing or wrong type")
	}
	if _, ok := info["title"]; !ok {
		t.Fatalf("info.title missing")
	}

	// paths checks
	paths, ok := v["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("paths object missing or wrong type")
	}

	required := []string{
		"/v1/core2/version",
		"/v1/core2/profiles",
		"/v1/core2/conformance/run",
		"/v1/core2/verify-result/validate",
		"/v1/core2/conformance-report/validate",
	}
	for _, p := range required {
		if _, ok := paths[p]; !ok {
			t.Fatalf("required path %s missing", p)
		}
	}

	// components.schemas exists
	components, ok := v["components"].(map[string]interface{})
	if !ok {
		t.Fatalf("components object missing or wrong type")
	}
	if _, ok := components["schemas"]; !ok {
		t.Fatalf("components.schemas missing")
	}
}
