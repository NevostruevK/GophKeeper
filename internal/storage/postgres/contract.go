package postgres

import "github.com/NevostruevK/GophKeeper/internal/models"

type Data interface {
	Decode(data []byte) error
	Type() models.MType
}
