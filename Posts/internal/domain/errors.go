package domain

import "errors"

// Errors that can be returned by repositories.
var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrCommentIsTooLong = errors.New("comment is too long")
)
