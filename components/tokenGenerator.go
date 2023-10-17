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
func generateToken() (string, error) {
         godotenv.Load()
         var secret = []byte(os.Getenv("TOKEN_SECRET"))

        // Generate access token
        token := jwt.New(jwt.SigningMethodHS256)
        tokenStr, err := token.SignedString(secret)

        // Return error if access token generation fails
        if err != nil {
		return tokenStr, err
        }
	return tokenStr, nil
}

