package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// OwnershipQuantity changes the quantity of an ownership, using increment, decrement or setter method.
func OwnershipQuantity(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
	ownershipUID := c.Query("ownershipUID")
	amountStr := c.Query("amount")
	changeType := c.Params("type")

	// Convert amount to int
	utils.UserLog(c, "converting " + amountStr + " to an int")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return Error(c, 400, "There was an error converting amount to Int")
	}
	if amount < 0 {
		return Error(c, 400, "Amount cannot be negative")
	}

	// Retreive the ownership
	utils.UserLog(c, "retrieving the ownership")
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}
	utils.UserLog(c, ownership.CustomItemName + " retrieved")

	// Make changes based on changeType
	switch changeType {
	case "increment":
		ownership.ItemQuantity += amount
		utils.UserLog(c, ownership.CustomItemName + " incremented")
	case "decrement":
		if ownership.ItemQuantity != 0 {
			ownership.ItemQuantity -= amount
		}
		if ownership.ItemQuantity < 0 {
			ownership.ItemQuantity = 0
			utils.UserLog(c, ownership.CustomItemName + " incremented")
		}
	case "set":
		ownership.ItemQuantity = amount
		utils.UserLog(c, ownership.CustomItemName + " set")
	default:
		return Error(c, 400, "Change type must be increment, decrement or set")
	}

	// Save new amount to the database, preload, make dto and return
	db.DB.Save(&ownership)
	utils.UserLog(c, "ownership saved")
	preloadOwnership(&ownership)
	dto := DTO("ownership", ownership)
	utils.UserLog(c, " success")
	return success(c, "Item found", dto)
}

// OwnershipDelete deletes an ownership from the database.
func OwnershipDelete(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
	ownershipUID := c.Query("ownershipUID")

	// Validate ownership
	utils.UserLog(c, "validating the ownership")
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Delete ownership and return
	db.DB.Delete(&ownership)
	if result := db.DB.Delete(&ownership); result.Error != nil {
		return Error(c, 500, "There was an error deleting the ownership")
	}
	utils.UserLog(c, "success")
	return success(c, "Ownership was successfully deleted")
}

// OwnershipEdit edits the fields of the ownership in the database.
func OwnershipEdit(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var ownership models.Ownership
	var data map[string]string
	user := c.Locals("user").(models.User)
	ownershipUID := c.Query("ownershipUID")

	// Parse request into data map
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate ownership
	utils.UserLog(c, "validating the ownership")
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

	// Save ownership and return
	db.DB.Save(&ownership)
	utils.UserLog(c, "success")
	return success(c, "Ownership was successfully updated")
}

// OwnershipCreate creates an ownership in the database.
func OwnershipCreate(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var ownershipCheck models.Ownership
	var data map[string]string
	var locationCheck models.Location
	var item models.Item

	// Parse JSON body
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	user := c.Locals("user").(models.User)
	qr := data["qr"]
	name := data["name"]

	// Check if fields are empty
	utils.UserLog(c, "validating fields are not empty")
	if qr == "" && name == "" {
		return Error(c, 400, "Missing field qr or name")
	}

	// Check if QR exists as an ownership
	utils.UserLog(c, "validating QR is not in use as ownership")
	result := db.DB.Where("item_qr = ? AND item_owner = ?", qr, user.UserUID).First(&ownershipCheck)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check if QR exists as location
	utils.UserLog(c, "validating QR is not in use as location")
	result = db.DB.Where("location_qr = ? AND location_owner = ?", qr, user.UserUID).First(&locationCheck)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check if custom item is already in use
	utils.UserLog(c, "validating custom item name is not in use")
	result = db.DB.Where("custom_item_name = ? AND item_owner = ?", name, user.UserUID).First(&ownershipCheck)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Create ownership
	utils.UserLog(c, "creating ownership")
	ownership, err := createOwnership(user.UserUID, item, qr, name)
	if err != nil {
		return Error(c, 400, err.Error())
	}

	// Preload ownership, add to dto and return
	preloadOwnership(&ownership)
	dto := DTO("ownership", ownership)
	utils.UserLog(c, "success")
	return success(c, "Ownership was successfully created", dto)
}

// OwnershipSetLocation sets the location of the ownership in the database.
func OwnershipSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	var ownership models.Ownership
	user := c.Locals("user").(models.User)
	locationQR := c.Query("location_qr") // TODO fix to be location UID
	ownershipUID := c.Query("ownershipUID")

	// Validate the location
	utils.UserLog(c, "validating the existance of location")
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	utils.UserLog(c, "validating the existance of ownership")
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, user.UserUID).First(&ownership)
	code, err = recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location, save and return
	utils.UserLog(c, "setting location")
	ownership.ItemLocation = location.LocationUID
	db.DB.Save(&ownership)
	utils.UserLog(c, "success")
	return success(c, "Ownership set in "+location.LocationName)
}

// OwnershipSearch searches for items based on users query.
func OwnershipSearch(c *fiber.Ctx) error {
	// Initialize variables
	var ownerships []models.Ownership
	var data map[string]string
	user := c.Locals("user").(models.User)

	// Parse JSON body
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}
	name := data["name"]
	tags := data["tags"]

	// Split up tags by commas
	tagsFormat := strings.Split(strings.TrimSpace(tags), ",")

	// Add ownership name and tags to query
	query := db.DB.Where("item_owner = ? AND custom_item_name LIKE ?", user.UserUID, "%"+name+"%")
	for _, tag := range tagsFormat {
		query = query.Where("item_tags LIKE ?", "%"+tag+"%")
	}

	// Search for query
	if err := query.Find(&ownerships).Error; err != nil {
		return Error(c, 404, "Not found")
	}

	// Preload locations with data
	for i := range ownerships {
		preloadOwnership(&ownerships[i])
	}

	// Add to dto and return
	dto := DTO("ownership", ownerships)
	return success(c, "Items found", dto)
}
