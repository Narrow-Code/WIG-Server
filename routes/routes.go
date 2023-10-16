/*
* Package routes defines the HTTP routes for the WIG-Server application.
*/
package routes

import( 
	controller "WIG-Server/controller"
	"github.com/gofiber/fiber/v2"
)
/*
* Setup configures the routes on a Fiber application.
*/
func Setup(app *fiber.App){
	app.Post("/users/signup", controller.Signup)
}
