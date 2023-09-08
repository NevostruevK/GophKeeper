package memory

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

type DataStore struct {
	SpecLen int64
	Spec    sync.Map
	Data    sync.Map
}

func NewDataStore() *DataStore {
	return &DataStore{
		Spec: sync.Map{},
		Data: sync.Map{},
	}
}

func (s *DataStore) AddSpecs(specs []models.Spec) {
	for _, spec := range specs {
		s.Spec.Store(spec.ID, spec)
		atomic.AddInt64(&s.SpecLen, 1)
	}
}

func (s *DataStore) GetSpecs(_ context.Context, _ uuid.UUID) ([]models.Spec, error) {
	specs := make([]models.Spec, 0, atomic.LoadInt64(&s.SpecLen))
	var err error
	s.Spec.Range(func(key any, value any) bool {
		spec, ok := value.(models.Spec)
		if !ok {
			err = ErrTypeAssert
			return false
		}
		specs = append(specs, spec)
		return true
	})
	return specs, err
}

func (s *DataStore) GetSpecsOfType(_ context.Context, _ uuid.UUID, mType models.MType) ([]models.Spec, error) {
	specs := make([]models.Spec, 0, atomic.LoadInt64(&s.SpecLen))
	var err error
	s.Spec.Range(func(key any, value any) bool {
		spec, ok := value.(models.Spec)
		if !ok {
			err = ErrTypeAssert
			return false
		}
		if spec.Type == mType {
			specs = append(specs, spec)
		}
		return true
	})
	return specs, err
}

func (s *DataStore) AddRecord(_ context.Context, _ uuid.UUID, r *models.Record) (*models.DataSpec, error) {
	id := uuid.New()
	s.Data.Store(id, r.Data)
	ds := &models.DataSpec{ID: id, DataSize: len(r.Data)}
	if _, ok := s.Spec.LoadOrStore(id, *(r.ToSpec(*ds))); !ok {
		atomic.AddInt64(&s.SpecLen, 1)
	}
	return &models.DataSpec{ID: id, DataSize: len(r.Data)}, nil
}
