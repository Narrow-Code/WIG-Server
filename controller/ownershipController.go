/* Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
 */
package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

/*
IncrementOwnership increases the ownerships quantity by the designated value

@param c *fiber.Ctx
*/
func OwnershipQuantity(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")
	amountStr := c.Query("amount")
	changeType := c.Params("type")

	// Convert amount to int
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return Error(c, 400, messages.ConversionError)
	}
	if amount < 0 {
		return Error(c, 400, messages.NegativeError)
	}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Valide and retreive the ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check type of change
	switch changeType {
	case "increment":
		ownership.ItemQuantity += amount
	case "decrement":
		ownership.ItemQuantity -= amount
		if ownership.ItemQuantity < 0 {
			ownership.ItemQuantity = 0
		}
	case "set":
		ownership.ItemQuantity = amount
	default:
		return Error(c, 400, messages.InvalidChangeType)
	}

	// Save new amount to the database and create response
	db.DB.Save(&ownership)

	ownershipDTO := DTO("ownership", ownership)
	return Success(c, "Item found", ownershipDTO)
}

func OwnershipDelete(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	db.DB.Delete(&ownership)

	// Check for errors after the delete operation
	if result := db.DB.Delete(&ownership); result.Error != nil {
		return Error(c, 500, messages.ErrorDeletingOwnership)
	}

	// Ownership successfully deleted
	return Success(c, messages.OwnershipDelete)
}

func OwnershipEdit(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)

	if err != nil {
		return Error(c, code, err.Error())
	}

	// Add new fields
	ownership.CustomItemName = c.Query("custom_item_name")
	ownership.CustItemImg = c.Query("custom_item_img")
	ownership.OwnedCustDesc = c.Query("custom_item_description")
	ownership.ItemTags = c.Query("item_tags")

	db.DB.Save(&ownership)

	// Ownership successfully updated
	return Success(c, messages.OwnershipUpdated)
}

func OwnershipCreate(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return Error(c, code, err.Error())
	}

	// convert uid to uint
	itemUID, err := strconv.ParseUint(c.Query("item_uid"), 10, 64)
	if err != nil {
		return Error(c, 400, messages.ConversionError)
	}

	ownership, err := createOwnership(userUID, uint(itemUID))
	if err != nil {
		return Error(c, code, err.Error())
	}

	ownershipDTO := DTO("ownershipUID", ownership.OwnershipUID)
	return Success(c, messages.OwnershipCreated, ownershipDTO)
}

func OwnershipSetLocation(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	locationQR := c.Query("location_qr")
	ownershipUID := c.Query("ownershipUID")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err = recordExists("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var ownership models.Ownership
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	ownership.ItemLocation = location.LocationUID
	db.DB.Save(&ownership)

	// return success
	return Success(c, "Ownership set in "+location.LocationName)
}
