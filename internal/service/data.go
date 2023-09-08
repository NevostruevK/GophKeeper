package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
)

const serviceTimeOut = time.Second

func (s *Service) GetData(ctx context.Context, spec models.Spec) (models.Entry, error) {
	ctx, cancel := context.WithTimeout(ctx, serviceTimeOut)
	defer cancel()
	ds := models.DataSpec{ID: spec.ID, DataSize: spec.DataSize}

	entry, err := s.memory.GetEntry(ctx, &ds)
	if !errors.Is(err, storage.ErrNotFound) {
		return entry, err
	}

	data, err := s.client.Keeper.GetData(ctx, ds)
	if err != nil {
		return nil, err
	}

	entry = models.NewEntry(spec.Type)
	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(entry)
	if err != nil {
		return nil, err
	}

	s.memory.AddEntry(ctx, spec.ID, entry)
	return entry, err
}
