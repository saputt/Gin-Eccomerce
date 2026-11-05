package auth

import "errors"

var (
	ErrInvalidCredential = errors.New("Invalid email or password")
	ErrInvalidToken      = errors.New("invalid or expired token")
)
