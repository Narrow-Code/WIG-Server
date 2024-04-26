// middleware handles all middle functions between the API calls.
package middleware

import (
	"WIG-Server/controller"
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
* AppAuth checks that the AppAuth header is valid.
 */
func AppAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Load environment variables
		godotenv.Load()

		// Get header value
		headerValue := c.Get("AppAuth")

		// Check if header value matches app secret
		if headerValue != os.Getenv("APP_SECRET") {
			return controller.Error(c, fiber.StatusBadRequest, "Unauthorized")
		}

		// Continue to the next middleware
		utils.Log("authorized")
		return c.Next()
	}
}

/*
* ValidateToken checks that the users token is still valid.
 */
func ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		utils.Log("validating token")
		token := c.Get("Authorization")

		// Check if token is missing
		if token == "" {
			return controller.Error(c, fiber.StatusBadRequest, "Token missing")
		}

		// Query database for user with the given token
		var user models.User
		result := db.DB.Where("token = ?", token).First(&user)

		// Check if user with token exists
		if result.Error != nil {
			return controller.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		// Store user in context locals
		c.Locals("user", user)

		// Continue to the next middleware
		utils.Log(user.Username + " authorized")
		return c.Next()
	}
}
