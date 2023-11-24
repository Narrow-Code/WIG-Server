package utils

import (
	"WIG-Server/dto"

	"github.com/gofiber/fiber/v2"
)

/*
returnSuccess returns a 200 success code, a 'true' success status and a message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param message The success message to return via fiber.
@return error - An error, if any, that occurred during the process.
*/
func Success(c *fiber.Ctx, message string, dtos ...dto.DTO) error {
	responseMap := fiber.Map{
		"message": message,
	}

	for _, dto := range dtos {
		responseMap[dto.Name] = dto.Data
	}

	return c.Status(200).JSON(responseMap)
}

/*
returnError returns the given error code, a 'false' success status and message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param code The error code to return via fiber
@param message The error message to return via fiber
@return error - An error, if any, that occurred during the process.
*/
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message})
}
