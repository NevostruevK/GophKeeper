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

func TestMemory_AddUser(t *testing.T) {
	st := memory.NewUserStore()
	t.Run("add user ok", func(t *testing.T) {
		id1, err := st.AddUser(context.Background(), *models.NewUser("test_login1", "test_password1"))
		require.NoError(t, err)
		id2, err := st.AddUser(context.Background(), *models.NewUser("test_login2", "test_password2"))
		require.NoError(t, err)
		assert.NotEqual(t, id1, id2)
	})
	t.Run("add the same user err", func(t *testing.T) {
		id3, err := st.AddUser(context.Background(), *models.NewUser("test_login3", "test_password3"))
		require.NoError(t, err)
		assert.NotEqual(t, id3, uuid.Nil)
		id4, err := st.AddUser(context.Background(), *models.NewUser("test_login3", "test_password4"))
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrDuplicateLogin)
		assert.Equal(t, id4, uuid.Nil)
	})
}

func TestMemory_GetUser(t *testing.T) {
	st := memory.NewUserStore()
	t.Run("get user ok", func(t *testing.T) {
		id1, err := st.AddUser(context.Background(), *models.NewUser("test_login1", "test_password1"))
		require.NoError(t, err)
		id2, err := st.GetUser(context.Background(), *models.NewUser("test_login1", "test_password1"))
		require.NoError(t, err)
		assert.Equal(t, id1, id2)
	})
	t.Run("no user err", func(t *testing.T) {
		id, err := st.GetUser(context.Background(), *models.NewUser("non-existent_login", "test_password"))
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, id, uuid.Nil)
	})
	t.Run("wrong password err", func(t *testing.T) {
		_, err := st.AddUser(context.Background(), *models.NewUser("test_login4", "test_password4"))
		require.NoError(t, err)
		id, err := st.GetUser(context.Background(), *models.NewUser("test_login4", "wrong_password"))
		require.Error(t, err)
		assert.ErrorIs(t, err, storage.ErrWrongPassword)
		assert.Equal(t, id, uuid.Nil)
	})
}
