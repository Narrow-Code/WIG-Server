package middleware

import (
	"github.com/joho/godotenv"
        "os"
	"github.com/gofiber/fiber/v2"
	"WIG-Server/messages"
)

func AppAuthHeaderCheck() fiber.Handler {
	// Get AppAuth secret
	godotenv.Load()
        
	return func(c *fiber.Ctx) error{
		headerValue := c.Get("AppAuth")
		
		if headerValue != os.Getenv("APP_SECRET"){
			return c.Status(400).JSON(
				fiber.Map{
					"success":false,
					"message":messages.AccessDenied})
		}
		
		return c.Next()
	}

}
