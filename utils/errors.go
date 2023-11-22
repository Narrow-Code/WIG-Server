package utils

import "github.com/gofiber/fiber/v2"

/*
returnError returns the given error code, a 'false' success status and message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param code The error code to return via fiber
@param message The error message to return via fiber
@return error - An error, if any, that occurred during the process.
*/
func NewError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"success":false,
		"message":message})
}
