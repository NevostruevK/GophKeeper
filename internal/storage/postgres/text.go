package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) AddText(ctx context.Context, t *models.Text, m *models.Meta, d *models.Description) (uuid.UUID, error) {
	return s.AddData(ctx, t, m, d)
}

func (s *Storage) GetText(ctx context.Context, m *models.Meta) (*models.Text, error) {
	d, err := s.GetData(ctx, m)
	if err != nil {
		return nil, err
	}
	t := &models.Text{}
	err = t.Decode(d)
	return t, err
}
