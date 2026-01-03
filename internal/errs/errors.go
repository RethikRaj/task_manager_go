package errs

import "errors"

var (
	ErrTitleRequired = errors.New("title is required")                    // ErrTitleRequired is returned when the title is required.
	ErrTitleTooLong  = errors.New("title must be at most 200 characters") // ErrTitleTooLong is returned when the title is too long.
)
