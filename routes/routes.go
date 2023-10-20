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
*
* @param app *fiber.App - The Fiber application instance on which the routes will be configured.
*/
func Setup(app *fiber.App){
	app.Post("/users/signup", controller.Signup)
	app.Get("/users/salt", controller.GetSalt)
	app.Post("/users/login", controller.PostLogin)
	app.Get("/users/login", controller.GetLogin)
}
