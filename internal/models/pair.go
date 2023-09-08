package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Pair структра для хранения пар login/password.
type Pair struct {
	Login    string
	Password string
}

func NewPair(login, password string) *Pair {
	return &Pair{login, password}
}

func (p *Pair) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(p)
	return err
}

func (p *Pair) Type() MType {
	return PAIR
}

func (p *Pair) String() string {
	return fmt.Sprintf("Login: %s Password: %s", p.Login, p.Password)
}

func (p *Pair) Show() string {
	return p.String()
}

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
