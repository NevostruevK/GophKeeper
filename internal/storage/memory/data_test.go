package memory_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory_GetEntry(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Get entry err (not found)", func(t *testing.T) {
		ds := &models.DataSpec{ID: uuid.New()}
		_, err := st.GetEntry(context.Background(), ds)
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
	t.Run("Get entry ok", func(t *testing.T) {
		out := models.NewText([]byte("some text"))
		ds := &models.DataSpec{ID: uuid.New()}
		st.AddEntry(context.Background(), ds.ID, out)

		in, err := st.GetEntry(context.Background(), ds)
		require.NoError(t, err)
		assert.Equal(t, out, in)
	})
	t.Run("Get entry err (type assert)", func(t *testing.T) {
		ds := &models.DataSpec{ID: uuid.New()}
		st.Data.Store(ds.ID, []byte("some text"))
		_, err := st.GetEntry(context.Background(), ds)
		require.Error(t, err)
		assert.ErrorIs(t, err, memory.ErrTypeAssert)
	})
}

func TestMemory_GetData(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Get data err", func(t *testing.T) {
		ds := &models.DataSpec{ID: uuid.New()}
		_, err := st.GetData(context.Background(), ds)
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrNotFound)
	})
	t.Run("Get data ok", func(t *testing.T) {
		out := models.Data([]byte("some text"))
		ds := &models.DataSpec{ID: uuid.New()}
		st.Data.Store(ds.ID, out)
		in, err := st.GetData(context.Background(), ds)
		require.NoError(t, err)
		assert.Equal(t, out, in)
	})

	t.Run("Get data ok", func(t *testing.T) {
		out := models.NewText([]byte("some text"))
		ds := &models.DataSpec{ID: uuid.New()}
		st.AddEntry(context.Background(), ds.ID, out)

		_, err := st.GetData(context.Background(), ds)
		require.Error(t, err)
		assert.ErrorIs(t, err, memory.ErrTypeAssert)
	})
}
