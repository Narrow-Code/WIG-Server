// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
        "github.com/gofiber/fiber/v2"
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
func validateToken(c *fiber.Ctx, uid string, token string){

}
