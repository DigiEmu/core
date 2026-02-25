package snapshot

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"digiemu-core/internal/canonicaljson"
)

type HashV1 string

// HashV1FromState computes SHA-256 over canonicaljson.Marshal(state) and returns uppercase hex string.
func HashV1FromState(state any) (HashV1, error) {
	b, err := canonicaljson.Marshal(state)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return HashV1(strings.ToUpper(fmt.Sprintf("%x", sum))), nil
}
