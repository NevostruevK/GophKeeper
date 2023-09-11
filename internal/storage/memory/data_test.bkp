package memory_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemory_GetData(t *testing.T) {
	st := memory.NewDataStore()
	t.Run("Get text ok", func(t *testing.T) {
		dOut := models.Data([]byte("some text"))
		ds := &models.DataSpec{ID: uuid.New()}
		st.AddData(context.Background(), ds.ID, dOut)

		dIn, err := st.GetData(context.Background(), ds)
		require.NoError(t, err)
		assert.Equal(t, dOut, dIn)
	})
}
