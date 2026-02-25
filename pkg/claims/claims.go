package claims

// ClaimID is a stable identifier for a claim.
type ClaimID string

// Claim is a minimal public claim representation.
type Claim struct {
	ID      ClaimID `json:"id"`
	Subject string  `json:"subject"`
	Payload any     `json:"payload"` // v1.0: payload is opaque to the kernel
}

// Validate performs minimal invariants for v1.0.
func (c Claim) Validate() error {
	if c.ID == "" {
		return ErrInvalidID
	}
	if c.Subject == "" {
		return ErrInvalidSubject
	}
	return nil
}
