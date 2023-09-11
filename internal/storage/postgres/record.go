package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) crypt(d models.Data) []byte {
	return s.crypto.Crypt(d)
}

// AddRecord Add record in database.
func (s *Storage) AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (*models.DataSpec, error) {
	d := s.crypt(r.Data)
	ds := &models.DataSpec{DataSize: len(d)}
	err := s.QueryRow(ctx, sql.InsertRecord, userID, r.Type, r.Title, d, ds.DataSize).Scan(&ds.ID)
	return ds, err
}

// GetData get data from database.
func (s *Storage) GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error) {
	d := make([]byte, 0, ds.DataSize)
	err := s.QueryRow(ctx, sql.SelectData, ds.ID).Scan(&d)
	if err != nil {
		return nil, err
	}
	return s.crypto.Decrypt(d)
}

// GetSpecs get specs from database.
func (s *Storage) GetSpecs(ctx context.Context, userID uuid.UUID) ([]models.Spec, error) {
	var count int
	if err := s.QueryRow(ctx, "SELECT COUNT(*) FROM records WHERE user_id = $1", userID).Scan(&count); err != nil {
		return nil, err
	}
	rows, err := s.Query(ctx, sql.SelectSpecs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scan(ctx, rows, count)
}

// GetSpecsOfType get specs of different types from database.
func (s *Storage) GetSpecsOfType(ctx context.Context, userID uuid.UUID, mType models.MType) ([]models.Spec, error) {
	var count int
	if err := s.QueryRow(ctx, "SELECT COUNT(*) FROM records WHERE user_id = $1 AND type = $2", userID, mType).Scan(&count); err != nil {
		return nil, err
	}
	rows, err := s.Query(ctx, sql.SelectSpecsOfType, userID, mType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scan(ctx, rows, count)
}

func scan(_ context.Context, rows pgx.Rows, count int) ([]models.Spec, error) {
	specs := make([]models.Spec, 0, count)
	for rows.Next() {
		s := models.Spec{}
		if err := rows.Scan(&s.ID, &s.Type, &s.Title, &s.DataSize); err != nil {
			return specs, err
		}
		specs = append(specs, s)
	}
	return specs, rows.Err()
}
