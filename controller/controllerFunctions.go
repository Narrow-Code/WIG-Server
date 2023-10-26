package controller

import (
        "github.com/gofiber/fiber/v2"
)

func returnError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"success":false,
		"message":message})
}
