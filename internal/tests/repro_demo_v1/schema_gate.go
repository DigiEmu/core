package repro_demo_v1

import (
	"encoding/json"
	"fmt"
)

// Offline schema gate:
// We validate a limited but strong subset of JSON Schema:
// - schema must declare required keys
// - properties must define basic types for those keys
//
// This avoids adding external dependencies while still catching drift.

type jsonSchema struct {
	Type       any                   `json:"type"`
	Required   []string              `json:"required"`
	Properties map[string]schemaProp `json:"properties"`
}

type schemaProp struct {
	Type any `json:"type"`
}

func schemaTypeToSet(v any) map[string]bool {
	set := map[string]bool{}
	switch t := v.(type) {
	case string:
		set[t] = true
	case []any:
		for _, it := range t {
			if s, ok := it.(string); ok {
				set[s] = true
			}
		}
	}
	return set
}

func ValidateAgainstSchemaSubset(schemaBytes, docBytes []byte) error {
	var s jsonSchema
	if err := json.Unmarshal(schemaBytes, &s); err != nil {
		return fmt.Errorf("schema json parse: %w", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(docBytes, &doc); err != nil {
		return fmt.Errorf("doc json parse: %w", err)
	}

	// required fields present
	for _, k := range s.Required {
		if _, ok := doc[k]; !ok {
			return fmt.Errorf("missing required field: %s", k)
		}
	}

	// type checks for required fields (basic)
	for _, k := range s.Required {
		prop, ok := s.Properties[k]
		if !ok {
			// schema might omit; skip (still enforce presence)
			continue
		}
		want := schemaTypeToSet(prop.Type)
		if len(want) == 0 {
			continue
		}
		v := doc[k]
		// map json types to schema types
		actual := ""
		switch v.(type) {
		case string:
			actual = "string"
		case bool:
			actual = "boolean"
		case float64:
			actual = "number"
		case []any:
			actual = "array"
		case map[string]any:
			actual = "object"
		case nil:
			actual = "null"
		default:
			actual = "unknown"
		}
		if !want[actual] {
			return fmt.Errorf("field %s has type %s, expected one of %v", k, actual, keys(want))
		}
	}
	return nil
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
