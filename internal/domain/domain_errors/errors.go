package domain_errors

import "errors"

var (
	ErrFlightInPast            = errors.New("cannot create flight in the past")
	ErrInvalidStatusTransition = errors.New("invalid flight status transition")
)
