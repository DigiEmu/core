package verify

import "regexp"

const sha256HexPattern = `^[A-Fa-f0-9]{64}$`

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
	ok, err := regexp.MatchString(sha256HexPattern, s)
	return err == nil && ok
}
