package memory_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory_AddSpecs(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Add records ok", func(t *testing.T) {
		titles := []string{"title1", "title2", "title3"}
		specs := make([]models.Spec, 0, len(titles))
		for _, t := range titles {
			spec := models.NewSpec(models.TEXT, t)
			spec.ID = uuid.New()
			specs = append(specs, *spec)
		}
		st.AddSpecs(specs)
		assert.Equal(t, len(titles), int(atomic.LoadInt64(&st.SpecCount)))
	})
}

func TestMemory_GetSpecs(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Get Specs ok", func(t *testing.T) {
		titles := []string{"title1", "title2", "title3"}
		sOut := make([]models.Spec, 0, len(titles))
		for _, t := range titles {
			spec := models.NewSpec(models.TEXT, t)
			spec.ID = uuid.New()
			sOut = append(sOut, *spec)
		}
		st.AddSpecs(sOut)
		sIn, err := st.GetSpecs(context.Background(), uuid.New())
		require.NoError(t, err)
		assert.ElementsMatch(t, sOut, sIn)
	})
	t.Run("Get Specs err", func(t *testing.T) {
		st.Spec.Store(uuid.New(), "TypeAssert value")
		_, err := st.GetSpecs(context.Background(), uuid.New())
		require.Error(t, err)
		assert.ErrorIs(t, err, memory.ErrTypeAssert)
	})
}

func TestMemory_GetSpecsOfType(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Get Specs ok", func(t *testing.T) {
		titles := []string{"title1", "title2", "title3"}
		sOut := make([]models.Spec, 0, len(titles))
		for _, t := range titles {
			spec := models.NewSpec(models.TEXT, t)
			spec.ID = uuid.New()
			sOut = append(sOut, *spec)
		}
		st.AddSpecs(sOut)
		sIn, err := st.GetSpecsOfType(context.Background(), uuid.New(), models.TEXT)
		require.NoError(t, err)
		assert.ElementsMatch(t, sOut, sIn)
	})
	t.Run("Get Specs err", func(t *testing.T) {
		st.Spec.Store(uuid.New(), "TypeAssert value")
		_, err := st.GetSpecsOfType(context.Background(), uuid.New(), models.TEXT)
		require.Error(t, err)
		assert.ErrorIs(t, err, memory.ErrTypeAssert)
	})
}
