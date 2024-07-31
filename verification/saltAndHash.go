package verification

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

const (
	algorithm   = "PBKDF2WithHmacSHA512"
	iterations  = 120000
	keyLength   = 32
	secret      = "JesusIsKing"
	saltLength  = 23
)

// GenerateSalt generates a random salt of specified length
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateHash generates a PBKDF2 hash of the password with the provided salt
func GenerateHash(password string, salt []byte) (string, error) {
	// Combine salt and secret
	saltAndSecret := append(salt, []byte(secret)...)

	// Generate PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), saltAndSecret, iterations, keyLength, sha512.New)

	return hex.EncodeToString(hash), nil
}

