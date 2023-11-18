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
	// User Routes
	app.Post("/user/signup", controller.PostSignup)
	app.Get("/user/salt", controller.GetSalt)
	app.Post("/user/login", controller.PostLogin)
	app.Get("/user/validate", controller.PostLoginCheck)

	// Scanner Routes
	app.Post("/scan/barcode", controller.GetBarcode)
	app.Get("/scan/check-type", controller.CheckQR)
	
	// Ownership Routes
	app.Post("/ownership/create", controller.CreateOwnership)
	app.Put("/ownership/quantity/:type", controller.ChangeQuantity)
	app.Put("/ownership/edit", controller.EditOwnership)
	app.Put("/ownership/set-location", controller.SetOwnershipLocation)
	app.Delete("/ownership/delete", controller.DeleteOwnership)

	// Location Routes
	app.Post("/location/create/:type", controller.CreateLocation)
	app.Put("/location/set-location", controller.SetLocationLocation)
	app.Put("/location/edit", controller.EditLocation)
}
