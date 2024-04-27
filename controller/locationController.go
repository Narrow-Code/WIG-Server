package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// LocationCreate creates a location using a QR code, a location name and the type of location.
func LocationCreate(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	var ownershipCheck models.Ownership
	user := c.Locals("user").(models.User)
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	log.Printf("controller#LocationCreate: User %d called LocationCreate", user.UserUID)

	// Check for empty fields
	utils.UserLog(c, "checking for empty fields")
	if locationQR == "" || locationName == "" {
		return Error(c, 400, "The locationQR or locationName field is empty")
	}

	// Check if QR exists as a location
	utils.UserLog(c, "checking if qr " + locationQR + " exists as a location")
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, user.UserUID).First(&location)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Check if QR exists as an ownership
	utils.UserLog(c, "checking if qr " + locationQR + " exists as an ownership")
	result = db.DB.Where("item_qr = ? AND item_owner = ?", locationQR, user.UserUID).First(&ownershipCheck)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate location name is not in use
	utils.UserLog(c, "checking if " + locationName + " is in use as a location name")
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, user.UserUID).First(&location)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Create location, add to DTO and return
	utils.Log("creating location " + locationName)
	location = createLocation(locationName, user, locationQR)
	locationDTO := DTO("location", &location)
	utils.Log("success")
	return success(c, "Location has been added successfully", locationDTO)
}

// LocationSetParent sets the location of a specific location.
func LocationSetParent(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	var setLocation models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("location_uid")
	setLocationUID := c.Query("set_location_uid")

	// Verify locations are not the same
	utils.UserLog(c, "checking that both locations are not the same")
	if locationUID == setLocationUID {
		return Error(c, 400, "Cannot set location in itself")
	}

	// Validate the location exists
	utils.UserLog(c, "validating that location exists")
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the set location exists
	utils.UserLog(c, "validating the set location exists")
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, user.UserUID).First(&setLocation)
	code, err = recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location parent and save
	location.Parent = setLocation.Location.LocationUID
	db.DB.Save(&location)
	utils.UserLog(c, location.LocationName + " is set in " + setLocation.LocationName)

	// return success
	utils.UserLog(c, "success")
	return success(c, location.LocationName+" set in "+setLocation.LocationName)
}

func LocationEdit(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var data map[string]string
	var location models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Parse request into data map
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate location exists
	utils.UserLog(c, "validating location exists")
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Add new fields
	location.LocationName = data["locationName"]
	location.LocationDescription = data["locationDescription"]
	location.LocationTags = data["locationTags"]
	location.LocationQR = data["qr"]

	// Save location to database and return
	db.DB.Save(&location)

	utils.UserLog(c, "success")
	return success(c, "Location updated successfully")
}

// Returns all ownerships and locations stored in a location.
func LocationUnpack(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var location models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Validate location exists
	utils.UserLog(c, "validating location exists")
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Get inventoryDTO, add to dto and return
	utils.UserLog(c, "unpacking location")
	inventoryDTO := getInventoryDTO(location, user)
	dto := DTO("inventory", inventoryDTO)
	utils.UserLog(c, "success")
	return success(c, "Unpacked", dto)
}

// Searches for locations based on users query.
func LocationSearch(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	user := c.Locals("user").(models.User)
	var locations []models.Location
	var data map[string]string

	// Parse JSON body
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}
	locationName := data["name"]
	locationTags := data["tags"]

	// Set tag format, split up by commas
	tagsFormat := strings.Split(strings.TrimSpace(locationTags), ",")

	// Add locationName and tags to query
	utils.UserLog(c, "adding locationName and tags to query")
	query := db.DB.Where("location_owner = ? AND location_name LIKE ?", user.UserUID, "%"+locationName+"%")
	for _, tag := range tagsFormat {
		query = query.Where("location_tags LIKE ?", "%"+tag+"%")
	}

	// Search for query
	utils.UserLog(c, "searching for locations")
	if err := query.Find(&locations).Error; err != nil {
		return Error(c, 404, "Not found")
	}

	// Preload locations with data
	utils.UserLog(c, "preloading locations")
	for i := range locations {
		preloadLocation(&locations[i])
	}

	// Add to DTO and return
	dto := DTO("locations", locations)
	utils.UserLog(c, "success")
	return success(c, "Items found", dto)
}

// Returns the entire inventory for a user.
func LocationGetInventory(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var locations models.Location
	user := c.Locals("user").(models.User)

	// Get default location
	utils.UserLog(c, "rerieving default location")
	db.DB.Where("location_uid = ?", db.DefaultLocationUUID).First(&locations)

	// Get Inventory dto, add to dto and return
	utils.UserLog(c, "unpacking inventory")
	inventoryDTO := getInventoryDTO(locations, user)
	dto := DTO("inventory", inventoryDTO)
	utils.UserLog(c, "success")
	return success(c, "Inventory returned", dto)
}
