package memory

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

type DataStore struct {
	SpecLen     int64
	Spec        sync.Map
	Data        sync.Map
	Description sync.Map
}

func NewDataStore() *DataStore {
	return &DataStore{
		Spec:        sync.Map{},
		Data:        sync.Map{},
		Description: sync.Map{},
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

func (s *DataStore) AddRecord(_ context.Context, _ uuid.UUID, r *models.Record) (uuid.UUID, error) {
	id := uuid.New()
	s.Data.Store(id, r.Data)

	if _, ok := s.Spec.LoadOrStore(id, *(r.ToSpec(id))); !ok {
		atomic.AddInt64(&s.SpecLen, 1)
	}
	s.Description.Store(id, r.Description)
	return id, nil
}
