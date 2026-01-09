package errs

import "errors"

var (
	ErrTitleRequired           = errors.New("title is required")                    // ErrTitleRequired is returned when the title is required.
	ErrTitleTooLong            = errors.New("title must be at most 200 characters") // ErrTitleTooLong is returned when the title is too long.
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailAlreadyExist       = errors.New("email already exists")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)
