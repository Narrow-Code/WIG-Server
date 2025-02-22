package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"WIG-Server/verification"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// usernameRegex is the regex expression to check username requirements
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)

// emailRegex is the regex expression to check email requirements
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// UserSalt retrieves the users salt.
func UserSalt(c *fiber.Ctx) error {
	// Initialize variables
	utils.Log("call began")
	var user models.User
	username := c.Params("username")

	// Check if username is empty
	utils.Log("checking for empty fields")
	if username == "" {
		return Error(c, 400, "Username is empty and required")
	}

	// Query database for username
	utils.Log("query database for " + user.Username)
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}

	// Add to dto and return
	utils.Log("return for salt for " + user.Username + " successful")
	dto := DTO("salt", user.Salt)
	return success(c, "Salt returned successfully", dto)
}

// UserValidate validates the users token is still valid.
func UserValidate(c *fiber.Ctx) error {
	utils.UserLog(c, "user token authorized")
	return success(c, "Authorized")
}

// UserLogin handles the user Login logic. Returning a token.
func UserLogin(c *fiber.Ctx) error {
	// Initialize variables
	utils.Log("call began")
	var data map[string]string
	var user models.User
	
	// Parse JSON body
	utils.Log("parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	username := data["username"]
	hash := data["hash"]


	// Check for empty fields
	utils.Log("checking for empty fields")
	if username == "" || hash == "" {
		return Error(c, 400, "Username or hash is empty and required")
	}

	// Check that user exists
	utils.Log("query username")
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}

	// Check if hash matches
	utils.Log("validating hash")
	if hash != user.Hash {
		utils.Log("Stored hash: " + user.Hash)
		utils.Log("Sent in hash: " + hash)
		return Error(c, 400, "The username and passwords do not match")
	}

	// Generate token
	utils.Log("generating token for " + user.Username)
	user.Token = verification.GenerateToken(user.Username, user.Hash)
	if user.Token == "error" {
		return Error(c, 400, "There was an error generating user token")
	}

	// Check for email verification
	utils.Log("checking email verification")
	if user.EmailConfirm == "false" {
		return Error(c, 400, "Please verify email")
	}

	// Save to database, add to dto's and return
	db.DB.Save(&user)
	utils.Log("user log in successful")
	tokenDTO := DTO("token", user.Token)
	uidDTO := DTO("uid", user.UserUID)
	return success(c, "Login was successful", tokenDTO, uidDTO)
}

/*
* Handles user registration requests.
* It performs various checks such as data validation and database uniqueness before creating a new user record.
 */
func UserSignup(c *fiber.Ctx) error {
	// Initialize variables
	utils.Log("call began")
	var data map[string]string
	var user models.User

	// Parse request into data map
	utils.Log("parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}

	// Check for empty fields
	utils.Log("checking for empty fields")
	if data["username"] == "" || data["email"] == "" {
		return Error(c, 400, "Username or email is empty and required")
	}
	if data["salt"] == "" || data["hash"] == "" {
		return Error(c, 400, "Salt or hash is empty and required")
	}

	// Check if username is in use
	utils.Log("checking if username is in use")
	result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Username: "+err.Error())
	}

	// Check if email is in use
	utils.Log("checking if email is in use")
	result = db.DB.Where("email = ?", data["email"]).First(&user)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Email: "+err.Error())
	}

	// Check username and email requierments
	utils.Log("validating username and email match requirements")
	if !usernameRegex.MatchString(data["username"]) {
		return Error(c, 400, "Username does not match requirements")
	}
	if !emailRegex.MatchString(data["email"]) {
		return Error(c, 400, "Email does not match requirements")
	}

	// Run DNS check on Email
	utils.Log("running DNS check on email for verification")
	domain := strings.Split(data["email"], "@")[1]
	_, err = net.LookupMX(domain)
	if err != nil {
		return Error(c, 400, "Email domain does not exist")
	}
	
	// create user
	user = createUser(data)
	utils.Log("registration for " + user.Username + " was successful")
	
	// send verification email
        hosted := os.Getenv("SELF_HOSTED")
        if hosted == "false" {
		verification.SendVerificationEmail(user)
        }

	return success(c, "Signup was successful. Self Hosted: " + hosted)
}

// Ping performs a health check
func Ping(c *fiber.Ctx) error {
	utils.Log("health check performed")
	return success(c, "Ping was successful")
}

// ResendVerificationEmail will cancel the old verification email and resend it
func ResendVerificationEmail(c *fiber.Ctx) error {
	utils.Log("call began")
	var data map[string]string
	var user models.User

	// Parse JSON body
	utils.Log("parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	email := data["email"]

	// Check that user exists
	utils.Log("query username")
	result := db.DB.Where("email = ?", email).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}
	
	if user.EmailConfirm == "true" {
		return Error(c, 400, "Email already verified")
	}

	// Resend verification email
	verification.SendVerificationEmail(user)

	return success(c, "Verification email was resent")
}

// ResetPassword will send a Reset Password email to the user
func ResetPassword(c *fiber.Ctx) error {
	utils.Log("call began")
	var data map[string]string
	var user models.User
	
	// Parse JSON body
	utils.Log("parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	email := data["email"]

	// Check that user exists
	utils.Log("query username")
	result := db.DB.Where("email = ?", email).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}

	// Resend verification email
	verification.SendResetPasswordEmail(user)

	return success(c, "Reset password email was resent")
}

func VerificationEmail(c *fiber.Ctx) error {
	uid := c.Params("uid")
	
	var emailVerification models.EmailVerification
	var user models.User

	result := db.DB.Where("verification_token = ?", uid).First(&emailVerification)

	if result.Error != nil || result.RowsAffected == 0 || emailVerification.EmailExpiresAt.Before(time.Now()) {
		return c.SendString("The email verification link you followed does not exist or has expired.")
	}

	db.DB.Where("user_uid = ?", emailVerification.EmailUserID).First(&user)

	user.EmailConfirm = "true"

	db.DB.Save(&user)
	db.DB.Delete(&emailVerification)
        return c.SendString(user.Username + " you're email has been verified")
}


