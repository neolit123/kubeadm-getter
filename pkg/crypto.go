// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func createHash(key []byte) []byte {
	hasher := sha256.New()
	hasher.Write(key)
	return hasher.Sum(nil)
}

func createRandBytes(sz uint32) ([]byte, error) {
	bytes := make([]byte, sz)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// CreateHandshakeBytes ...
func CreateHandshakeBytes() ([]byte, error) {
	bytes, err := createRandBytes(HandShakeSize)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CreateToken ...
func CreateToken() (string, error) {
	bytes, err := createRandBytes(32)
	if err != nil {
		return "", err
	}
	bytes = createHash(bytes)
	return hex.EncodeToString(bytes[:3]) + "." + hex.EncodeToString(bytes[3:11]), nil
}

// EncryptBytes ...
func EncryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := createRandBytes(uint32(gcm.NonceSize()))
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptBytes ...
func DecryptBytes(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return nil, fmt.Errorf("size of data is less than the nonce")
	}

	nonce, out := data[:nonceSize], data[nonceSize:]
	out, err = gcm.Open(nil, nonce, out, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}
