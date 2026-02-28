package digiemu

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorClassification_Is(t *testing.T) {
	t.Run("hash mismatch", func(t *testing.T) {
		err := HashMismatch("verify.hash", "hashes/manifest.sha256", "aaa", "bbb", nil)
		if !errors.Is(err, ErrHashMismatch) {
			t.Fatalf("expected errors.Is HashMismatch")
		}
		wrapped := fmt.Errorf("wrapper: %w", err)
		if !errors.Is(wrapped, ErrHashMismatch) {
			t.Fatalf("expected errors.Is HashMismatch through wrapper")
		}
	})

	t.Run("io error", func(t *testing.T) {
		cause := fmt.Errorf("permission denied")
		err := IOError("read.file", "bundle/VERIFY_RESULT.json", cause)
		if !errors.Is(err, ErrIO) {
			t.Fatalf("expected errors.Is IO")
		}
	})

	t.Run("bundle invalid", func(t *testing.T) {
		err := BundleInvalid("verify.bundle", "bundle/", nil)
		if !errors.Is(err, ErrBundleInvalid) {
			t.Fatalf("expected errors.Is BundleInvalid")
		}
	})

	t.Run("schema invalid", func(t *testing.T) {
		err := SchemaInvalid("verify.schema", "schemas/VERIFY_RESULT_SCHEMA_v1.json", "report.failures[0]", nil)
		if !errors.Is(err, ErrSchemaInvalid) {
			t.Fatalf("expected errors.Is SchemaInvalid")
		}
	})

	t.Run("usage", func(t *testing.T) {
		err := UsageError("cli.verify", "", fmt.Errorf("bad flag"))
		if !errors.Is(err, ErrUsage) {
			t.Fatalf("expected errors.Is Usage")
		}
	})
}

func TestCodeOf(t *testing.T) {
	t.Run("detail", func(t *testing.T) {
		err := HashMismatch("verify.hash", "hashes/manifest.sha256", "aaa", "bbb", nil)
		code, ok := CodeOf(err)
		if !ok {
			t.Fatalf("expected ok")
		}
		if code != CodeHashMismatch {
			t.Fatalf("expected %q, got %q", CodeHashMismatch, code)
		}
	})

	t.Run("wrapped detail", func(t *testing.T) {
		inner := SchemaInvalid("verify.schema", "schemas/VERIFY_RESULT_SCHEMA_v1.json", "report.failures[0]", nil)
		err := fmt.Errorf("wrap: %w", inner)
		code, ok := CodeOf(err)
		if !ok {
			t.Fatalf("expected ok")
		}
		if code != CodeSchemaInvalid {
			t.Fatalf("expected %q, got %q", CodeSchemaInvalid, code)
		}
	})

	t.Run("sentinel", func(t *testing.T) {
		code, ok := CodeOf(ErrIO)
		if !ok {
			t.Fatalf("expected ok")
		}
		if code != CodeIO {
			t.Fatalf("expected %q, got %q", CodeIO, code)
		}
	})
}
