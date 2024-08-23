package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func (params *HashParams) HashPasword(password string) (encodedHash string, err error) {

	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	// TODO PEPPER THIS MF AFTER HASING https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	// https://pkg.go.dev/crypto/hmac
	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// TODO remove all this rubish and only keep hash and salt
	hashWithDetails := fmt.Sprintf("$%s$%s", b64Salt, b64Hash)

	encodedHash = base64.RawStdEncoding.EncodeToString([]byte(hashWithDetails))

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (params *HashParams) ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	salt, hashedPassword, err := params.DecodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	hashedUserInput := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	return subtle.ConstantTimeCompare(hashedPassword, hashedUserInput) == 1, nil
}

func (params *HashParams) DecodeHash(encodedHash string) (salt, hash []byte, err error) {
	decodedHash, err := base64.RawStdEncoding.Strict().DecodeString(encodedHash)
	if err != nil {
		return nil, nil, err
	}

	values := strings.Split(string(decodedHash), "$")
	if len(values) != 3 {
		return nil, nil, ErrInvalidHash
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[1])
	if err != nil {
		return nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[2])
	if err != nil {
		return nil, nil, err
	}
	params.KeyLength = uint32(len(hash))

	return salt, hash, nil
}
