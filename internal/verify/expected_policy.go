package verify

import "regexp"

var sha256Hex = regexp.MustCompile(`^[A-Fa-f0-9]{64}$`)

// IsPlaceholderExpected returns true for known placeholder values or empty string.
func IsPlaceholderExpected(s string) bool {
	if s == "" {
		return true
	}
	switch s {
	case "REPLACE_AFTER_ASSEMBLE", "REPLACE_AFTER_COMPUTE":
		return true
	default:
		return false
	}
}

// IsRealHash reports whether s looks like a 64-hex sha256 digest.
func IsRealHash(s string) bool {
	return sha256Hex.MatchString(s)
}
