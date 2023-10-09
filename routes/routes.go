package routes

import( 
	controller "WIG-Server/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Post("/users/signup", controller.Signup)
}
