package models

type Description struct {
	Description  []byte
	IsCompressed bool
}

func NewDescription(description []byte) *Description {
	return &Description{Description: description}
}
