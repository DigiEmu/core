package main

import (
	"errors"

	derr "digiemu-core/pkg/digiemu"
)

// Phase A2 contract (locked):
// - ok: 0
// - verify_fail: 1
// - usage: 2
// - io: 3
// - internal: 4
func exitCodeForError(err error) int {
	if err == nil {
		return 0
	}

	switch {
	case errors.Is(err, derr.ErrUsage):
		return 2
	case errors.Is(err, derr.ErrIO):
		return 3
	case errors.Is(err, derr.ErrInternal):
		return 4
	case errors.Is(err, derr.ErrVerifyFailed),
		errors.Is(err, derr.ErrHashMismatch),
		errors.Is(err, derr.ErrSchemaInvalid),
		errors.Is(err, derr.ErrBundleInvalid),
		errors.Is(err, derr.ErrFileMissing):
		return 1
	default:
		// Unknown errors are treated as internal.
		return 4
	}
}

func isUsageError(err error) bool {
	return errors.Is(err, derr.ErrUsage)
}
