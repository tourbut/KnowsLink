// Sealer encrypts secrets the relay must store (webhook signing secrets, later routine keys) with AES-256-GCM.
package relay

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

// Sealer holds the KNOWSLINK_SEAL_KEY cipher; the key lives only in process memory, never in state or logs.
type Sealer struct{ aead cipher.AEAD }

const sealPrefix = "sealed:v1:"

// NewSealer reads a base64 32-byte key. An empty key returns nil so the features that need it stay off.
func NewSealer(key string) (*Sealer, error) {
	if key == "" {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(raw) != 32 {
		return nil, errors.New("KNOWSLINK_SEAL_KEY must be base64 of 32 bytes")
	}
	block, err := aes.NewCipher(raw)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Sealer{aead}, nil
}

// Seal binds the ciphertext to aad (the owning subscription or agent ID), so a value copied to another record does not open.
func (s *Sealer) Seal(plain, aad string) string {
	nonce := make([]byte, s.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	return sealPrefix + base64.RawStdEncoding.EncodeToString(s.aead.Seal(nonce, nonce, []byte(plain), []byte(aad)))
}

// Open fails without a key, on another record's aad, on tampering and on a value sealed under a different key.
func (s *Sealer) Open(sealed, aad string) (string, error) {
	raw, ok := strings.CutPrefix(sealed, sealPrefix)
	if s == nil || !ok {
		return "", fault("sealed_unreadable")
	}
	b, err := base64.RawStdEncoding.DecodeString(raw)
	n := s.aead.NonceSize()
	if err != nil || len(b) < n {
		return "", fault("sealed_unreadable")
	}
	plain, err := s.aead.Open(nil, b[:n], b[n:], []byte(aad))
	if err != nil {
		return "", fault("sealed_unreadable")
	}
	return string(plain), nil
}
