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

// The regex expression to check username requirements
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)

// The regex expression to check email requirements
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

/*
GetSalt handles user login requesting salt.
It checks for existing username and returns the salt.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occurred during the registration process.
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
	code, err := RecordExists("Username", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	saltDTO := DTO("salt", user.Salt)

	return Success(c, "Salt returned successfully", saltDTO)
}

/*
PostLoginCheck handles user login checks.
It checks if the user is logged in at initial start of application, making sure passwords have not changed.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occured during the registration procces.
*/
func UserValidate(c *fiber.Ctx) error {
	return Success(c, "Authorized")
}

/*
PostLogin handles user login requests.
If successful, it returns a JSON response with a success message and access token.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occurred during the registration process.
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
	code, err := RecordExists("Username", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check if hash matches then generate token
	if data["hash"] != user.Hash {
		return Error(c, 400, "The username and passwords do not match")
	}

	user.Token = utils.GenerateToken(user.Username, user.Hash)

	tokenDTO := DTO("token", user.Token)
	uidDTO := DTO("uid", user.UserUID)

	db.DB.Save(&user)

	return Success(c, "Login was successful", tokenDTO, uidDTO)
}

/*
PostSignup handles user registration requests.
It performs various checks such as data validation and database uniqueness before creating a new user record.
/If successful, it returns a JSON response with a success message and the user data.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occurred during the registration process.
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
	code, err := recordNotInUse("Username", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Query for email in database
	result = db.DB.Where("email = ?", data["email"]).First(&user)
	code, err = recordNotInUse("Email", result)
	if err != nil {
		return Error(c, code, err.Error())
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
	}

	db.DB.Create(&user)

	// TODO send verification email

	return Success(c, "Signup was successful")
}
