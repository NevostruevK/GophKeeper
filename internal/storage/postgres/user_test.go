package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()
	t.Run("Add user ok", func(t *testing.T) {
		user1, err := addUser(ctx, st, &ids)
		require.NoError(t, err)
		user2, err := addUser(ctx, st, &ids)
		require.NoError(t, err)
		assert.NotEqual(t, user1.ID, user2.ID)
	})
	t.Run("Add the same login error", func(t *testing.T) {
		user, err := addUser(ctx, st, &ids)
		require.NoError(t, err)

		id, err := st.AddUser(ctx, models.User{Login: user.Login, Password: "any password"})
		require.Error(t, err)
		assert.True(t, errors.Is(err, storage.ErrDuplicateLogin))
		assert.Equal(t, uuid.Nil, id)
	})
	t.Run("Add short login error", func(t *testing.T) {
		id, err := st.AddUser(ctx, models.User{Login: "short", Password: "any password"})
		require.Error(t, err)
		assert.False(t, errors.Is(err, storage.ErrDuplicateLogin))
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestPostgres_GetUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, err := newStorage(ctx)
	require.NoError(t, err)

	ids := idsDB{make([]uuid.UUID, 0, 4)}
	defer func() {
		require.NoError(t, deleteFromDB(ctx, st, ids.ids))
	}()
	t.Run("Get user ok", func(t *testing.T) {
		user := models.NewUser("test_login", "test_password")
		idAdd, err := st.AddUser(ctx, *user)
		require.NoError(t, err)
		ids.ids = append(ids.ids, idAdd)
		idGet, err := st.GetUser(ctx, *user)
		require.NoError(t, err)
		assert.Equal(t, idAdd, idGet)
	})
}
