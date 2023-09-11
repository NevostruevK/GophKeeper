package postgres

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

/*
const(
	key = "secretKeyForGophKeeper"
	nonce = "_GophKeeper_"
)
*/
/*
type dataDB []byte

func newDataDB(size int) dataDB {
	return make([]byte, 0, size)
}
*/
/*
func dataToDB(d models.Data) ([]byte) {
    aesblock, err := aes.NewCipher([]byte(key))
    if err != nil {
        return nil, err
    }

    aesgcm, err := cipher.NewGCM(aesblock)
    if err != nil {
        return nil, err
    }
	dst := aesgcm.Seal(nil, []byte(nonce), d, nil)
	return dst, nil
	return
	// TODO зашифровать данные
//	return []byte(d), nil
}
*/
func (s *Storage) crypt(d models.Data) []byte {
	return s.crypto.Crypt(d)
}

/*
	func (db dataDB) toData() (models.Data, error) {
		// TODO расшифровать данные
		return models.Data(db), nil
	}
*/
func (s *Storage) AddRecord(ctx context.Context, userID uuid.UUID, r *models.Record) (*models.DataSpec, error) {
	/*
	   d, err := dataToDB(r.Data)

	   	if err != nil {
	   		return nil, err
	   	}
	*/d := s.crypt(r.Data)
	ds := &models.DataSpec{DataSize: len(d)}
	err := s.QueryRow(ctx, sql.InsertRecord, userID, r.Type, r.Title, d, ds.DataSize).Scan(&ds.ID)
	return ds, err
}

func (s *Storage) GetData(ctx context.Context, ds *models.DataSpec) (models.Data, error) {
	//	d := newDataDB(ds.DataSize)
	d := make([]byte, 0, ds.DataSize)
	err := s.QueryRow(ctx, sql.SelectData, ds.ID).Scan(&d)
	if err != nil {
		return nil, err
	}
	return s.crypto.Encrypt(d)
	// return d.toData()
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
		if err := rows.Scan(&s.ID, &s.Type, &s.Title, &s.DataSize); err != nil {
			return specs, err
		}
		specs = append(specs, s)
	}
	return specs, rows.Err()
}
