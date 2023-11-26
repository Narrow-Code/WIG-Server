// Defines the HTTP routes for the WIG-Server application.
package routes

import (
	controller "WIG-Server/controller"

	"github.com/gofiber/fiber/v2"
)

/*
* Configures the routes on a Fiber application.
*
* @param app The Fiber application instance on which the routes will be configured.
 */
func Setup(app *fiber.App) {
	// User Routes
	app.Post("/user/signup", controller.UserSignup)
	app.Get("/user/salt", controller.UserSalt)
	app.Post("/user/login", controller.UserLogin)
	app.Post("/app/validate", controller.UserValidate)

	// Scanner Routes
	app.Post("/app/scan/barcode", controller.ScanBarcode)
	app.Get("/app/scan/check-qr", controller.ScanCheckQR)

	// Ownership Routes
	app.Post("/app/ownership/create", controller.OwnershipCreate)
	app.Put("/app/ownership/quantity/:type", controller.OwnershipQuantity)
	app.Put("/app/ownership/edit", controller.OwnershipEdit)
	app.Put("/app/ownership/set-location", controller.OwnershipSetLocation)
	app.Delete("/app/ownership/delete", controller.OwnershipDelete)

	// Location Routes
	app.Post("/app/location/create/:type", controller.LocationCreate)
	app.Put("/app/location/set-location", controller.LocationSetLocation)
	app.Put("/app/location/edit", controller.LocationEdit)
}
