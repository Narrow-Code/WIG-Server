package config

import (
	"github.com/joho/godotenv"
        "os"
	"github.com/dgrijalva/jwt-go"

)

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

