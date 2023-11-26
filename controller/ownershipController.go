package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

/*
* Changes the quantity of an ownership, using increment, decrement or setter method.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func OwnershipQuantity(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	ownershipUID := c.Query("ownershipUID")
	amountStr := c.Query("amount")
	changeType := c.Params("type")

	// Convert amount to int
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return Error(c, 400, "There was an error converting amount to Int")
	}
	if amount < 0 {
		return Error(c, 400, "Amount cannot be negative")
	}

	// Valide and retreive the ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err := RecordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check type of change
	switch changeType {
	case "increment":
		ownership.ItemQuantity += amount
	case "decrement":
		if ownership.ItemQuantity == 0 {
		} else {
			ownership.ItemQuantity -= amount
			if ownership.ItemQuantity < 0 {
				ownership.ItemQuantity = 0
			}
		}
	case "set":
		ownership.ItemQuantity = amount
	default:
		return Error(c, 400, "Change type must be increment, decrement or set")
	}

	// Save new amount to the database and create response
	db.DB.Save(&ownership)

	ownershipDTO := DTO("ownership", ownership)
	return Success(c, "Item found", ownershipDTO)
}

/*
* Deletes an ownership from the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func OwnershipDelete(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	ownershipUID := c.Query("ownershipUID")

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err := RecordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	db.DB.Delete(&ownership)

	// Check for errors after the delete operation
	if result := db.DB.Delete(&ownership); result.Error != nil {
		return Error(c, 500, "There was an error deleting the ownership")
	}

	// Ownership successfully deleted
	return Success(c, "Ownership was successfully deleted")
}

/*
* Edits the fields of the ownership in the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func OwnershipEdit(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	ownershipUID := c.Query("ownershipUID")

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err := RecordExists("Ownership", result)

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
	return Success(c, "Ownership was successfully updated")
}

/*
* Creates an ownership in the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func OwnershipCreate(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)

	// convert uid to uint
	itemUID, err := strconv.ParseUint(c.Query("item_uid"), 10, 64)
	if err != nil {
		return Error(c, 400, "There was an error converting itemUID to Uint")
	}

	ownership, err := createOwnership(userUID, uint(itemUID))
	if err != nil {
		return Error(c, 400, err.Error())
	}

	ownershipDTO := DTO("ownershipUID", ownership.OwnershipUID)
	return Success(c, "Ownership was successfully created", ownershipDTO)
}

/*
* Sets the locatino of the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func OwnershipSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	locationQR := c.Query("location_qr")
	ownershipUID := c.Query("ownershipUID")

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err := RecordExists("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var ownership models.Ownership
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = RecordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	ownership.ItemLocation = location.LocationUID
	db.DB.Save(&ownership)

	// return success
	return Success(c, "Ownership set in "+location.LocationName)
}
