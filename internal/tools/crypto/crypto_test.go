package crypto

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	t.Run("test ok", func(t *testing.T) {
		data := []byte("some data for test crypto package")
		key := lo.RandomString(32, lo.AllCharset)
		nonce := lo.RandomString(12, lo.AllCharset)
		c, err := NewCrypto([]byte(key), []byte(nonce))
		require.NoError(t, err)
		encrypted := c.Crypt(data)
		decrypted, err := c.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, data, decrypted)
	})
	t.Run("key err", func(t *testing.T) {
		key := lo.RandomString(30, lo.AllCharset)
		nonce := lo.RandomString(12, lo.AllCharset)
		_, err := NewCrypto([]byte(key), []byte(nonce))
		require.Error(t, err)
	})
}
