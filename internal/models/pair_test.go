package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModelsPair(t *testing.T) {
	t.Run("test String", func(t *testing.T) {
		pair := models.NewPair("Login", "Password")
		s := pair.String()
		assert.Equal(t, "Login: Login Password: Password", s)
	})
	t.Run("test Show", func(t *testing.T) {
		pair := models.NewPair("Login", "Password")
		s := pair.Show()
		assert.Equal(t, "Login: Login Password: Password", s)
	})
}

func TestModelsPair_IsReadyForStorage(t *testing.T) {
	const (
		loginIsEmpty    = "login is empty"
		passwordIsShort = "password is short (at least 4 characters)"
	)
	type result struct {
		bool
		string
	}
	tests := []struct {
		name string
		obj  models.Pair
		want result
	}{
		{
			name: "test ok",
			obj:  *models.NewPair("valid_login", "valid_password"),
			want: result{true, ""},
		},
		{
			name: "test err (login is empty)",
			obj:  *models.NewPair("", "valid_password"),
			want: result{false, loginIsEmpty},
		},
		{
			name: "test err (password is short (at least 4 characters))",
			obj:  *models.NewPair("valid_login", "pas"),
			want: result{false, passwordIsShort},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok, s := tt.obj.IsReadyForStorage(); ok != tt.want.bool || s != tt.want.string {
				t.Errorf("IsReadyForStorage(%v) got (%v , %v), want %v", tt.obj, ok, s, tt.want)
			}
		})
	}
}
