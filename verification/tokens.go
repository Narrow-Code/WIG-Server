// utils provides functions that may be regularly used throughout the WIG-Application.
package verification

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

/* GenerateVerificationToken generates a randomized UUID and expiration time.
* 
* @param user The user in which the token is being generated for.
* @return string The generated token
* @return time.Time the time in which the token expires.
*/
func GenerateVerificationToken(user models.User) (string, time.Time) {
	var emailVerification models.EmailVerification
	result := db.DB.Where("user_id = ?", user.UserUID).First(&emailVerification)

	if result != nil {
		db.DB.Delete(&emailVerification)
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour)

	emailVerification = models.EmailVerification{
		VerificationToken: 		token,
		EmailUserID:   	user.UserUID,
		EmailExpiresAt: 	expiresAt,
	}

	// Create user and return
	db.DB.Create(&emailVerification)

	return token, expiresAt
}

/* GeneratePasswordToken generates a randomized UUID and expiration time.
* 
* @param user The user in which the token is being generated for.
* @return string The generated token
* @return time.Time the time in which the token expires.
*/
func GeneratePassowrdToken(user models.User) (string, time.Time) {
	var passwordChange models.PasswordChange
	result := db.DB.Where("password_user_id = ?", user.UserUID).First(&passwordChange)

	if result != nil {
		db.DB.Delete(&passwordChange)
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour)

	passwordChange = models.PasswordChange{
		PasswordChangeToken: 	token,
		PasswordUserID:   	user.UserUID,
		PasswordExpiresAt: 	expiresAt,
	}

	// Create user and return
	db.DB.Create(&passwordChange)

	return token, expiresAt
}
