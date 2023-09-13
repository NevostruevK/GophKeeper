// package memory storage in memory.
package memory

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/google/uuid"
)

// GetData load data.
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

// AddEntry store entry.
func (s *DataStore) AddEntry(_ context.Context, id uuid.UUID, d models.Entry) {
	s.Data.Store(id, d)
}

// GetEntry load entry.
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
