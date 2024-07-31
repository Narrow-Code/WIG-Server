package verification

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

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
func GenerateSalt() (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, byteValue := range salt {
		sb.WriteString(fmt.Sprintf("%02x", byteValue))
	}

	fmt.Printf("Salt length: %d\n", len(sb.String()))
	
	return sb.String(), nil
}

// GenerateHash generates a PBKDF2 hash of the password with the provided salt
func GenerateHash(password string, salt string) (string, error) {
	// Combine salt and secret
	saltAndSecret := salt + secret

	// Generate PBKDF2 hash
	hash := pbkdf2.Key([]byte(password), []byte(saltAndSecret), iterations, keyLength, sha512.New)

	return hex.EncodeToString(hash), nil
}

