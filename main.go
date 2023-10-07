package main

import (
	routes "WIG-Server/routes"
	db "WIG-Server/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":30001")
}

