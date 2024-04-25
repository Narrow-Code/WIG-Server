package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// LocationCreate creates a location using a QR code, a location name and the type of location.
func LocationCreate(c *fiber.Ctx) error {
	// Initialize variables
	var location models.Location
	user := c.Locals("user").(models.User)
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	log.Printf("controller#LocationCreate: User %d called LocationCreate", user.UserUID)

	// Check for empty fields
	if locationQR == "" || locationName == "" {
		return Error(c, 400, "The locationQR or locationName field is empty")
	}

	// Validate location QR code is not in use
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, user.UserUID).First(&location)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate location name is not in use
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, user.UserUID).First(&location)
	code, err = recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Create location, add to DTO and return
	location = createLocation(locationName, user, locationQR)
	locationDTO := DTO("location", &location)
	return success(c, "Location has been added successfully", locationDTO)
}

// LocationSetLocation sets the location of a specific location.
func LocationSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	var location models.Location
	var setLocation models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("location_uid")
	setLocationUID := c.Query("set_location_uid")

	// Verify locations are not the same
	if locationUID == setLocationUID {
		return Error(c, 400, "Cannot set location in itself")
	}

	// Validate the location exists
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the set location exists
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, user.UserUID).First(&setLocation)
	code, err = recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location parent and save
	location.Parent = setLocation.Location.LocationUID
	db.DB.Save(&location)

	// return success
	return success(c, location.LocationName+" set in "+setLocation.LocationName)
}

func LocationEdit(c *fiber.Ctx) error {
	// Initialize variables
	var data map[string]string
	var location models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Parse request into data map
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate location exists
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
	return success(c, "Location updated successfully")
}

// Returns all ownerships and locations stored in a location.
func UnpackLocation(c *fiber.Ctx) error {
	// Initialize variables
	var location models.Location
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Validate location exists
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Get inventoryDTO, add to dto and return
	inventoryDTO := getInventoryDTO(location, user)
	dto := DTO("inventory", inventoryDTO)
	return success(c, "Unpacked", dto)
}

// Searches for locations based on users query.
func LocationSearch(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	var locations []models.Location
	var data map[string]string

	// Parse JSON body	
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}

	// Unpack variables
	locationName := data["name"]
	locationTags := data["tags"]

	// Set tag format, split up by commas
	tagsFormat := strings.Split(strings.TrimSpace(locationTags), ",")

	// Add locationName and tags to query
	query := db.DB.Where("location_owner = ? AND location_name LIKE ?", user.UserUID, "%"+locationName+"%")
	for _, tag := range tagsFormat {
		query = query.Where("location_tags LIKE ?", "%"+tag+"%")
	}

	// Search for query
	if err := query.Find(&locations).Error; err != nil {
		return Error(c, 404, "Not found")
	}

	// Preload locations with data
	for i := range locations {
		preloadLocation(&locations[i])
	}

	// Add to DTO and return
	dto := DTO("locations", locations)
	return success(c, "Items found", dto)
}

// Returns the entire inventory for a user.
func ReturnInventory(c *fiber.Ctx) error {
	// Initialize variables
	var locations models.Location
	user := c.Locals("user").(models.User)

	// Get default location
	db.DB.Where("location_uid = ?", db.DefaultLocationUUID).First(&locations)

	// Get Inventory dto, add to dto and return
	inventoryDTO := getInventoryDTO(locations, user)
	dto := DTO("inventory", inventoryDTO)
	return success(c, "Inventory returned", dto)
}
