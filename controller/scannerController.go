package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/upcitemdb"
	"WIG-Server/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
* Takes a barcode and searches to see if an item in the database exists with the barcode.
* If an item does not exist, it makes API call to upcitemdb.com to search barcode.
* If the item exists at upcitemdb, it creates a new item with that data.
* Then after all, it creates an ownership with item and userdata.
 */
func ScannerBarcode(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")	
	var item models.Item
	var ownerships []models.Ownership
	user := c.Locals("user").(models.User)
	barcode := c.Params("barcode")

	// Check if barcode is empty and convert to int
	utils.UserLog(c, "checking for empty fields")
	if barcode == "" {
		return Error(c, 400, "Barcode is empty and required")
	}
	utils.UserLog(c, "converting barcode to int")
	barcodeCheck, err := strconv.Atoi(barcode)
	if err != nil || barcodeCheck < 0 {
		return Error(c, 400, "There was an error converting barcode to an Int")
	}

	// Check if item exists
	utils.UserLog(c, "checking if item exists")
	result := db.DB.Where("barcode = ?", barcode).First(&item)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// If item doesn't exist call upcitemdb and add to database
	if result.Error == gorm.ErrRecordNotFound {
		// Check if API limit has been reached
		utils.UserLog(c, "call upcitemdb")
		code := upcitemdb.GetBarcode(barcode)
		if code == 429 {
			return Error(c, code, "API limit reached")
		}

		// Check if item was added to database
		utils.UserLog(c, "checking if item was added to database")
		result = db.DB.Where("barcode = ?", barcode).First(&item)
		if result.Error == gorm.ErrRecordNotFound {
			return Error(c, 404, "Item was not found in the database")
		}
	}

	// Search Ownership by uid
	utils.UserLog(c, "checking for existing ownership")
	db.DB.Where("item_number = ? AND item_owner = ?", item.ItemUid, user.UserUID).Find(&ownerships)

	// If no ownership exists, create ownership
	if len(ownerships) == 0 {
		utils.UserLog(c, "creating ownership")
		ownership, err := createOwnership(user.UserUID, item, "", "")
		if err != nil {
			return Error(c, 400, err.Error())
		}
		ownerships = append(ownerships, ownership)
	}

	// Preload ownerships, add to dto and return
	for i := range ownerships {
		preloadOwnership(&ownerships[i])
	}
	dto := DTO("ownership", ownerships)
	utils.UserLog(c, "success")
	return success(c, "Item found", dto)
}

// ScannerCheckQR takes a QR code as parameters, and checks whether it is an item, location or an unused QR.
func ScannerCheckQR(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
	qr := c.Params("qr")

	// Check for empty fields
	utils.UserLog(c, "checking for empty fields")
	if qr == "" {
		return Error(c, 400, "QR is empty and required")
	}

	// Check if qr exists as location
	utils.UserLog(c, "checking if qr exists as location")
	result := db.DB.Where("location_qr = ? AND location_owner = ?", qr, user.UserUID).First(&location)
	emptyUID := [16]byte{}
	if location.LocationUID != emptyUID {
		utils.UserLog(c, "returning location")
		return success(c, "LOCATION")
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Check if qr exists as ownership
	utils.UserLog(c, "checking if qr exists as ownership")
	result = db.DB.Where("item_qr = ? AND item_owner = ?", qr, user.UserUID).First(&ownership)
	if ownership.OwnershipUID != emptyUID {
		utils.UserLog(c, "returning ownership")
		return success(c, "OWNERSHIP")
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Return as unused QR code
	utils.UserLog(c, "returning new")
	return success(c, "NEW")
}

// ScannerQRLocation takes a QR code and returns its corresponding location
func ScannerQRLocation(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	user := c.Locals("user").(models.User)
	qr := c.Params("qr")

	// Check if QR is empty
	utils.UserLog(c, "checking for empty fields")
	if qr == "" {
		return Error(c, 400, "QR is empty and required")
	}

	// Check if location exists in local database
	utils.UserLog(c, "checking if location exists")
	result := db.DB.Where("location_qr = ? AND location_owner = ?", qr, user.UserUID).First(&location)
	if result.Error == gorm.ErrRecordNotFound {
		return Error(c, 400, "Item was not found in the database")
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Preload location, add to dto and return
	preloadLocation(&location)
	dto := DTO("location", location)
	utils.UserLog(c, "success")
	return success(c, "Location returned", dto)
}

// ScannerQROwnership takes a QR code and returns its corresponding ownership
func ScannerQROwnership(c *fiber.Ctx) error {
	// Initialize variables
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
	qr := c.Params("qr")

	// Check if QR is empty
	if qr == "" {
		return Error(c, 400, "QR is empty and required")
	}

	// Check if item exists in local database
	result := db.DB.Where("item_qr = ? AND item_owner = ?", qr, user.UserUID).First(&ownership)
	if result.Error == gorm.ErrRecordNotFound {
		return Error(c, 400, "Item was not found in the database")
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return Error(c, 400, "internal server error")
	}

	// Preload location, add to dto and return
	preloadOwnership(&ownership)
	dto := DTO("ownership", ownership)
	return success(c, "Ownership returned", dto)
}
