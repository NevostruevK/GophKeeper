package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddData(t *testing.T) {
	text := models.NewText([]byte("some text"), false)

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
		meta := models.NewMeta(user.ID, text.Type(), "some_title")
		id, err := st.AddData(ctx, text, meta, nil)
		require.NoError(t, err)
		assert.Equal(t, meta.ID, id)
	})

	t.Run("Add data for unknown user error", func(t *testing.T) {
		meta := models.NewMeta(uuid.New(), models.TEXT, "some_title")
		id, err := st.AddData(ctx, text, meta, nil)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestStorage_GetData(t *testing.T) {
	text := models.NewText([]byte("some text"), false)
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

		meta := models.NewMeta(user.ID, text.Type(), "some_title")
		_, err = st.AddData(ctx, text, meta, nil)
		require.NoError(t, err)

		dataDB, err := st.GetData(ctx, meta)
		require.NoError(t, err)

		actual := &models.Text{}
		err = actual.Decode(dataDB)
		require.NoError(t, err)
		assert.Equal(t, text, actual)
	})
	t.Run("Get data for non-existent ID", func(t *testing.T) {
		meta := models.NewMeta(uuid.New(), text.Type(), "some_title")
		_, err := st.GetData(ctx, meta)
		assert.Error(t, err)
	})
}
