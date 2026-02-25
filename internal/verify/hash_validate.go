package verify

// IsValidSHA256HexLower reports whether s is exactly a 64-character lowercase hex string.
// Allowed characters: [0-9a-f].
func IsValidSHA256HexLower(s string) bool {
	if len(s) != 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			continue
		}
		return false
	}
	return true
}
