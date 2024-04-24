package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"net"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// The regex expression to check username requirements
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)

// The regex expression to check email requirements
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

/*
* Retrieves the users salt. 
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func UserSalt(c *fiber.Ctx) error {
	// Get parameters
	username := c.Query("username")
	if username == "" {
		return Error(c, 400, "Username is empty and required")
	}

	// Query database for username
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	saltDTO := DTO("salt", user.Salt)

	return Success(c, "Salt returned successfully", saltDTO)
}

/*
* Validates the users token is still valid. 
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func UserValidate(c *fiber.Ctx) error {
	return Success(c, "Authorized")
}

/*
* Handles the user Login logic. Returning a token. 
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func UserLogin(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Check for empty fields
	if data["username"] == "" || data["hash"] == "" {
		return Error(c, 400, "Username or hash is empty and required")
	}

	// Check that user exists
	var user models.User
	result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check if hash matches then generate token
	if data["hash"] != user.Hash {
		return Error(c, 400, "The username and passwords do not match")
	}

	user.Token = utils.GenerateToken(user.Username, user.Hash)

	if user.Token == "error" {
		return Error(c, 400, "There was an error generating user token")
	}

	tokenDTO := DTO("token", user.Token)
	uidDTO := DTO("uid", user.UserUID)

	db.DB.Save(&user)

	return Success(c, "Login was successful", tokenDTO, uidDTO)
}

/* 
* Handles user registration requests.
* It performs various checks such as data validation and database uniqueness before creating a new user record.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func UserSignup(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}

	// Check for empty fields
	if data["username"] == "" || data["email"] == "" {
		return Error(c, 400, "Username or email is empty and required")
	}
	if data["salt"] == "" || data["hash"] == "" {
		return Error(c, 400, "Sale or hash is empty and required")
	}

	// Query for username in database
	var user models.User
	result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Username: " + err.Error())
	}

	// Query for email in database
	result = db.DB.Where("email = ?", data["email"]).First(&user)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, "Email: " + err.Error())
	}

	// Check username requirements
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

	// Set user model
	user = models.User{
		Username: data["username"],
		Email:    data["email"],
		Salt:     data["salt"],
		Hash:     data["hash"],
		UserUID:  uuid.New(),
	}

	db.DB.Create(&user)

	// TODO send verification email

	return Success(c, "Signup was successful")
}

func Ping (c *fiber.Ctx) error {
	return Success(c, "Ping was successful")
}
