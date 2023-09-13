package service

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/api/grpc/client"
	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/google/uuid"
)

// Service service for data sinchronisation.
type Service struct {
	memory *memory.DataStore
	client *client.Client
}

// NewService returns new Service.
func NewService(client *client.Client) *Service {
	return &Service{
		memory.NewDataStore(),
		client,
	}
}

// LoadSpecs provides specs of data for user.
func (s Service) LoadSpecs(ctx context.Context, typ models.MType) ([]models.Spec, error) {
	var (
		specs []models.Spec
		err   error
	)
	if typ == models.NOTIMPLEMENT {
		specs, err = s.memory.GetSpecs(ctx, uuid.Nil)
	} else {
		specs, err = s.memory.GetSpecsOfType(ctx, uuid.Nil, typ)
	}
	if err != nil {
		return nil, err
	}
	if len(specs) > 0 {
		return specs, nil
	}
	if typ == models.NOTIMPLEMENT {
		specs, err = s.client.Keeper.GetSpecs(ctx)
	} else {
		specs, err = s.client.Keeper.GetSpecsOfType(ctx, typ)
	}
	if err != nil {
		return nil, err
	}
	s.memory.AddSpecs(specs)
	return specs, err
}
