/*
* Package main is the entry point for the WIG-Server application.
*
* This package initializes the database connection, sets up routes, and starts the server.
*/

package main

import (
	routes "WIG-Server/routes"
	db "WIG-Server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
        "os"
)

/*
* main is the entry point of the application.
*
* It connects to the database, sets up routes, and starts the server.
*/
func main() {
	// Get port number
	godotenv.Load()
        var port = os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// Connect to the database
	db.Connect()

	// Create a new Fiber app
	app := fiber.New()

	// Setup routes
	routes.Setup(app)

	// Start the server and lsiten on port 30001
	app.Listen(":" + port)
}

