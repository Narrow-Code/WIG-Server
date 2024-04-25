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
	createLocation(locationName, user, locationQR)
	locationDTO := DTO("location", &location)
	return success(c, "Location has been added successfully", locationDTO)
}

// LocationSetLocation sets the location of a specific location.
func LocationSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationUID := c.Query("location_uid")
	setLocationUID := c.Query("set_location_uid")

	// Verify locations are not the same
	if locationUID == setLocationUID {
		return Error(c, 400, "Cannot set location in itself")
	}

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var setLocation models.Location
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, user.UserUID).First(&setLocation)
	code, err = recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	location.Parent = setLocation.Location.LocationUID
	db.DB.Save(&location)

	// return success
	return success(c, location.LocationName+" set in "+setLocation.LocationName)
}

// LocationEdit edits the fields of the location.
func LocationEdit(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate ownership
	var location models.Location
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

	db.DB.Save(&location)

	// Ownership successfully updated
	return success(c, "Location updated successfully")
}

// Returns all ownerships and locations stored in a location.
func UnpackLocation(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	ownerships, locations := GetAllFromLocation(location, user)

	ownershipDTO := DTO("ownerships", ownerships)
	locationDTO := DTO("locations", locations)

	return success(c, "Unpacked", ownershipDTO, locationDTO)
}

// Searches for locations based on users query.
func LocationSearch(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing the JSON")
	}
	name := data["name"]
	tags := data["tags"]
	var locations []models.Location

	tagsFormat := strings.Split(strings.TrimSpace(tags), ",")

	query := db.DB.Where("location_owner = ? AND location_name LIKE ?", user.UserUID, "%"+name+"%")

	for _, tag := range tagsFormat {
		query = query.Where("location_tags LIKE ?", "%"+tag+"%")
	}

	if err := query.Find(&locations).Error; err != nil {
		return Error(c, 404, "Not found")
	}

	for i := range locations {
		preloadLocation(&locations[i])
	}

	locationDTO := DTO("locations", locations)
	return success(c, "Items found", locationDTO)
}

// Returns the entire inventory for a user.
func ReturnInventory(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)

	var locations models.Location
	db.DB.Where("location_uid = ?", db.DefaultLocationUUID).First(&locations)

	inventory := ReturnAllInventory(locations, user)
	inventoryDTO := DTO("inventory", inventory)

	return success(c, "Inventory returned", inventoryDTO)
}
