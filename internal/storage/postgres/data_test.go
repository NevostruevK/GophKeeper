package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	func TestStorage_AddData(t *testing.T) {
		text := models.NewText([]byte("some text"))

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		st, err := newStorage(ctx)
		require.NoError(t, err)

		ids := idsDB{make([]uuid.UUID, 0, 4)}
		defer func() {
			require.NoError(t, deleteFromDB(ctx, st, ids.ids))
		}()

		t.Run("Add Data ok", func(t *testing.T) {
			user, err := addUser(ctx, st, &ids)
			require.NoError(t, err)
			r := models.NewRecord(user.ID, text.Type(), "some_title")
			id, err := st.AddData(ctx, text, r, nil)
			require.NoError(t, err)
			assert.Equal(t, r.ID, id)
		})

		t.Run("Add data for non-existent user error", func(t *testing.T) {
			r := models.NewRecord(uuid.New(), models.TEXT, "some_title")
			id, err := st.AddData(ctx, text, r, nil)
			assert.Error(t, err)
			assert.Equal(t, uuid.Nil, id)
		})
	}
*/
func TestStorage_GetData(t *testing.T) {
	dOut := models.Data([]byte("some data"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()
	t.Run("Get data ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		r := models.NewRecord(models.TEXT, "some_title", dOut, nil)
		ds, err := st.AddRecord(ctx, user.ID, r)
		require.NoError(t, err)

		dIn, err := st.GetData(ctx, ds)
		require.NoError(t, err)
		assert.Equal(t, dOut, dIn)
	})
	t.Run("Get data for non-existent ID", func(t *testing.T) {
		ds := &models.DataSpec{ID: uuid.New()}
		_, err := st.GetData(ctx, ds)
		assert.Error(t, err)
	})
}
