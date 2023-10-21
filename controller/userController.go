/* Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
*/

package controller

import (
	components "WIG-Server/components"
	"github.com/gofiber/fiber/v2"
	models "WIG-Server/models"
	db "WIG-Server/config"
	messages "WIG-Server/messages"
	"gorm.io/gorm"
	"regexp"
	"net"
	"strings"
)

/*
* The regex expression to check username requirements
*/
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,20}$`)

/*
* The regex expression to check email requirements
*/
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

/*
* GetSalt handles user login requesting salt.
* It checks for existing username and returns the salt.
*
* @param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
*
* @return error - An error, if any, that occurred during the registration process.
*/
func GetSalt(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
        err := c.BodyParser(&data)

        // Error with JSON request
        if err != nil {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message":messages.ErrorParsingRequest,
				"salt":""})
                        }
	// Query for username
	var user models.User
	result := db.DB.Where("username = ?", data["username"]).First(&user)


	// Check if user is found
	if result.Error == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(
			fiber.Map{
				"success": false,
				"message": messages.UsernameDoesNotExist,
				"salt":""})
	} else if result.Error != nil {
		return c.Status(400).JSON(
                       	fiber.Map{
                               	"success":false,
                               	"message":messages.ErrorWithConnection,
                               	"salt":""}) 
	} else {
		return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":messages.SaltReturned,          
                                "salt":user.UserSalt})
	}
}

/*
* GetLogin handles user login checks.
* It checks if the user is logged in at initial start of application, making sure passwords have not changed.
*
* @param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
*
* @return error - An error, if any, that occured during the registration procces.
*/
func GetLogin(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)

	// Error with JSON request
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.ErrorParsingRequest})
			}

	// Query for UID
	var user models.User

	// Check if UID and token exist
        if data["uid"] == "" {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message": messages.UIDEmpty})
        }

        if data["token"] == "" {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message":messages.TokenEmpty})
        }


        result := db.DB.Where("user_uid = ?", data["uid"]).First(&user)

	// Check if UID was found
	if result.Error == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(
                        fiber.Map{
                                "success": false,
                                "message": messages.RecordNotFound})
        } else if result.Error != nil {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message":messages.ErrorWithConnection}) 
        }
	
	// Validate token
	if !components.ValidateToken(user.Username, user.UserHash, data["token"]) {
		return c.Status(400).JSON(
                fiber.Map{
                        "success":false,
                        "message":messages.ErrorToken})
		}
	
	// Token matches
	return c.Status(200).JSON(
		fiber.Map{
			"success":true,
			"message":messages.TokenPass})
}

/*
* Login handles user login requests.
* If successful, it returns a JSON response with a success message and access token.
*
* @param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
*
* @return error - An error, if any, that occurred during the registration process.
*/
func PostLogin(c *fiber.Ctx) error {
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)

	// Error with JSON request
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.ErrorParsingRequest,
				"token":"",
				"uid":""})
	}

	// Check if username and hash match
        var user models.User

	if data["username"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message": messages.UsernameEmpty,
				"token":"",
				"uid":""})
	}

	if data["hash"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.HashMissing,
				"token":"",
				"uid":""})
	}

        result := db.DB.Where("username = ?", data["username"]).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
                return c.Status(404).JSON(
                        fiber.Map{
                                "success": false,
                                "message": messages.UsernameDoesNotExist,
                                "token": "",
                                "uid": "",})
        } else if result.Error != nil {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message":messages.ErrorWithConnection,
                                "token":"",  
                                "uid":""}) 
        } else {
		if data["hash"] != user.UserHash {
			return c.Status(400).JSON(
				fiber.Map{
					"success":false,
					"message":messages.UsernamePasswordDoNotMatch,
					"token":"",  
                                	"uid":""})
		}
	}
	// Generate token
	token := components.GenerateToken(user.Username, user.UserHash)
	
	return c.Status(200).JSON(
		fiber.Map{
			"success":true,
			"message":messages.UserLoginSuccess,
			"token":token,  
                        "uid":user.UserUID})
}

/*
* Signup handles user registration requests.
* It performs various checks such as data validation and database uniqueness before creating a new user record.
* If successful, it returns a JSON response with a success message and the user data.
*
* @param c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
*
* @return error - An error, if any, that occurred during the registration process.
*/
func Signup(c *fiber.Ctx) error {
	// Parse request into data map 
	var data map[string]string	
	err := c.BodyParser(&data)

	// Error with JSON request
	if err != nil {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.ErrorParsingRequest})
			}

	// Username empty error
	if data["username"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.UsernameEmpty})
	}

	// Email empty error
	if data["email"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.EmailEmpty})
		}
	
	// Salt empty error
	if data["salt"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.SaltMissing})
		}

	// Hash empty error
	if data["hash"] == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.HashMissing})
		}

	// Query for username in database
	var user models.User
	userResult := db.DB.Where("username = ?", data["username"]).First(&user)
 
	if userResult.RowsAffected != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.UsernameInUse})
		}

	// Query for email
	emailResult := db.DB.Where("user_email = ?", data["email"]).First(&user)

	// Error with connection
	if emailResult.RowsAffected != 0 {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.EmailInUse})
		}

	// Check username requirements
	if !usernameRegex.MatchString(data["username"]) {
		return c.Status(400).JSON(
			fiber.Map{
				"success":false,
				"message":messages.ErrorUsernameRequirements})
	}
	
	// TODO Check email validity
	if !emailRegex.MatchString(data["email"]) {
                return c.Status(400).JSON(
                        fiber.Map{
                                "success":false,
                                "message":messages.ErrorEmailRequirements})
	} else {
		// Run DNS check on Email
		domain := strings.Split(data["email"], "@")[1]
		_, err := net.LookupMX(domain)
		
		if err != nil {
			return c.Status(400).JSON(
				fiber.Map{
					"success":false,
					"message":messages.ErrorEmailRequirements})
	}



	}
	// Set up fields
	user = models.User{
		Username: data["username"],
		UserEmail: data["email"],
		UserSalt: data["salt"],
		UserHash: data["hash"],
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
