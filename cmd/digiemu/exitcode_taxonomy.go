package main

import (
	"errors"

	derr "digiemu-core/pkg/digiemu"
)

// exitCodeForTaxonomy maps the stable error taxonomy to the CLI exit-code space.
//
// Phase A note: this is a hook point only; it is intentionally not wired into
// command execution yet to avoid changing CLI UX/semantics.
func exitCodeForTaxonomy(err error) (code int, ok bool) {
	if err == nil {
		return 0, false
	}

	switch {
	case errors.Is(err, derr.ErrHashMismatch):
		return 2, true

	case errors.Is(err, derr.ErrBundleInvalid),
		errors.Is(err, derr.ErrSchemaInvalid),
		errors.Is(err, derr.ErrFileMissing),
		errors.Is(err, derr.ErrIO),
		errors.Is(err, derr.ErrVerifyFailed),
		errors.Is(err, derr.ErrUsage):
		return 4, true

	case errors.Is(err, derr.ErrInternal):
		return 5, true

	default:
		return 0, false
	}
}
