package utils

import "github.com/gofiber/fiber/v2"

/*
returnSuccess returns a 200 success code, a 'true' success status and a message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param message The success message to return via fiber.
@return error - An error, if any, that occurred during the process.
*/
func NewSuccess (c *fiber.Ctx, message string, fields ...ResponseField) error {
	responseMap := fiber.Map{
		"success":true,
		"message":message,
	}

	for _, field := range fields {
		responseMap[field.Field] = field.Response
	}

	return c.Status(200).JSON(responseMap)
}

