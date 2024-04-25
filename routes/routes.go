// routes defines the HTTP routes for the WIG-Server application.
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
	app.Get("/ping", controller.Ping)

	// User Routes
	app.Post("/user/signup", controller.UserSignup)
	app.Get("/user/salt", controller.UserSalt)
	app.Post("/user/login", controller.UserLogin)
	app.Post("/app/validate", controller.UserValidate)

	// Scanner Routes
	app.Post("/app/scan/barcode", controller.ScannerBarcode)
	app.Get("/app/scan/check-qr", controller.ScannerCheckQR)
	app.Get("/app/scan/qr/location", controller.ScannerQRLocation)
	app.Get("/app/scan/qr/ownership", controller.ScannerQROwnership)

	// Ownership Routes
	app.Post("/app/ownership/create", controller.OwnershipCreate)
	app.Put("/app/ownership/quantity/:type", controller.OwnershipQuantity)
	app.Put("/app/ownership/edit", controller.OwnershipEdit)
	app.Put("/app/ownership/set-location", controller.OwnershipSetLocation)
	app.Delete("/app/ownership/delete", controller.OwnershipDelete)
	app.Post("/app/ownership/search", controller.OwnershipSearch)

	// Location Routes
	app.Post("/app/location/create", controller.LocationCreate)
	app.Put("/app/location/set-location", controller.LocationSetParent)
	app.Put("/app/location/edit", controller.LocationEdit)
	app.Post("/app/location/unpack", controller.LocationUnpack)
	app.Post("/app/location/search", controller.LocationSearch)

	// Borrower Routes
	app.Post("/app/borrower/create", controller.BorrowerCreate)
	app.Post("/app/borrower/checkout", controller.BorrowerCheckout)
	app.Post("/app/borrower/checkin", controller.BorrowerCheckin)
	app.Get("/app/borrower/get", controller.BorrowerGetAll)
	app.Get("/app/borrower/getcheckedout", controller.BorrowerGetInventory)

	app.Get("/app/inventory", controller.LocationGetInventory)
}
