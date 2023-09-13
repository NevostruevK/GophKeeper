package memory

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

// DataStore provides storage for data in memory.
type DataStore struct {
	SpecCount int64    // conunt of specs
	Spec      sync.Map // map for specs
	Data      sync.Map // map for entries
}

// NewDataStore returns DataStore.
func NewDataStore() *DataStore {
	return &DataStore{
		Spec: sync.Map{},
		Data: sync.Map{},
	}
}

// AddSpecs strore specs.
func (s *DataStore) AddSpecs(specs []models.Spec) {
	for _, spec := range specs {
		s.Spec.Store(spec.ID, spec)
		atomic.AddInt64(&s.SpecCount, 1)
	}
}

// GetSpecs load specs.
func (s *DataStore) GetSpecs(_ context.Context, _ uuid.UUID) ([]models.Spec, error) {
	specs := make([]models.Spec, 0, atomic.LoadInt64(&s.SpecCount))
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

// GetSpecsOfType load specs for different types.
func (s *DataStore) GetSpecsOfType(_ context.Context, _ uuid.UUID, mType models.MType) ([]models.Spec, error) {
	specs := make([]models.Spec, 0, atomic.LoadInt64(&s.SpecCount))
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

// AddRecord store record.
func (s *DataStore) AddRecord(_ context.Context, _ uuid.UUID, r *models.Record) (*models.DataSpec, error) {
	id := uuid.New()
	s.Data.Store(id, r.Data)
	ds := &models.DataSpec{ID: id, DataSize: len(r.Data)}
	if _, ok := s.Spec.LoadOrStore(id, *(r.ToSpec(*ds))); !ok {
		atomic.AddInt64(&s.SpecCount, 1)
	}
	return &models.DataSpec{ID: id, DataSize: len(r.Data)}, nil
}
