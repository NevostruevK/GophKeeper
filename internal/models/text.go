package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func (t *Text) String() string {
	return fmt.Sprintf("text size %d", len(t.Text))
}

func (t *Text) Show() string {
	return string(t.Text)
}

func (t *Text) IsReadyForStorage() (bool, string) {
	const textIsEmpty = "text is empty"
	if len(t.Text) > 0 {
		return true, ""
	}
	return false, textIsEmpty
}
