package verification

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/hmac"
	"encoding/hex"
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
	saltAndSecret := append(salt, []byte(secret)...)

	// Generate PBKDF2 hash
	hash := pbkdf2(password, saltAndSecret, iterations, keyLength)

	return hex.EncodeToString(hash), nil
}

// pbkdf2 implements PBKDF2 with HMAC-SHA512
func pbkdf2(password string, salt []byte, iterations, keyLength int) []byte {
	key := make([]byte, keyLength)
	hmac := hmac.New(sha512.New, []byte(password))
	for i := 0; i < iterations; i++ {
		hmac.Reset()
		hmac.Write(salt)
		hmac.Write(key)
		hmac.Sum(key[:0])
	}
	return key
}
