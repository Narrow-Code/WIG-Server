/*
* The middleware package handles all middle functions between the API calls.
 */
package middleware

import (
	"WIG-Server/controller"
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
* AppAuthHeaderCheck checks that the AppAuth header is valid.
 */
func AppAuth() fiber.Handler {
	// Get AppAuth secret
	godotenv.Load()

	return func(c *fiber.Ctx) error {
		headerValue := c.Get("AppAuth")

		if headerValue != os.Getenv("APP_SECRET") {
			return controller.Error(c, 400, messages.AccessDenied)
		}

		return c.Next()
	}
}

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
			return controller.Error(c, fiber.StatusUnauthorized, messages.AccessDenied)
		}

		c.Locals("uid", strconv.FormatUint(uint64(user.UserUID), 10))
		return c.Next()
	}
}
