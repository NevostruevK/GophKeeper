package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) AddCard(ctx context.Context, c *models.Card, m *models.Meta, d *models.Description) (uuid.UUID, error) {
	return s.AddData(ctx, c, m, d)
}

func (s *Storage) GetCard(ctx context.Context, m *models.Meta) (*models.Card, error) {
	d, err := s.GetData(ctx, m)
	if err != nil {
		return nil, err
	}
	c := &models.Card{}
	err = c.Decode(d)
	return c, err
}
