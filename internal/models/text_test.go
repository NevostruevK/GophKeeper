package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModelsText(t *testing.T) {
	t.Run("test String", func(t *testing.T) {
		text := models.NewText([]byte("text"))
		s := text.String()
		assert.Equal(t, "text size 4", s)
	})
	t.Run("test Show", func(t *testing.T) {
		text := models.NewText([]byte("text"))
		s := text.Show()
		assert.Equal(t, "text", s)
	})
}

func TestModelsText_IsReadyForStorage(t *testing.T) {
	const textIsEmpty = "text is empty"
	type result struct {
		bool
		string
	}
	tests := []struct {
		name string
		obj  models.Text
		want result
	}{
		{
			name: "test ok",
			obj:  *models.NewText([]byte("text")),
			want: result{true, ""},
		},
		{
			name: "test err",
			obj:  *models.NewText(nil),
			want: result{false, textIsEmpty},
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
