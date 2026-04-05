package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// NewID builds a deterministic ID from a stable prefix + input.
// IMPORTANT: input must itself be deterministic and stable.
func NewID(prefix string, input string) string {
	sum := sha256.Sum256([]byte(prefix + "|" + input))
	return prefix + "_" + hex.EncodeToString(sum[:16])
}

// NewIDParts is a helper for composing stable IDs from multiple deterministic parts.
func NewIDParts(prefix string, parts ...string) string {
	return NewID(prefix, strings.Join(parts, "|"))
}
