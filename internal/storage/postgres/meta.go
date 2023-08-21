package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddMeta(ctx context.Context, m *models.Meta) (id uuid.UUID, err error) {
	return id, s.QueryRow(ctx, sql.InsertMeta, m.ID, m.UserID, m.MType, m.Title, m.DataSize, m.DescriptionSize).Scan(&id)
}

func (s *Storage) GetMetas(ctx context.Context, userID uuid.UUID) ([]models.Meta, error) {
	var count int
	if err := s.QueryRow(ctx, "SELECT COUNT(*) FROM metas WHERE user_id = $1", userID).Scan(&count); err != nil {
		return nil, err
	}
	rows, err := s.Query(ctx, sql.SelectMetasForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scan(ctx, rows, count)
}

func (s *Storage) GetMetasForType(ctx context.Context, userID uuid.UUID, mType models.MType) ([]models.Meta, error) {
	var count int
	if err := s.QueryRow(ctx, "SELECT COUNT(*) FROM metas WHERE user_id = $1 AND type = $2", userID, mType).Scan(&count); err != nil {
		return nil, err
	}
	rows, err := s.Query(ctx, sql.SelectMetasForUserAndType, userID, mType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scan(ctx, rows, count)
}

func scan(_ context.Context, rows pgx.Rows, count int) ([]models.Meta, error) {
	metas := make([]models.Meta, 0, count)
	for rows.Next() {
		meta := models.Meta{}
		if err := rows.Scan(&meta.ID, &meta.MType, &meta.Title, &meta.DataSize, &meta.DescriptionSize); err != nil {
			return metas, err
		}
		metas = append(metas, meta)
	}
	return metas, rows.Err()
}
