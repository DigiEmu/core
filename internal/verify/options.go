package verify

// Options controls additional operator behaviors beyond pure verification.
type Options struct {
	// WriteExpected controls whether the operator may write expected_hash_v1 into snapshot.json
	// when a placeholder value is present.
	WriteExpected bool
}
