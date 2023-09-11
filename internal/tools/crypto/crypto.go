package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

type Crypto struct {
	cipher.AEAD
	nonce []byte
}

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

func (c *Crypto) Crypt(src []byte) []byte {
	return c.Seal(nil, c.nonce, src, nil)
}

func (c *Crypto) Encrypt(src []byte) ([]byte, error) {
	return c.Open(nil, c.nonce, src, nil)
}
