package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type EncryptionService struct {
	aesBlock cipher.Block
	gcm      cipher.AEAD
}

func NewEncryptionService() (*EncryptionService, error) {
	s := os.Getenv("CRYPTO_SECRET")

	sByte := []byte(s)

	if len(sByte) > 32 || len(sByte) < 1 {
		log.Fatal("Invalid Crypto Key")
		return nil, errors.New("invalid crypto key")
	}

	if len(sByte) < 32 {
		padding := make([]byte, 32-len(sByte))
		sByte = append(sByte, padding...)
	}

	aesBlock, err := aes.NewCipher(sByte)
	if err != nil {
		return nil, fmt.Errorf("something went wrong when creating aes block: %w", err)
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, fmt.Errorf("error setting GCM mode: %w", err)
	}

	return &EncryptionService{
		aesBlock: aesBlock,
		gcm:      gcm,
	}, nil
}

func (e EncryptionService) Encrypt(t string) (string, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("error generating the nonce ", err)
		return "", err
	}

	ciphertext := e.gcm.Seal(nonce, nonce, []byte(t), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func (e EncryptionService) Decrypt(encodedData string) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 data: %w", err)
	}

	nonceSize := e.gcm.NonceSize()
	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	plaintext, err := e.gcm.Open(nil, []byte(nonce), []byte(encryptedData), nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting data: %w", err)
	}

	return string(plaintext), nil

}
