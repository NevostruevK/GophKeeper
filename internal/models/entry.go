package models

type Entry interface {
	String() string
	Show() string
	IsReadyForStorage() (bool, string)
}

func NewEntry(typ MType) Entry {
	switch typ {
	case CARD:
		return &Card{}
	case FILE:
		return &File{}
	case PAIR:
		return &Pair{}
	default:
		return &Text{}
	}
}
