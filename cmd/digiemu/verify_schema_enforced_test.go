package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestVerify_JSON_ValidatesAgainstSchemaV1(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	bin := buildDigiemuBinaryTempDir(t, repoRoot)

	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")
	code, stdout, stderr := runCLI(t, bin, "verify", "--bundle", bundle, "--json")
	if code != 0 {
		t.Fatalf("verify failed: exit=%d\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}

	var instance any
	if err := json.Unmarshal([]byte(stdout), &instance); err != nil {
		t.Fatalf("invalid JSON output: %v\nstdout=%s", err, stdout)
	}

	schemaPath := filepath.Join(repoRoot, "docs", "VERIFY_RESULT_SCHEMA_v1.json")
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	b = stripUTF8BOM(b)

	// Compile using the schema's $id to match internal resolution.
	schemaID := "https://digiemu.dev/schemas/VERIFY_RESULT_SCHEMA_v1.json"
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaID, strings.NewReader(string(b))); err != nil {
		t.Fatalf("add schema resource: %v", err)
	}

	schema, err := compiler.Compile(schemaID)
	if err != nil {
		t.Fatalf("compile schema: %v", err)
	}

	if err := schema.Validate(instance); err != nil {
		t.Fatalf("schema validation failed: %v\nstdout=%s", err, stdout)
	}
}
