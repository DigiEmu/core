package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isAbsSlashPath(p string) bool {
	if p == "" {
		return false
	}
	// Unix-style absolute.
	if strings.HasPrefix(p, "/") {
		return true
	}
	// Windows drive absolute, normalized to slash form (e.g. C:/foo).
	if len(p) >= 3 && isAlpha(p[0]) && p[1] == ':' && p[2] == '/' {
		return true
	}
	return false
}

func stableTraceEntry(t *testing.T, repoRoot, s string) string {
	t.Helper()

	usedPrefix := ""
	p := s
	if strings.HasPrefix(s, "used:") {
		usedPrefix = "used:"
		p = strings.TrimPrefix(s, "used:")
	}

	// Make trace stable across OSes:
	// - Convert any backslashes to forward slashes (even on Unix).
	// - If the path is absolute, make it repo-relative when possible.
	pSlash := strings.ReplaceAll(p, "\\", "/")
	repoSlash := strings.ReplaceAll(repoRoot, "\\", "/")

	// Make absolute paths stable by converting them to repo-relative.
	if isAbsSlashPath(pSlash) {
		pOS := filepath.FromSlash(pSlash)
		repoOS := filepath.FromSlash(repoSlash)
		if rel, err := filepath.Rel(repoOS, pOS); err == nil {
			relSlash := filepath.ToSlash(rel)
			if relSlash != "." && !strings.HasPrefix(relSlash, "..") && !strings.Contains(relSlash, ":") {
				pSlash = relSlash
			}
		}
	} else if filepath.IsAbs(p) {
		// Fallback for OS-native absolute paths.
		if rel, err := filepath.Rel(repoRoot, p); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
			pSlash = filepath.ToSlash(rel)
		}
	}

	pSlash = strings.TrimPrefix(pSlash, "./")
	return usedPrefix + pSlash
}

func stableTrace(t *testing.T, repoRoot string, trace []string) []string {
	t.Helper()
	out := make([]string, 0, len(trace))
	for _, s := range trace {
		out = append(out, stableTraceEntry(t, repoRoot, s))
	}
	return out
}
