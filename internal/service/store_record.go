package service

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/NevostruevK/GophKeeper/internal/models"
)

//func (s *Service) StoreRecord(ctx context.Context, re models.RecordEntry) (*models.DataSpec, error){
func (s *Service) StoreEntry(ctx context.Context, typ models.MType, title string, entry models.Entry) (*models.DataSpec, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(entry)
	if err != nil {
		return nil, err
	}
	//	return buff.Bytes(), err

	r := models.NewRecord(typ, title, buff.Bytes())
	ds, err := s.client.Keeper.AddRecord(ctx, r)
	if err != nil {
		return nil, err
	}
	s.memory.AddEntry(ctx, ds.ID, entry)
	s.memory.AddSpecs([]models.Spec{*r.ToSpec(*ds)})
	return ds, nil
}
