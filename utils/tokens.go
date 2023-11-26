// Provides functions that may be regularly used throughout the WIG-Application.
package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

/*
* Generates a randomized authentication token for API calls between the WIG-Application and server.
*
* @param username The username.
* @param hash The hashed password.
*
* @return string The generated authentication token.
 */
func GenerateToken(username string, hash string) string {
	godotenv.Load()
	var secret = []byte(os.Getenv("TOKEN_SECRET"))

	// Generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"hash":     hash,
	})
	tokenStr, err := token.SignedString(secret)

	// Return error if access token generation fails
	if err != nil {
		return "error"
	}
	return tokenStr
}
