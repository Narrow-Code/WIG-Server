package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {
	return c.SendString("Signup route works")
}
