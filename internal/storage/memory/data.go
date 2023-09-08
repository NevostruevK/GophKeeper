package memory

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/google/uuid"
)

func (s *DataStore) AddData(_ context.Context, id uuid.UUID, d models.Data) {
	s.Data.Store(id, d)
}

func (s *DataStore) GetData(_ context.Context, ds *models.DataSpec) (models.Data, error) {
	v, ok := s.Data.Load(ds.ID)
	if !ok {
		return nil, storage.ErrNotFound
	}
	d, ok := v.(models.Data)
	if !ok {
		return nil, ErrTypeAssert
	}
	return d, nil
}

func (s *DataStore) AddEntry(_ context.Context, id uuid.UUID, d models.Entry) {
	s.Data.Store(id, d)
}

func (s *DataStore) GetEntry(_ context.Context, ds *models.DataSpec) (models.Entry, error) {
	v, ok := s.Data.Load(ds.ID)
	if !ok {
		return nil, storage.ErrNotFound
	}
	d, ok := v.(models.Entry)
	if !ok {
		return nil, ErrTypeAssert
	}
	return d, nil
}
