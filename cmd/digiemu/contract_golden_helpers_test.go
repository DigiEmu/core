package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func stableTraceEntry(t *testing.T, repoRoot, s string) string {
	t.Helper()

	usedPrefix := ""
	p := s
	if strings.HasPrefix(s, "used:") {
		usedPrefix = "used:"
		p = strings.TrimPrefix(s, "used:")
	}

	// Make absolute paths stable by converting them to repo-relative.
	if filepath.IsAbs(p) {
		if rel, err := filepath.Rel(repoRoot, p); err == nil && !strings.HasPrefix(rel, "..") {
			p = rel
		}
	}

	p = filepath.ToSlash(p)
	return usedPrefix + p
}

func stableTrace(t *testing.T, repoRoot string, trace []string) []string {
	t.Helper()
	out := make([]string, 0, len(trace))
	for _, s := range trace {
		out = append(out, stableTraceEntry(t, repoRoot, s))
	}
	return out
}
