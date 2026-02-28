package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunVerify_JSONCanonical_StdoutOnlyJSON(t *testing.T) {
	tmp := t.TempDir()
	ref := "demo"
	root := filepath.Join(tmp, "snapshots", ref)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Minimal bundle: snapshot.json only. expected_hash_v1 can be any non-empty string.
	content := "{\n" +
		"  \"expected_hash_v1\": \"NOPE\",\n" +
		"  \"foo\": 1\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join(root, "snapshot.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write snapshot.json: %v", err)
	}

	args := []string{
		"--ref", ref,
		"--data", tmp,
		"--fixture-root", tmp,
		"--json=canonical",
	}

	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("expected exit=1 (verify_fail), got %d", exit)
	}

	if got := stderr.String(); got != "" {
		t.Fatalf("expected empty stderr for verify_fail in json mode, got: %q", got)
	}

	out := stdout.String()
	if strings.HasPrefix(out, "OK ") || strings.HasPrefix(out, "FAIL ") {
		t.Fatalf("expected no human output on stdout in json mode")
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("stdout not valid json: %v\nstdout=%q", err, out)
	}
	if ok, _ := decoded["ok"].(bool); ok {
		t.Fatalf("expected ok=false")
	}
}

func TestRunVerify_JSONPrettyFlag_StdoutOnlyJSON(t *testing.T) {
	tmp := t.TempDir()
	ref := "demo"
	root := filepath.Join(tmp, "snapshots", ref)
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	content := "{\n" +
		"  \"expected_hash_v1\": \"NOPE\",\n" +
		"  \"foo\": 1\n" +
		"}\n"
	if err := os.WriteFile(filepath.Join(root, "snapshot.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write snapshot.json: %v", err)
	}

	args := []string{
		"--ref", ref,
		"--data", tmp,
		"--fixture-root", tmp,
		"--json",
	}

	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO(args, &stdout, &stderr)
	if exit != 1 {
		t.Fatalf("expected exit=1 (verify_fail), got %d", exit)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("expected empty stderr for verify_fail in json mode, got: %q", got)
	}
	out := stdout.String()
	if !strings.Contains(out, "\n  ") {
		t.Fatalf("expected pretty JSON with indentation")
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("stdout not valid json: %v", err)
	}
}

func TestRunVerify_JSONMode_UsageErrorsGoToStderr(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exit := runVerifyWithIO([]string{"--json=canonical"}, &stdout, &stderr)
	if exit != 2 {
		t.Fatalf("expected exit=2 (usage), got %d", exit)
	}
	if stdout.String() != "" {
		t.Fatalf("expected empty stdout on usage error")
	}
	if stderr.String() == "" {
		t.Fatalf("expected stderr to contain usage output")
	}
}
