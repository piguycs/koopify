package user

import "errors"

const PG_UNIQUE_VIOLATION = "23505"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
