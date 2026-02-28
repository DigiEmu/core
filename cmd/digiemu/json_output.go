package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"digiemu-core/internal/canonicaljson"
)

type jsonMode string

const (
	jsonModeNone      jsonMode = ""
	jsonModePretty    jsonMode = "pretty"
	jsonModeCanonical jsonMode = "canonical"
)

func normalizeJSONFlagArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--json" {
			out = append(out, "--json=pretty")
			continue
		}
		if a == "--json=true" {
			out = append(out, "--json=pretty")
			continue
		}
		if a == "--json=false" {
			continue
		}
		out = append(out, a)
	}
	return out
}

func parseJSONMode(s string) (jsonMode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return jsonModeNone, nil
	}
	if s == "true" {
		return jsonModePretty, nil
	}
	switch jsonMode(s) {
	case jsonModePretty, jsonModeCanonical:
		return jsonMode(s), nil
	default:
		return jsonModeNone, fmt.Errorf("invalid --json mode %q (expected %q or %q)", s, jsonModePretty, jsonModeCanonical)
	}
}

func writePrettyJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func writeCanonicalJSON(w io.Writer, v any) error {
	// Normalize through encoding/json (omitempty, embedded flattening, etc.),
	// then canonicalize the generic value so map key ordering is stable.
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var anyV any
	if err := json.Unmarshal(b, &anyV); err != nil {
		return err
	}
	cb, err := canonicaljson.Marshal(anyV)
	if err != nil {
		return err
	}
	if _, err := w.Write(cb); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}
