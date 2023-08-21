package models

import (
	"bytes"
	"encoding/gob"
)

// Card структра для хранения банковских карт.
type Card struct {
	Owner           []byte
	ExpirationMonth uint8
	ExpirationYear  uint16
	Bank            []byte
	Number          uint64
	CVV             uint16
}

func NewCard(owner, bank []byte, number uint64, expMonth uint8, expYear, cvv uint16) *Card {
	return &Card{
		Owner:           owner,
		Bank:            bank,
		Number:          number,
		ExpirationYear:  expYear,
		ExpirationMonth: expMonth,
		CVV:             cvv,
	}
}

func (c *Card) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(c)
	return err
}

func (c *Card) Type() MType {
	return CARD
}
