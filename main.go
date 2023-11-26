/*
* Package main is the entry point for the WIG-Server application.
*
* This package initializes the database connection, sets up routes, and starts the server.
*/
package main

import (
	"WIG-Server/routes"
	"WIG-Server/db"
	"github.com/gofiber/fiber/v2"
	"WIG-Server/middleware"
)

/*
* main is the entry point of the application.
*
* It connects to the database, sets up routes, and starts the server.
*/
func main() {
	db.Connect()
	app := fiber.New()
	app.Use(middleware.AppAuth())
	loggedRoutes := app.Group("/app")
	loggedRoutes.Use(middleware.ValidateToken())
	routes.Setup(app)
	app.Listen(":" + db.GetPort()) 
}

