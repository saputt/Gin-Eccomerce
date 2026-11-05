package user

import "errors"

var (
	ErrEmailAlreadyExist = errors.New("Email already exist")
	ErrUserNotFound      = errors.New("User not found")
)
