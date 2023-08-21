package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) AddPair(ctx context.Context, p *models.Pair, m *models.Meta, d *models.Description) (uuid.UUID, error) {
	return s.AddData(ctx, p, m, d)
}

func (s *Storage) GetPair(ctx context.Context, m *models.Meta) (*models.Pair, error) {
	d, err := s.GetData(ctx, m)
	if err != nil {
		return nil, err
	}
	p := &models.Pair{}
	err = p.Decode(d)
	return p, err
}
