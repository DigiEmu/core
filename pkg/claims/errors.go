package claims

import "errors"

var (
	ErrInvalidID      = errors.New("invalid claim id")
	ErrInvalidSubject = errors.New("invalid claim subject")
)
