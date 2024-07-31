// routes defines the HTTP routes for the WIG-Server application.
package routes

import (
	controller "WIG-Server/controller"
	"WIG-Server/verification"

	"github.com/gofiber/fiber/v2"
)

/*
* Configures the routes on a Fiber application.
*
* @param app The Fiber application instance on which the routes will be configured.
 */
func Setup(app *fiber.App) {
	app.Get("/ping", controller.Ping)
	app.Get("/verification/:uid", controller.VerificationEmail)
	app.All("/resetpassword/:uid", verification.ResetPasswordPage)

	// User Routes
	app.Post("/user/signup", controller.UserSignup)
	app.Get("/user/:username/salt", controller.UserSalt)
	app.Post("/user/login", controller.UserLogin)
	app.Get("/app/validate", controller.UserValidate)
	app.Post("/user/verification", controller.ResendVerificationEmail)
	app.Post("/user/reset-password", controller.ResetPassword)

	// Scanner Routes
	app.Post("/app/scan/:barcode", controller.ScannerBarcode)
	app.Get("/app/scan/check", controller.ScannerCheckQR)
	app.Get("/app/scan/location", controller.ScannerQRLocation)
	app.Get("/app/scan/ownership", controller.ScannerQROwnership)

	// Ownership Routes
	app.Post("/app/ownership", controller.OwnershipCreate)
	app.Put("/app/ownership", controller.OwnershipEdit)
	app.Delete("/app/ownership", controller.OwnershipDelete)
	app.Put("/app/ownership/:ownershipUID/quantity/:type", controller.OwnershipQuantity) // TODO fix
	app.Put("/app/ownership/:ownershipUID/set-parent", controller.OwnershipSetLocation)
	app.Post("/app/ownership/search", controller.OwnershipSearch)

	// Location Routes
	app.Post("/app/location", controller.LocationCreate)
	app.Put("/app/location", controller.LocationEdit)
	app.Delete("/app/location", controller.LocationDelete)
	app.Put("/app/location/:locationUID/set-parent", controller.LocationSetParent)
	app.Get("/app/location/:locationUID", controller.LocationUnpack)
	app.Post("/app/location/search", controller.LocationSearch)

	// Borrower Routes
	app.Post("/app/borrower", controller.BorrowerCreate)
	app.Get("/app/borrower", controller.BorrowerGetAll)
	app.Delete("/app/borrower", controller.BorrowerDelete)
	app.Post("/app/borrower/:borrowerUID/checkout", controller.BorrowerCheckout)
	app.Post("/app/borrower/check-in", controller.BorrowerCheckin)
	app.Get("/app/borrower/checked-out", controller.BorrowerGetInventory)

	app.Get("/app/inventory", controller.LocationGetInventory)
}
