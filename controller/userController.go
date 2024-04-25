package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"net"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// usernameRegex is the regex expression to check username requirements
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)

// emailRegex is the regex expression to check email requirements
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// UserSalt retrieves the users salt.
func UserSalt(c *fiber.Ctx) error {
	// Initialize variables
	var user models.User
	username := c.Query("username")

	// Check if username is empty
	if username == "" {
		return Error(c, 400, "Username is empty and required")
	}

	// Query database for username
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}

	// Add to dto and return
	dto := DTO("salt", user.Salt)
	return success(c, "Salt returned successfully", dto)
}

// UserValidate validates the users token is still valid.
func UserValidate(c *fiber.Ctx) error {
	return success(c, "Authorized")
}

// UserLogin handles the user Login logic. Returning a token.
func UserLogin(c *fiber.Ctx) error {
	// Initialize variables
	var data map[string]string
	var user models.User
	username := data["username"]

	// Parse JSON body
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Check for empty fields
	if data["username"] == "" || data["hash"] == "" {
		return Error(c, 400, "Username or hash is empty and required")
	}

	// Check that user exists
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, "Username " + err.Error())
	}

	// Check if hash matches
	if data["hash"] != user.Hash {
		return Error(c, 400, "The username and passwords do not match")
	}

	// Generate token
	user.Token = utils.GenerateToken(user.Username, user.Hash)
	if user.Token == "error" {
		return Error(c, 400, "There was an error generating user token")
	}

	// Save to database, add to dto's and return
	db.DB.Save(&user)
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
	var data map[string]string
	var user models.User

	// Parse request into data map
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}

	// Check for empty fields
	if data["username"] == "" || data["email"] == "" {
		return Error(c, 400, "Username or email is empty and required")
	}
	if data["salt"] == "" || data["hash"] == "" {
		return Error(c, 400, "Salt or hash is empty and required")
	}

	// Check if username is in use
	result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Username: "+err.Error())
	}

	// Check if email is in use
	result = db.DB.Where("email = ?", data["email"]).First(&user)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Email: "+err.Error())
	}

	// Check username and email requierments
	if !usernameRegex.MatchString(data["username"]) {
		return Error(c, 400, "Username does not match requirements")
	}
	if !emailRegex.MatchString(data["email"]) {
		return Error(c, 400, "Email does not match requirements")
	}

	// Run DNS check on Email
	domain := strings.Split(data["email"], "@")[1]
	_, err = net.LookupMX(domain)
	if err != nil {
		return Error(c, 400, "Email domain does not exist")
	}

	// TODO send verification email

	// Create user and return
	user = createUser(data)
	return success(c, "Signup was successful")
}

// Ping performs a health check
func Ping(c *fiber.Ctx) error {
	return success(c, "Ping was successful")
}
