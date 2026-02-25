package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type replayGolden struct {
	Snapshot any      `json:"snapshot"`
	Claims   any      `json:"claims"`
	Trace    []string `json:"trace"`
}

func normalizeReplayGolden(t *testing.T, repoRoot string, raw []byte) []byte {
	t.Helper()

	var g replayGolden
	if err := json.Unmarshal(raw, &g); err != nil {
		t.Fatalf("invalid replay json: %v\nraw=%s", err, string(raw))
	}

	toStableTrace := func(s string) string {
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

	for i := range g.Trace {
		g.Trace[i] = toStableTrace(g.Trace[i])
	}

	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		t.Fatalf("marshal normalized replay json: %v", err)
	}
	return append(b, '\n')
}

func TestReplay_JSON_Golden(t *testing.T) {
	bin := buildBinaryOnce(t)

	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatal(err)
	}

	bundle := filepath.Join(repoRoot, "data", "test-fixtures", "snapshots", "demo")

	code, stdout, stderr := runCLI(t, bin, "replay", "--bundle", bundle, "--json")
	if code != 0 {
		t.Fatalf("replay failed: exit=%d\nstdout=%s\nstderr=%s", code, stdout, stderr)
	}

	got := normalizeReplayGolden(t, repoRoot, []byte(stdout))

	goldenPath := filepath.Join(repoRoot, "cmd", "digiemu", "testdata", "replay_demo_golden.json")

	if _, err := os.Stat(goldenPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Fatalf("golden file created — verify and re-run test")
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	want := normalizeReplayGolden(t, repoRoot, golden)

	if !bytes.Equal(want, got) {
		t.Fatalf("replay output changed — golden mismatch\n--- want ---\n%s\n--- got ---\n%s", string(want), string(got))
	}
}
