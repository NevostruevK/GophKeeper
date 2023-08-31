package models

import (
	"bytes"
	"encoding/gob"
)

// Pair структра для хранения пар login/password.
type Pair struct {
	Login    []byte
	Password []byte
}

func NewPair(login, password []byte) *Pair {
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
