/*
* Package components provides functions that may be regularly used throughout the WIG-Application
*/
package components

import (
	"github.com/joho/godotenv"
        "os"
	"github.com/dgrijalva/jwt-go"

)

/*
* generateToken generates a randomized authentication token for API calls between the WIG-Application and server.
*
* @return string - The generated authentication token.
* @return error - An error, if any, during the token generation process.
*/
func GenerateToken(username string, salt string, email string) string {
         godotenv.Load()
         var secret = []byte(os.Getenv("TOKEN_SECRET"))

        // Generate access token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
	"salt": salt,
	"email": email,
    })
        tokenStr, err := token.SignedString(secret)

        // Return error if access token generation fails
        if err != nil {
		return ""
        }
	return tokenStr
}

