package models

import (
	"fmt"
)

// Text структра для хранения произвольных текстовых данных.
type Text struct {
	Text []byte
}

// NewText returns Text.
func NewText(text []byte) *Text {
	return &Text{text}
}

// String prints Text.
func (t *Text) String() string {
	return fmt.Sprintf("text size %d", len(t.Text))
}

// Show shows Text information.
func (t *Text) Show() string {
	return string(t.Text)
}

// IsReadyForStorage check Text for ready to store.
func (t *Text) IsReadyForStorage() (bool, string) {
	const textIsEmpty = "text is empty"
	if len(t.Text) > 0 {
		return true, ""
	}
	return false, textIsEmpty
}
