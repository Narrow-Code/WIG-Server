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
	app.Post("/users/signup", controller.PostSignup)
	app.Get("/users/salt", controller.GetSalt)
	app.Post("/users/login", controller.PostLogin)
	app.Post("/users/login/check", controller.PostLoginCheck)
	app.Post("/items/barcode", controller.GetBarcode)

	app.Post("/code/check", controller.CheckQR)

	app.Put("/ownership/quantity/:type", controller.ChangeQuantity)
	app.Delete("/ownership/delete", controller.DeleteOwnership)
	app.Put("/ownership/edit", controller.EditOwnership)
	app.Put("/ownership/location/set", controller.SetOwnershipLocation)
	app.Post("/ownership/create", controller.CreateOwnership)

	app.Post("/location/create/:type", controller.CreateLocation)
	app.Put("/location/location/set", controller.SetLocationLocation)
}
