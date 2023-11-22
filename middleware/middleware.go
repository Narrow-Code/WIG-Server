/*
* The middleware package handles all middle functions between the API calls.
 */
package middleware

import (
	"WIG-Server/messages"
	"WIG-Server/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
* AppAuthHeaderCheck checks that the AppAuth header is valid.
 */
func AppAuthHeaderCheck() fiber.Handler {
	// Get AppAuth secret
	godotenv.Load()
        
	return func(c *fiber.Ctx) error{
		headerValue := c.Get("AppAuth")
		
		if headerValue != os.Getenv("APP_SECRET"){
			return utils.NewError(c, 400, messages.AccessDenied)
		}
		
		return c.Next()
	}

}
