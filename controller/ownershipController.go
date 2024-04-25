package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"strings"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// OwnershipQuantity changes the quantity of an ownership, using increment, decrement or setter method.
func OwnershipQuantity(c *fiber.Ctx) error {
	// Initialize variables
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
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

	// Retreive the ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Make changes based on changeType
	switch changeType {
	case "increment":
		ownership.ItemQuantity += amount
	case "decrement":
		if ownership.ItemQuantity != 0 {
			ownership.ItemQuantity -= amount
		}
		if ownership.ItemQuantity < 0 {
			ownership.ItemQuantity = 0
		}
	case "set":
		ownership.ItemQuantity = amount
	default:
		return Error(c, 400, "Change type must be increment, decrement or set")
	}

	// Save new amount to the database, preload, make dto and return
	db.DB.Save(&ownership)
	preloadOwnership(&ownership)
	dto := DTO("ownership", ownership)
	return success(c, "Item found", dto)
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
	user := c.Locals("user").(models.User)
	ownershipUID := c.Query("ownershipUID")

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	db.DB.Delete(&ownership)

	// Check for errors after the delete operation
	if result := db.DB.Delete(&ownership); result.Error != nil {
		return Error(c, 500, "There was an error deleting the ownership")
	}

	// Ownership successfully deleted
	return success(c, "Ownership was successfully deleted")
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
	user := c.Locals("user").(models.User)
	ownershipUID := c.Query("ownershipUID")

	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err := recordExists(result)

	if err != nil {
		return Error(c, code, err.Error())
	}

	// Add new fields
	ownership.CustomItemName = data["customItemName"]
	ownership.CustItemImg = data["customItemImg"]
	ownership.OwnedCustDesc = data["customItemDescription"]
	ownership.ItemTags = data["itemTags"]
	ownership.ItemQR = data["qr"]

	db.DB.Save(&ownership)

	// Ownership successfully updated
	return success(c, "Ownership was successfully updated")
}

/*
* Creates an ownership in the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func OwnershipCreateNoItem(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	user := c.Locals("user").(models.User)
	qr := data["qr"]
	name := data["name"]

	if data["qr"] == "" && data["name"] == "" {
		return Error(c, 400, "Missing field qr or name")
	}

	var ownershipCheck models.Ownership
	result := db.DB.Where("item_qr = ? AND item_owner = ?", qr, user.UserUID).First(&ownershipCheck)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	var locationCheck models.Location
	result = db.DB.Where("location_qr = ? AND location_owner = ?", qr, user.UserUID).First(&locationCheck)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	result = db.DB.Where("custom_item_name = ? AND item_owner = ?", name, user.UserUID).First(&ownershipCheck)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	var item models.Item

	ownership, err := createOwnership(user.UserUID, item, qr, name)
	if err != nil {
		return Error(c, 400, err.Error())
	}

	preloadOwnership(&ownership)
	ownershipDTO := DTO("ownership", ownership)
	return success(c, "Ownership was successfully created", ownershipDTO)
}

/*
* Sets the location of the ownership in the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func OwnershipSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationQR := c.Query("location_qr")
	ownershipUID := c.Query("ownershipUID")

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var ownership models.Ownership
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err = recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	ownership.ItemLocation = location.LocationUID
	db.DB.Save(&ownership)

	// return success
	return success(c, "Ownership set in "+location.LocationName)
}

/*
* Searches for items based on users query.
*
* @param c The FIber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func OwnershipSearch(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}
	name := data["name"]
	tags := data["tags"]
	var ownerships []models.Ownership

	tagsFormat := strings.Split(strings.TrimSpace(tags), ",")

	query := db.DB.Where("item_owner = ? AND custom_item_name LIKE ?", user.UserUID, "%"+name+"%")

	for _, tag := range tagsFormat {
		query = query.Where("item_tags LIKE ?", "%"+tag+"%")
	}

	if err := query.Find(&ownerships).Error; err != nil {
		return Error(c, 404, "Not found")
	}

	for i := range ownerships {
		preloadOwnership(&ownerships[i])
	}

	ownershipDTO := DTO("ownership", ownerships)
	return success(c, "Items found", ownershipDTO)
}
