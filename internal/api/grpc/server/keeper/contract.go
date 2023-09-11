package keeper

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

// DataStore interface for gRPC KeeperServer.
type DataStore interface {
	GetSpecs(ctx context.Context, userID uuid.UUID) ([]models.Spec, error)
	GetSpecsOfType(ctx context.Context, userID uuid.UUID, mType models.MType) ([]models.Spec, error)
	AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (*models.DataSpec, error)
	GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error)
}
