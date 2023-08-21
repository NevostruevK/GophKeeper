package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) AddFile(ctx context.Context, f *models.File, m *models.Meta, d *models.Description) (uuid.UUID, error) {
	return s.AddData(ctx, f, m, d)
}

func (s *Storage) GetFile(ctx context.Context, m *models.Meta) (*models.File, error) {
	d, err := s.GetData(ctx, m)
	if err != nil {
		return nil, err
	}
	f := &models.File{}
	err = f.Decode(d)
	return f, err
}
