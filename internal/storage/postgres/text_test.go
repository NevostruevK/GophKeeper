package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddText(t *testing.T) {
	text := models.NewText([]byte("some text"), false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()

	t.Run("Add text ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		meta := models.NewMeta(user.ID, models.TEXT, "some_title")
		id, err := st.AddText(ctx, text, meta, nil)
		require.NoError(t, err)
		assert.Equal(t, meta.ID, id)
	})

	t.Run("Add text for non-existent ID error", func(t *testing.T) {
		meta := models.NewMeta(uuid.New(), models.TEXT, "some_title")
		id, err := st.AddText(ctx, text, meta, nil)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestStorage_GetText(t *testing.T) {
	text := models.NewText([]byte("some text"), false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()

	t.Run("Get text ok", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		meta := models.NewMeta(user.ID, models.TEXT, "some_title")
		_, err = st.AddText(ctx, text, meta, nil)
		require.NoError(t, err)

		actual, err := st.GetText(ctx, meta)
		require.NoError(t, err)
		assert.Equal(t, text, actual)
	})

	t.Run("Get text for non-existent ID", func(t *testing.T) {
		meta := models.NewMeta(uuid.New(), text.Type(), "some_title")
		actual, err := st.GetText(ctx, meta)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}
