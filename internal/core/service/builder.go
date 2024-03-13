package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/pdaccess/pvault/internal/core/ports"
)

type ServiceImpl struct {
	storageGSM  cipher.AEAD
	persistence ports.Persistance
	validators  []ports.TokenValidator
}

func New(storeageKey []byte, persistence ports.Persistance, tokenValidators ...ports.TokenValidator) (ports.Service, error) {
	block, err := aes.NewCipher(storeageKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	return &ServiceImpl{
		storageGSM:  aesGCM,
		persistence: persistence,
		validators:  tokenValidators,
	}, nil
}

func (s *ServiceImpl) encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, s.storageGSM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	ciphertext := s.storageGSM.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext, nil
}

func (s *ServiceImpl) decrypt(enc []byte) ([]byte, error) {
	nonceSize := s.storageGSM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := s.storageGSM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}
