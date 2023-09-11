package models_test

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsText_Encode(t *testing.T) {
	msg := []byte("some text")

	t.Run("encode is idempotent", func(t *testing.T) {
		text1 := models.NewText(msg)
		data1, err := text1.Encode()
		require.NoError(t, err)
		text2 := models.NewText(msg)
		data2, err := text2.Encode()
		require.NoError(t, err)
		assert.Equal(t, data1, data2)
	})
	t.Run("encode is different", func(t *testing.T) {
		text1 := models.NewText(msg)
		data1, err := text1.Encode()
		require.NoError(t, err)
		msg = append(msg, []byte("append")...)
		text2 := models.NewText(msg)
		data2, err := text2.Encode()
		require.NoError(t, err)
		assert.NotEqual(t, data1, data2)
	})
}

func TestModelsText_Decode(t *testing.T) {
	msg := []byte("some text")
	t.Run("decode ok", func(t *testing.T) {
		text := models.NewText(msg)
		data, err := text.Encode()
		require.NoError(t, err)
		dText := &models.Text{}
		require.NoError(t, dText.Decode(data))
		assert.Equal(t, text, dText)
	})
}
