package models

import (
	"github.com/NevostruevK/GophKeeper/internal/utils/cut"
	"github.com/google/uuid"
)

const (
	TitleSize = 256
)

type MType string

const (
	PAIR MType = "PAIR"
	TEXT MType = "TEXT"
	FILE MType = "FILE"
	CARD MType = "CARD"
)

const ()

type Meta struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MType           MType
	Title           string
	DescriptionSize int
	DataSize        int
}

func NewMeta(userID uuid.UUID, mType MType, title string) *Meta {
	return &Meta{
		ID:     uuid.New(),
		UserID: userID,
		MType:  mType,
		Title:  cut.Cut(title, TitleSize),
	}
}

func (m *Meta) SetDescriptionSize(dsc *Description) {
	if dsc != nil {
		m.DescriptionSize = len(dsc.Description)
	}
}

func (m *Meta) SetDataSize(size int) {
	m.DataSize = size
}
