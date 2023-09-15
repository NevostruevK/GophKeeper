package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	owner    = "owner"
	bank     = "bank"
	number   = 123456789123456789
	expMonth = 10
	expYear  = 2030
	cvv      = 123
)

func TestModelsCard(t *testing.T) {
	t.Run("test String", func(t *testing.T) {
		card := models.NewCard(owner, bank, number, expMonth, expYear, cvv)
		s := card.String()
		assert.Equal(t, "owner : bank : 123456789123456789 : 10/2030", s)
	})
	t.Run("test Show", func(t *testing.T) {
		card := models.NewCard(owner, bank, number, expMonth, expYear, cvv)
		s := card.Show()
		assert.Equal(t, " bank\n 123456789123456789\n valid 10/2030  CVV 123\n owner", s)
	})
}

func TestModelsCard_IsReadyForStorage(t *testing.T) {
	const (
		wrongCardOwner           = "wrong owner"
		wrongCardExpirationMonth = "wrong month of expiration (should be from 1 to 12)"
		wrongCardExpirationYear  = "wrong year of expiration (should be more then 2000)"
		wrongCardNumber          = "wrong card number"
		wrongCardCVV             = "wrong card CVV"
	)
	type result struct {
		bool
		string
	}
	tests := []struct {
		name string
		obj  models.Card
		want result
	}{
		{
			name: "test ok",
			obj:  *models.NewCard(owner, bank, number, expMonth, expYear, cvv),
			want: result{true, ""},
		},
		{
			name: "test err (wrongCardOwner)",
			obj:  *models.NewCard("", bank, number, expMonth, expYear, cvv),
			want: result{false, wrongCardOwner},
		},
		{
			name: "test err (wrongCardExpirationMonth)",
			obj:  *models.NewCard(owner, bank, number, 0, expYear, cvv),
			want: result{false, wrongCardExpirationMonth},
		},
		{
			name: "test err (wrongCardExpirationYear)",
			obj:  *models.NewCard(owner, bank, number, expMonth, 1999, cvv),
			want: result{false, wrongCardExpirationYear},
		},
		{
			name: "test err (wrongCardNumber)",
			obj:  *models.NewCard(owner, bank, 0, expMonth, expYear, cvv),
			want: result{false, wrongCardNumber},
		},
		{
			name: "test err (wrongCardCVV)",
			obj:  *models.NewCard(owner, bank, number, expMonth, expYear, 0),
			want: result{false, wrongCardCVV},
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
