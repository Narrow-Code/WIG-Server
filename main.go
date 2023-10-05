package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create instance of fiber
	app := fiber.New()

	// Create httphandler
	app.Get("/username", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{
			"success":true,
			"message": "Get username created",
		})
	})

	// Listen on port
	app.Listen(":3000")

}
