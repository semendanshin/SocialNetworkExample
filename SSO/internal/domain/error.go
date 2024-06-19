package domain

import "errors"

// Errors that can be returned by repositories.
var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidPassword = errors.New("invalid password")
)
