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
	app.Post("/user/signup", controller.UserSignup)
	app.Get("/user/salt", controller.UserSalt)
	app.Post("/user/login", controller.UserLogin)
	app.Post("/user/validate", controller.UserValidate)

	// Scanner Routes
	app.Post("/scan/barcode", controller.ScanBarcode)
	app.Get("/scan/check-qr", controller.ScanCheckQR)
	
	// Ownership Routes
	app.Post("/ownership/create", controller.OwnershipCreate)
	app.Put("/ownership/quantity/:type", controller.OwnershipQuantity)
	app.Put("/ownership/edit", controller.OwnershipEdit)
	app.Put("/ownership/set-location", controller.OwnershipSetLocation)
	app.Delete("/ownership/delete", controller.OwnershipDelete)

	// Location Routes
	app.Post("/location/create/:type", controller.LocationCreate)
	app.Put("/location/set-location", controller.LocationSetLocation)
	app.Put("/location/edit", controller.LocationEdit)
}
