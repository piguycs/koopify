package user

import "errors"

const PG_UNIQUE_VIOLATION = "23505"

var (
	ErrUserExists               = errors.New("user already exists")
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrUserNotFound             = errors.New("user not found")
	ErrAccountDeletionScheduled = errors.New("account scheduled for deletion")
	ErrDeletionPolicyNotFound   = errors.New("deletion policy not found")
	ErrCannotEditAdmin          = errors.New("cannot edit admin users")
)
