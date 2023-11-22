/*
* Package components provides functions that may be regularly used throughout the WIG-Application
*/
package utils

import (
	"github.com/joho/godotenv"
        "os"
	"github.com/dgrijalva/jwt-go"

)

/*
* generateToken generates a randomized authentication token for API calls between the WIG-Application and server.
*
* @param username The username of the user
* @param hash The hashed password of the user
*
* @return string - The generated authentication token.
* @return error - An error, if any, during the token generation process.
*/
func GenerateToken(username string, hash string) string {
         godotenv.Load()
         var secret = []byte(os.Getenv("TOKEN_SECRET"))

        // Generate access token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
	"hash": hash,
    })
        tokenStr, err := token.SignedString(secret)

        // Return error if access token generation fails
        if err != nil {
		return ""
        }
	return tokenStr
}

/*
* validateToken checks to see if a users token matches the protocol.
*
* @param username The username of the user
* @param hash The hashed password of the user
* @param token The authentication token to verify
* 
* @return bool - True if the token matches, false if not
*/
func ValidateToken(username string, hash string, token string) bool {
	if token == GenerateToken(username, hash) {
		return true
	} else {
		return false
	}
}
