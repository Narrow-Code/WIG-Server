package verification

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	algorithm   = "PBKDF2WithHmacSHA512"
	iterations  = 120000
	keyLength   = 32
	secret      = "JesusIsKing"
	saltLength  = 16
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
	saltString := ""
	for _, b := range salt {
		saltString += fmt.Sprintf("%02x", b)
	}
	saltAndSecret := saltString + secret

	// Generate PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), []byte(saltAndSecret), iterations, keyLength, sha512.New)

	return hex.EncodeToString(hash), nil
}

