package models

import (
	"fmt"
)

// Pair структра для хранения пар login/password.
type Pair struct {
	Login    string
	Password string
}

// NewPair returns Pair.
func NewPair(login, password string) *Pair {
	return &Pair{login, password}
}

// String prints Pair.
func (p *Pair) String() string {
	return fmt.Sprintf("Login: %s Password: %s", p.Login, p.Password)
}

// Show shows Pair information.
func (p *Pair) Show() string {
	return p.String()
}

// IsReadyForStorage check Pair's fields for ready to store.
func (p *Pair) IsReadyForStorage() (bool, string) {
	const (
		loginIsEmpty    = "login is empty"
		passwordIsShort = "password is short (at least 4 characters)"
		passwordMinLen  = 4
	)
	if p.Login == "" {
		return false, loginIsEmpty
	}
	if len(p.Password) < passwordMinLen {
		return false, passwordIsShort
	}
	return true, ""
}
