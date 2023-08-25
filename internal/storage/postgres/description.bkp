package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
)

func (s *Storage) AddDescription(ctx context.Context, d *models.Description, id uuid.UUID) error {
	_, err := s.Exec(ctx, sql.InsertDescription, id, d.Description, d.IsCompressed)
	return err
}

func (s *Storage) GetDescription(ctx context.Context, id uuid.UUID) (*models.Description, error) {
	d := models.Description{}
	err := s.QueryRow(ctx, sql.SelectDescription, id).Scan(&d.Description, &d.IsCompressed)
	return &d, err
}
