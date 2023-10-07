package routes

import( 
	controller "WIG-Server/controller"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App){
	app.Get("/users/signup", controller.Signup)
}
