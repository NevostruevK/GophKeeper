package storage

import "errors"

var (
	ErrDuplicateLogin = errors.New("duplicate login")
	ErrNotFound       = errors.New("not found")
	ErrWrongPassword  = errors.New("wrong password")
)
