// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/components"
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"

	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
returnError returns the given error code, success status and message through fiber to the application.

@param c c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@param code The error code to return via fiber
@param The error message to retrun via fiber

@return error - An error, if any, that occurred during the registration process.
*/
func returnError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"success":false,
		"message":message})
}

/*
validateToken checks if a users UID and token match and are valid.

@param c *fiber.Ctx - The fier context containing the HTTP request and response objects.
@param UID The users UID
@param token The users authentication token

@return error - An error that occured during the process or if the token does not match
*/
func validateToken(c *fiber.Ctx, uid string, token string) error{
// Check if UID and token exist
        if uid == "" {
                return returnError(c, 400, messages.UIDEmpty)
        }

        if token == "" {
                return returnError(c, 400, messages.TokenEmpty)
        }

        // Query for UID
        var user models.User
        result := db.DB.Where("user_uid = ?", uid).First(&user)

        // Check if UID was found
        if result.Error == gorm.ErrRecordNotFound {
                return returnError(c, 404, messages.RecordNotFound)

        } else if result.Error != nil {
                return returnError(c, 400, messages.ErrorWithConnection)
	}

        // Validate token
        if !components.ValidateToken(user.Username, user.Hash, token) {
                return returnError(c, 400, messages.ErrorToken)
                }
	
	return errors.New(messages.TokenPass)
}
