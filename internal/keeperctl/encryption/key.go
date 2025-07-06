package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

const NonceLength = 12

// Key is user's encryption key.
type Key struct {
	sum [sha256.Size]byte
}

// NewKey creates new encryption key.
func NewKey(username string, password creds.Password) Key {
	sum := sha256.Sum256([]byte(username + "@" + string(password)))

	return Key{sum: sum}
}

// Hash provides hash of the encryption key.
func (k Key) Hash() string {
	return hex.EncodeToString(k.sum[:])
}

// Encrypt encrypts provided message.
func (k Key) Encrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	gcm, err := k.getGCM()
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("ReadFull error: %w", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

// Decrypt decrypts provided data.
func (k Key) Decrypt(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	gcm, err := k.getGCM()
	if err != nil {
		return nil, err
	}

	nonce, payload := data[:NonceLength], data[NonceLength:]

	decrypted, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("can't decript key: %w", err)
	}

	return decrypted, nil
}

func (k Key) getGCM() (cipher.AEAD, error) {
	cipherBlock, err := aes.NewCipher(k.sum[:])
	if err != nil {
		return nil, fmt.Errorf("newCipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, fmt.Errorf("newgcm error: %w", err)
	}

	return gcm, nil
}
