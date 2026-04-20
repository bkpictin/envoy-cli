// Package encrypt provides AES-GCM encryption and decryption for
// environment variable values stored in envoy config files.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/user/envoy-cli/internal/config"
)

const encryptedPrefix = "enc:"

// deriveKey produces a 32-byte AES-256 key from an arbitrary passphrase.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// EncryptValue encrypts plaintext with the given passphrase and returns a
// base64-encoded ciphertext prefixed with "enc:".
func EncryptValue(passphrase, plaintext string) (string, error) {
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptValue decrypts a value previously encrypted with EncryptValue.
// Returns an error if the value is not encrypted or decryption fails.
func DecryptValue(passphrase, ciphertext string) (string, error) {
	if !IsEncrypted(ciphertext) {
		return "", errors.New("value is not encrypted")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext[len(encryptedPrefix):])
	if err != nil {
		return "", err
	}
	key := deriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	nonce, data := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// IsEncrypted reports whether the value carries the encrypted prefix.
func IsEncrypted(value string) bool {
	return len(value) > len(encryptedPrefix) && value[:len(encryptedPrefix)] == encryptedPrefix
}

// EncryptTarget encrypts all plain-text values in the given target.
func EncryptTarget(cfg *config.Config, target, passphrase string) error {
	envs, ok := cfg.Targets[target]
	if !ok {
		return errors.New("target not found: " + target)
	}
	for k, v := range envs {
		if IsEncrypted(v) {
			continue
		}
		enc, err := EncryptValue(passphrase, v)
		if err != nil {
			return err
		}
		envs[k] = enc
	}
	return nil
}

// DecryptTarget decrypts all encrypted values in the given target.
func DecryptTarget(cfg *config.Config, target, passphrase string) error {
	envs, ok := cfg.Targets[target]
	if !ok {
		return errors.New("target not found: " + target)
	}
	for k, v := range envs {
		if !IsEncrypted(v) {
			continue
		}
		dec, err := DecryptValue(passphrase, v)
		if err != nil {
			return err
		}
		envs[k] = dec
	}
	return nil
}
