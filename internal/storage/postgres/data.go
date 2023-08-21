package postgres

import (
	"bytes"
	"context"
	"crypto/aes"
	"encoding/gob"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const (
	KeySize   = 2 * aes.BlockSize
	NonceSize = 12
)

type DataDB struct {
	Data  []byte
	Key   []byte
	Nonce []byte
}

func NewDataDB(size int) *DataDB {
	return &DataDB{
		Data:  make([]byte, 0, size),
		Key:   make([]byte, 0, KeySize),
		Nonce: make([]byte, 0, NonceSize),
	}
}

func DataToDB(data Data) (*DataDB, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	key := lo.RandomString(KeySize, lo.AllCharset)
	nonce := lo.RandomString(NonceSize, lo.AllCharset)

	// TODO зашифровать данные

	return &DataDB{buff.Bytes(), []byte(key), []byte(nonce)}, nil
}

func (d *DataDB) ToData() ([]byte, error) {
	// TODO расшифровать данные
	return d.Data, nil
}

func (s *Storage) AddData(ctx context.Context, data Data, m *models.Meta, dsc *models.Description) (id uuid.UUID, err error) {
	d, err := DataToDB(data)
	if err != nil {
		return uuid.Nil, err
	}

	m.SetDescriptionSize(dsc)
	m.SetDataSize(len(d.Data))

	tx, err := s.Begin(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, sql.InsertMeta, m.ID, m.UserID, data.Type(), m.Title, m.DataSize, m.DescriptionSize).Scan(&id)
	if err != nil {
		return
	}

	_, err = tx.Exec(ctx, sql.InsertData, id, d.Data, d.Key, d.Nonce)
	if err != nil {
		return
	}

	if m.DescriptionSize > 0 {
		_, err = tx.Exec(ctx, sql.InsertDescription, id, dsc.Description, dsc.IsCompressed)
		if err != nil {
			return
		}
	}
	return id, tx.Commit(ctx)
}

func (s *Storage) GetData(ctx context.Context, meta *models.Meta) ([]byte, error) {
	d := NewDataDB(meta.DataSize)
	err := s.QueryRow(ctx, sql.SelectData, meta.ID).Scan(&d.Data, &d.Key, &d.Nonce)
	if err != nil {
		return nil, err
	}
	return d.ToData()
}
