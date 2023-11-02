// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/components"
	"github.com/gofiber/fiber/v2"
	"WIG-Server/models"
	"WIG-Server/db"

	"WIG-Server/messages"
	"gorm.io/gorm"
	"regexp"
	"net"
	"strings"
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
	
	// Check if parameters are empty
        if username == "" {
		return returnError(c, 400, messages.UsernameEmpty)
	}
	
	// Query for username
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	
	// Check if user is found
	if result.Error == gorm.ErrRecordNotFound {
		return returnError(c, 404, messages.UsernameDoesNotExist)

	} else if result.Error != nil {
		return returnError(c, 400, messages.ErrorWithConnection)

	} else {
		return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":messages.SaltReturned,          
                                "salt":user.Salt})
	}
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

        // Error with JSON request
        if err != nil {
                return returnError(c, 400, messages.ErrorParsingRequest)
        }

	// Validate Token
	err = validateToken(c, data["uid"], data["token"])	
	if err == nil {
		return validateToken(c, data["uid"], data["token"])
	}

	// Token matches
	return c.Status(200).JSON(
		fiber.Map{
			"success":true,
			"message":messages.TokenPass})
		
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

	// Error with JSON request
	if err != nil {
		return returnError(c, 400, messages.ErrorParsingRequest)
	}

	// Check if username and hash match
        var user models.User

	if data["username"] == "" {
		return returnError(c, 400, messages.UsernameEmpty)
	}

	if data["hash"] == "" {
		return returnError(c, 400, messages.HashMissing)
	}

        result := db.DB.Where("username = ?", data["username"]).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
                return returnError(c, 404, messages.UsernameDoesNotExist)
        } else if result.Error != nil {
                return returnError(c, 400, messages.ErrorWithConnection)
        } else {
		if data["hash"] != user.Hash {
			return returnError(c, 400, messages.UsernamePasswordDoNotMatch)
		}
	}

	// Generate token
	token := components.GenerateToken(user.Username, user.Hash)
	
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

	// Error with JSON request
	if err != nil {
		return returnError(c, 400, messages.ErrorParsingRequest)
			}

	// Username empty error
	if data["username"] == "" {
		return returnError(c, 400, messages.UsernameEmpty) 
	}

	// Email empty error
	if data["email"] == "" {
		return returnError(c, 400, messages.EmailEmpty)
		}
	
	// Salt empty error
	if data["salt"] == "" {
		return returnError(c, 400, messages.SaltMissing)
		}

	// Hash empty error
	if data["hash"] == "" {
		return returnError(c, 400, messages.HashMissing)
		}

	// Query for username in database
	var user models.User
	userResult := db.DB.Where("username = ?", data["username"]).First(&user)
 
	if userResult.RowsAffected != 0 {
		return returnError(c, 400, messages.UsernameInUse)
		}

	// Query for email
	emailResult := db.DB.Where("email = ?", data["email"]).First(&user)

	// Error with connection
	if emailResult.RowsAffected != 0 {
		return returnError(c, 400, messages.EmailInUse)
		}

	// Check username requirements
	if !usernameRegex.MatchString(data["username"]) {
		return returnError(c, 400, messages.ErrorUsernameRequirements)
	}
	
	// Check email validity
	if !emailRegex.MatchString(data["email"]) {
                return returnError(c, 400, messages.ErrorEmailRequirements)
	} else {
		// Run DNS check on Email
		domain := strings.Split(data["email"], "@")[1]
		_, err := net.LookupMX(domain)
		
		if err != nil {
			return returnError(c, 400, messages.ErrorEmailRequirements)
	}

	}
	// Set up fields
	user = models.User{
		Username: data["username"],
		Email: data["email"],
		Salt: data["salt"],
		Hash: data["hash"],
		EmailConfirm: "false",
	}

	// TODO get error for failure?

	db.DB.Create(&user)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message":messages.SignupSuccess,
	})

	// TODO Send verification email
}
