// package crypto provides encryption and decription of data.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// Crypto provides encryption and decription of data.
type Crypto struct {
	cipher.AEAD
	nonce []byte
}

// NewCrypto returns Crypto.
func NewCrypto(key, nonce []byte) (*Crypto, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	return &Crypto{aesgcm, nonce}, nil
}

// Crypt encrypt data.
func (c *Crypto) Crypt(src []byte) []byte {
	return c.Seal(nil, c.nonce, src, nil)
}

// Decrypt decrypt data.
func (c *Crypto) Decrypt(src []byte) ([]byte, error) {
	return c.Open(nil, c.nonce, src, nil)
}
