package auth

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

type UserStore interface {
	AddUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GetUser(ctx context.Context, user models.User) (uuid.UUID, error)
}
