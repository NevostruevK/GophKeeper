package storage

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

type Storage interface {
	AddUser(ctx context.Context, user models.User) (id uuid.UUID, err error)
	GetUser(ctx context.Context, user models.User) (uuid.UUID, error)
	GetSpecs(ctx context.Context, userID uuid.UUID) ([]models.Spec, error)
	GetSpecsOfType(ctx context.Context, userID uuid.UUID, mType models.MType) ([]models.Spec, error)
	AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (*models.DataSpec, error)
	GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error)
}
