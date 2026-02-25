package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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

	g.Trace = stableTrace(t, repoRoot, g.Trace)

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
