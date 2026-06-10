package domain

import "errors"

// Sentinel errors the layers return and the HTTP boundary maps to status codes.
// Lower layers wrap these with context (fmt.Errorf("...: %w", err)); handlers use
// errors.Is to translate them (Architecture Guidelines §3.2).
var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrInvalidCredentials is returned on a failed login. Use it for BOTH a wrong
	// password and an unknown email so the API never reveals which emails exist.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrEmailTaken is returned when registering an email that already exists.
	ErrEmailTaken = errors.New("email already registered")
)
