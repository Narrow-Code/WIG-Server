package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/upcitemdb"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
GetBarcode handles the functionality of returning any ownerships and items back after scanning a barcode.

@param c *fiber.Ctx
*/
func ScanBarcode(c *fiber.Ctx) error {
	// Initialize variables
	uid := c.Locals("uid").(string)
	barcode := c.Query("barcode")

	// Validate barcode
	if barcode == "" {
		return Error(c, 400, "Barcode is empty and required")
	}
	barcodeCheck, err := strconv.Atoi(barcode)
	if err != nil || barcodeCheck < 0 {
		return Error(c, 400, "There was an error converting barcode to an Int")
	}

	// Check if item exists in local database
	var item models.Item
	result := db.DB.Where("barcode = ?", barcode).First(&item)

	// If item isn't found, check api and add to
	if result.Error == gorm.ErrRecordNotFound {
		upcitemdb.GetBarcode(barcode)
		result = db.DB.Where("barcode = ?", barcode).First(&item)
		if result.Error == gorm.ErrRecordNotFound {
			return Error(c, 400, "Item was not found in the database")
		}
	}

	// If there is a connection error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Search Ownership by barcode
	var ownerships []models.Ownership
	db.DB.Where("item_number = ? AND item_owner = ?", item.ItemUid, uid).Find(&ownerships)

	// If no ownership exists, create ownership
	if len(ownerships) == 0 {
		ownership, err := createOwnership(uid, item.ItemUid)
		if err != nil {
			return Error(c, 400, err.Error())
		}
		ownership.ItemQuantity = 1
		ownerships = append(ownerships, ownership)
	}

	itemDTO := DTO("item", item.Name)
	ownershipDTO := DTO("ownership", ownerships)

	return Success(c, "Item found", itemDTO, ownershipDTO)
}

/*
CheckQR takes a QR code as parameter, and checks whether it is an item, location or a unused QR.

@param c *fiber.Ctx - The fier context containing the HTTP request and response objects.
@return error - An error that occured during the process or if the token does not match
*/
func ScanCheckQR(c *fiber.Ctx) error {
	// Initialize variables
	uid := c.Locals("uid")
	qr := c.Query("qr")

	// Check for empty fields
	if qr == "" {
		return Error(c, 400, "QR is empty and required")
	}

	// Check if qr exists as location
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", qr, uid).First(&location)
	if location.LocationUID != 0 {
		return Success(c, "LOCATION")
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Check if qr exists as ownership
	var ownership models.Ownership
	result = db.DB.Where("item_qr = ? AND item_owner = ?", qr, uid).First(&ownership)
	if ownership.OwnershipUID != 0 {
		return Success(c, "OWNERSHIP")
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	return Success(c, "NEW")
}
