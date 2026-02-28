package tests

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyResultSchemaLocked(t *testing.T) {
	root := repoRoot(t) // defined in cli_version_test.go

	schemaPath := filepath.Join(root, "schemas", "VERIFY_RESULT_SCHEMA_v1.json")
	lockPath := filepath.Join(root, "schemas", "VERIFY_RESULT_SCHEMA_v1.sha256")

	lockBytes, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("read lock file: %v", err)
	}
	lock := strings.TrimSpace(string(lockBytes))
	if lock == "" {
		t.Fatalf("empty lock file: %s", lockPath)
	}

	hash, err := sha256FileHexLower(schemaPath)
	if err != nil {
		t.Fatalf("hash schema: %v", err)
	}
	if hash != lock {
		t.Fatalf("schema lock mismatch: got %s want %s", hash, lock)
	}
}

func sha256FileHexLower(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
