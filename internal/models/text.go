package models

import (
	"bytes"
	"encoding/gob"
)

// Text структра для хранения произвольных текстовых данных.
type Text struct {
	Text []byte
}

func NewText(text []byte) *Text {
	return &Text{text}
}

func (t *Text) Encode() ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(t)
	return buff.Bytes(), err
}

func (t *Text) Decode(d []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(d))
	err := dec.Decode(t)
	return err
}

func (t *Text) Type() MType {
	return TEXT
}
