package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_GetData(t *testing.T) {
	dOut := models.Data([]byte("some data"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Run("Get data ok", func(t *testing.T) {
		id, err := addUser(ctx)
		require.NoError(t, err)

		r := models.NewRecord(models.TEXT, "some_title", dOut)
		ds, err := testStorage.AddRecord(ctx, id, r)
		require.NoError(t, err)

		dIn, err := testStorage.GetData(ctx, ds)
		require.NoError(t, err)
		assert.Equal(t, dOut, dIn)
	})
	t.Run("Get data for non-existent ID", func(t *testing.T) {
		ds := &models.DataSpec{ID: uuid.New()}
		_, err := testStorage.GetData(ctx, ds)
		assert.Error(t, err)
	})
}
