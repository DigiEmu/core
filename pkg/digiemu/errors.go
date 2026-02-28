package digiemu

import (
	"errors"
	"fmt"
)

// Code is a stable, machine-readable error code for reporting and CLI mapping.
type Code string

const (
	CodeBundleInvalid Code = "BUNDLE_INVALID"
	CodeVerifyFailed  Code = "VERIFY_FAILED"
	CodeSchemaInvalid Code = "SCHEMA_INVALID"
	CodeHashMismatch  Code = "HASH_MISMATCH"
	CodeFileMissing   Code = "FILE_MISSING"
	CodeIO            Code = "IO"
	CodeInternal      Code = "INTERNAL"
	CodeUsage         Code = "USAGE" // primarily CLI-layer
)

// Sentinel root errors (stable). Use errors.Is(err, ErrX) for classification.
var (
	ErrBundleInvalid = errors.New(string(CodeBundleInvalid))
	ErrVerifyFailed  = errors.New(string(CodeVerifyFailed))
	ErrSchemaInvalid = errors.New(string(CodeSchemaInvalid))
	ErrHashMismatch  = errors.New(string(CodeHashMismatch))
	ErrFileMissing   = errors.New(string(CodeFileMissing))
	ErrIO            = errors.New(string(CodeIO))
	ErrInternal      = errors.New(string(CodeInternal))
	ErrUsage         = errors.New(string(CodeUsage))
)

// CodeOf attempts to extract a stable Code from err.
//
// - If err (or any wrapped error) is a *Detail, its Code is returned.
// - Otherwise, it falls back to the stable sentinel set via errors.Is.
func CodeOf(err error) (Code, bool) {
	if err == nil {
		return "", false
	}

	var d *Detail
	if errors.As(err, &d) {
		return d.Code, true
	}

	switch {
	case errors.Is(err, ErrBundleInvalid):
		return CodeBundleInvalid, true
	case errors.Is(err, ErrVerifyFailed):
		return CodeVerifyFailed, true
	case errors.Is(err, ErrSchemaInvalid):
		return CodeSchemaInvalid, true
	case errors.Is(err, ErrHashMismatch):
		return CodeHashMismatch, true
	case errors.Is(err, ErrFileMissing):
		return CodeFileMissing, true
	case errors.Is(err, ErrIO):
		return CodeIO, true
	case errors.Is(err, ErrInternal):
		return CodeInternal, true
	case errors.Is(err, ErrUsage):
		return CodeUsage, true
	default:
		return "", false
	}
}

// Detail carries deterministic context and wraps a sentinel Root.
// Root MUST be one of the sentinels above for stable classification.
type Detail struct {
	Code  Code              // stable machine code
	Op    string            // operation label: "verify.bundle", "read.file", ...
	Path  string            // relevant path (relative if possible)
	Field string            // schema/json field if applicable
	Want  string            // expected value (hash, version, etc.)
	Got   string            // observed value
	Meta  map[string]string // optional; keep minimal and deterministic
	Cause error             // underlying error (io/json/etc.)
	Root  error             // sentinel root for errors.Is classification
}

func (e *Detail) Error() string {
	// Human string. Keep short & deterministic.
	base := fmt.Sprintf("%s", e.Code)
	if e.Op != "" {
		base += " op=" + e.Op
	}
	if e.Path != "" {
		base += " path=" + e.Path
	}
	if e.Field != "" {
		base += " field=" + e.Field
	}
	if e.Want != "" || e.Got != "" {
		base += " want=" + e.Want + " got=" + e.Got
	}
	if e.Cause != nil {
		base += ": " + e.Cause.Error()
	}
	return base
}

func (e *Detail) Unwrap() error {
	// Ensure errors.Is(err, Root) works.
	if e.Root != nil {
		return e.Root
	}
	return e.Cause
}

// --- Constructors (keep stable semantics; OK to extend fields later) ---

func newDetail(code Code, root error, op, path, field, want, got string, meta map[string]string, cause error) error {
	return &Detail{
		Code:  code,
		Root:  root,
		Op:    op,
		Path:  path,
		Field: field,
		Want:  want,
		Got:   got,
		Meta:  meta,
		Cause: cause,
	}
}

func BundleInvalid(op, path string, cause error) error {
	return newDetail(CodeBundleInvalid, ErrBundleInvalid, op, path, "", "", "", nil, cause)
}

func FileMissing(op, path string, cause error) error {
	return newDetail(CodeFileMissing, ErrFileMissing, op, path, "", "", "", nil, cause)
}

func HashMismatch(op, path, want, got string, cause error) error {
	return newDetail(CodeHashMismatch, ErrHashMismatch, op, path, "", want, got, nil, cause)
}

func SchemaInvalid(op, path, field string, cause error) error {
	return newDetail(CodeSchemaInvalid, ErrSchemaInvalid, op, path, field, "", "", nil, cause)
}

func VerifyFailed(op, path string, cause error) error {
	return newDetail(CodeVerifyFailed, ErrVerifyFailed, op, path, "", "", "", nil, cause)
}

func IOError(op, path string, cause error) error {
	return newDetail(CodeIO, ErrIO, op, path, "", "", "", nil, cause)
}

func InternalError(op string, cause error) error {
	return newDetail(CodeInternal, ErrInternal, op, "", "", "", "", nil, cause)
}

func UsageError(op, path string, cause error) error {
	return newDetail(CodeUsage, ErrUsage, op, path, "", "", "", nil, cause)
}
