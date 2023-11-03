// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/components"
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
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
func GetSalt(c *fiber.Ctx) error {
	// Get parameters
	username:= c.Query("username")
        if username == "" {return returnError(c, 400, messages.UsernameEmpty)}
	
	// Query database for username
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	code, err := recordExists("Username", result)
	if err != nil {return returnError(c, code, err.Error())}

	return c.Status(200).JSON(
                  fiber.Map{
                       "success":true,
                       "message":messages.SaltReturned,          
                       "salt":user.Salt})
}

/*
PostLoginCheck handles user login checks.
It checks if the user is logged in at initial start of application, making sure passwords have not changed.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occured during the registration procces.
*/
func PostLoginCheck(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
        if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Return success
	return returnSuccess(c, messages.TokenPass) 
}

/*
PostLogin handles user login requests.
If successful, it returns a JSON response with a success message and access token.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occurred during the registration process.
*/
func PostLogin(c *fiber.Ctx) error {
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}

	// Check for empty fields
	if data["username"] == "" {return returnError(c, 400, messages.UsernameEmpty)}
	if data["hash"] == "" {return returnError(c, 400, messages.HashMissing)}

	// Check that user exists
	var user models.User
        result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordExists("Username", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Check if hash matches then generate token
	if data["hash"] != user.Hash {return returnError(c, 400, messages.UsernamePasswordDoNotMatch)}
	token := components.GenerateToken(user.Username, user.Hash)
	
	// Retrun success
	return c.Status(200).JSON(
		fiber.Map{
			"success":true,
			"message":messages.UserLoginSuccess,
			"token":token,  
                        "uid":user.UserUID})
}

/*
PostSignup handles user registration requests.
It performs various checks such as data validation and database uniqueness before creating a new user record.
/If successful, it returns a JSON response with a success message and the user data.

@param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@return error - An error, if any, that occurred during the registration process.
*/
func PostSignup(c *fiber.Ctx) error {
	// Parse request into data map 
	var data map[string]string	
	err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}

	// Check for empty fields
	if data["username"] == "" {return returnError(c, 400, messages.UsernameEmpty)}
	if data["email"] == "" {return returnError(c, 400, messages.EmailEmpty)}
	if data["salt"] == "" {return returnError(c, 400, messages.SaltMissing)}
	if data["hash"] == "" {return returnError(c, 400, messages.HashMissing)}

	// Query for username in database
	var user models.User
	result := db.DB.Where("username = ?", data["username"]).First(&user)
	code, err := recordNotInUse("Username", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Query for email in database
	result = db.DB.Where("email = ?", data["email"]).First(&user)
	code, err = recordNotInUse("Email", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Check username requirements
	if !usernameRegex.MatchString(data["username"]){return returnError(c, 400, messages.ErrorUsernameRequirements)}
	if !emailRegex.MatchString(data["email"]) {return returnError(c, 400, messages.ErrorEmailRequirements)} 
	
	// Run DNS check on Email
	domain := strings.Split(data["email"], "@")[1]
	_, err = net.LookupMX(domain)
	if err != nil {return returnError(c, 400, messages.ErrorEmailRequirements)}

	// Set user model
	user = models.User{
		Username: data["username"],
		Email: data["email"],
		Salt: data["salt"],
		Hash: data["hash"],
		EmailConfirm: "false",
	}

	db.DB.Create(&user)

	// TODO send verification email

	return returnSuccess(c, messages.SignupSuccess) 
}
