package keeper

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

type DataStore interface {
	GetSpecs(ctx context.Context, userID uuid.UUID) ([]models.Spec, error)
	//	GetRecordsForType(ctx context.Context, userID uuid.UUID, mType models.MType) ([]models.Record, error)
	AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (uuid.UUID, error)
	GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error)
}
