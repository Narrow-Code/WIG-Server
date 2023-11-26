// The entry point for the WIG-Server application.
package main

import (
	"WIG-Server/routes"
	"WIG-Server/db"
	"github.com/gofiber/fiber/v2"
	"WIG-Server/middleware"
)

/*
* Connects to the database, sets up routes, and starts the backend server.
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

