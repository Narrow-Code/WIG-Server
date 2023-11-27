// The Handles all middle functions between the API calls.
package middleware

import (
	"WIG-Server/controller"
	"WIG-Server/db"
	"WIG-Server/models"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
* Checks that the AppAuth header is valid.
 */
func AppAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		godotenv.Load()
		headerValue := c.Get("AppAuth")

		if headerValue != os.Getenv("APP_SECRET") {
			return controller.Error(c, 400, "Unauthorized")
		}
		return c.Next()
	}
}

/*
* Checks that the users token is still valid.
 */
func ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		fmt.Println("WORKING")

		if token == "" {
			return controller.Error(c, fiber.StatusBadRequest, "Token missing")
		}

		var user models.User
		result := db.DB.Where("token = ?", token).First(&user)

		if result.Error != nil {
			return controller.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		c.Locals("uid", user)
		return c.Next()
	}
}
