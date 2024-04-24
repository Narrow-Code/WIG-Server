// utils provides functions that may be regularly used throughout the WIG-Application.
package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

/*
* GenerateToken generates a randomized authentication token for API calls between the WIG-Application and server.
* 
* @param username The username.
* @param hash The hashed password.
* @return string The generated authentication token.
 */
func GenerateToken(username string, hash string) string {
	// Load environment variables
	godotenv.Load()

	// Get token secret from environment
	tokenSecret := []byte(os.Getenv("TOKEN_SECRET"))

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"hash":     hash,
	})

	// Sign token with secret
	tokenStr, err := token.SignedString(tokenSecret)

	// Return error if access token signing fails
	if err != nil {
		return "error"
	}

	return tokenStr
}
