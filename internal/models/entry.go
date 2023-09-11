package models

// Entry interface for stored data.
type Entry interface {
	String() string
	Show() string
	IsReadyForStorage() (bool, string)
}

// NewEntry returns new Entry for differnt type.
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
