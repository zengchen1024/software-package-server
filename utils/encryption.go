package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

type Config struct {
	EncryptionKey string `json:"encryption_key"  required:"true"`
}

type SymmetricEncryption interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) ([]byte, error)
}

var Encryption SymmetricEncryption

func InitEncryption(key string) (err error) {
	Encryption, err = newSymmetricEncryption(key, "")

	return
}

func newSymmetricEncryption(key, nonce string) (SymmetricEncryption, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	se := symmetricEncryption{aead: gcm}

	if nonce != "" {
		nonce1, err := hex.DecodeString(nonce)
		if err != nil {
			return nil, err
		}
		if len(nonce1) != gcm.NonceSize() {
			return nil, errors.New(
				"the length of nonce for symmetric encryption is unmatched",
			)
		}
		se.nonce = nonce1
	}

	return &se, nil
}

type symmetricEncryption struct {
	aead  cipher.AEAD
	nonce []byte
}

func (se *symmetricEncryption) Encrypt(plaintext []byte) (string, error) {
	nonce := se.nonce
	if nonce == nil {
		nonce = make([]byte, se.aead.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return "", err
		}
	}

	return base64.StdEncoding.EncodeToString(
		se.aead.Seal(nonce, nonce, plaintext, nil),
	), nil
}

func (se *symmetricEncryption) Decrypt(ciphertext string) ([]byte, error) {
	content, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	nonceSize := se.aead.NonceSize()
	if len(content) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, text := content[:nonceSize], content[nonceSize:]

	return se.aead.Open(nil, nonce, text, nil)
}
