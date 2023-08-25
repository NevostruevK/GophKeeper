package postgres

import (
	"context"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/*
func (s *Storage) AddRecord(ctx context.Context, r *models.Record) (id uuid.UUID, err error) {
	data := []byte("data")
	r.SetDataSize(len(data))
	return id, s.QueryRow(ctx, sql.InsertRecordWithoutDescription, r.ID, r.UserID, r.MType, r.Title, data, r.DataSize).Scan(&id)
}
*/

/*
func (s *Storage) AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (id uuid.UUID, err error) {
	d, err := DataToDB(data)
	if err != nil {
		return uuid.Nil, err
	}
	r.SetDataSize(len(d.Data))

	if dsc != nil {
		err = s.QueryRow(ctx, sql.InsertRecord, r.ID, r.UserID, data.Type(), r.Title, d.Data, r.DataSize, dsc.Description).Scan(&id)
	} else {
		err = s.QueryRow(ctx, sql.InsertRecordWithoutDescription, r.ID, r.UserID, data.Type(), r.Title, d.Data, r.DataSize).Scan(&id)
	}
	return id, err
}
*/

type dataDB []byte

func newDataDB(size int) dataDB {
	return make([]byte, 0, size)
}

func dataToDB(d models.Data) (dataDB, error) {
	// TODO зашифровать данные
	return dataDB(d), nil
}

func (db dataDB) toData() (models.Data, error) {
	// TODO расшифровать данные
	return models.Data(db), nil
}

func (s *Storage) AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (*models.DataSpec, error) {
	d, err := dataToDB(r.Data)
	if err != nil {
		return nil, err
	}
	ds := &models.DataSpec{DataSize: len(d)}
	if r.Description != nil {
		//		INSERT INTO records (id, user_id, type, title, data, data_size, has_description, description)
		//		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'true', $6)
		err = s.QueryRow(ctx, sql.InsertRecord, userID, r.Type, r.Title, d, ds.DataSize, r.Description).Scan(&ds.ID)
	} else {
		//		INSERT INTO records (id, user_id, type, title, data, data_size, has_description)
		//		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'false')
		err = s.QueryRow(ctx, sql.InsertRecordWithoutDescription, userID, r.Type, r.Title, []byte(d), ds.DataSize).Scan(&ds.ID)
	}
	return ds, err
}

func (s *Storage) GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error) {
	d := newDataDB(ds.DataSize)
	//	dd := make([]byte, ds.DataSize)
	err := s.QueryRow(ctx, sql.SelectData, ds.ID).Scan(&d)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//	return models.Data(dd), nil
	return d.toData()
}

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
		if err := rows.Scan(&s.ID, &s.Type, &s.Title, &s.DataSize, &s.HasDescription); err != nil {
			return specs, err
		}
		specs = append(specs, s)
	}
	return specs, rows.Err()
}
