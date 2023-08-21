package models_test

import (
	"errors"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsUser_CountHash(t *testing.T) {
	t.Run("count hash ok", func(t *testing.T) {
		user, err := models.NewUserDB("test_login", "test_password")
		require.NoError(t, err)
		err = user.CountHash()
		require.NoError(t, err)
		hash, err := models.CountHash(*user)
		require.NoError(t, err)
		assert.True(t, len(hash) == models.HashSize)
		assert.Equal(t, hash, user.Hash)
	})
}

func TestModelsUser_CheckHash(t *testing.T) {
	t.Run("check hash ok", func(t *testing.T) {
		user, err := models.NewUserDB("test_login", "test_password")
		require.NoError(t, err)
		err = user.CheckHash()
		require.NoError(t, err)
	})
	t.Run("wrong hash error", func(t *testing.T) {
		user1, err := models.NewUserDB("test_login", "test_password")
		require.NoError(t, err)
		user2, err := models.NewUserDB("test_login", "another_password")
		require.NoError(t, err)
		user1.Hash = user2.Hash
		err = user1.CheckHash()
		assert.Error(t, err)
		assert.True(t, errors.Is(err, models.ErrWrongHashForUser))
	})
}
