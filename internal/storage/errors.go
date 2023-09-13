package storage

import "errors"

var (
	// ErrDuplicateLogin error: duplicate login
	ErrDuplicateLogin = errors.New("duplicate login")
	// ErrNotFound error: not found
	ErrNotFound = errors.New("not found")
	// ErrWrongPassword error: wrong password
	ErrWrongPassword = errors.New("wrong password")
)
