package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsUser_CheckHash(t *testing.T) {
	t.Run("check hash ok", func(t *testing.T) {
		password := "test_password"
		user, err := models.NewUserDB("test_login", password)
		require.NoError(t, err)
		err = user.CheckHash(password)
		require.NoError(t, err)
	})
	t.Run("wrong hash error", func(t *testing.T) {
		login := "test_login"
		password1 := "test_password"
		password2 := "another_password"
		user1, err := models.NewUserDB(login, password1)
		require.NoError(t, err)
		user2, err := models.NewUserDB(login, password2)
		require.NoError(t, err)
		user1.Hash = user2.Hash
		err = user1.CheckHash(password1)
		assert.Error(t, err)
	})
}
