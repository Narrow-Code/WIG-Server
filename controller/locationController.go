package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func LocationCreate(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	locationType := c.Params("type")
	log.Printf("controller#LocationCreate: User %s called LocationCreate", userUID)

	// Check location type exists
	if locationType != "bin" && locationType != "bag" && locationType != "area" {
		return Error(c, 400, "Location type must be bin, bag or area")
	}

	// convert uid to uint
	userUIDInt, err := strconv.ParseUint(userUID, 10, 64)
	if err != nil {
		return Error(c, 400, "Error converting userUID to uint")
	}

	// Check for empty fields
	if locationQR == "" || locationName == "" {
		return Error(c, 400, "The locationQR or locationName field is empty")
	}

	// Validate location QR code is not in use
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err := recordNotInUse("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Valide location name is not in use
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, userUID).First(&location)
	code, err = recordNotInUse("Location Name", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// create location
	location = models.Location{
		LocationName:  locationName,
		LocationOwner: uint(userUIDInt),
		LocationType:  locationType,
		LocationQR:    locationQR,
	}

	db.DB.Create(&location)

	return Success(c, messages.LocationAdded)
}

func LocationSetLocation(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	locationUID := c.Query("location_uid")
	setLocationUID := c.Query("set_location_uid")

	// Verify locations are not the same
	if locationUID == setLocationUID {
		return Error(c, 400, messages.LocationSelfError)
	}

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, userUID).First(&location)
	code, err := RecordExists("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var setLocation models.Location
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, userUID).First(&setLocation)
	code, err = RecordExists("Location", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	location.LocationLocation = &setLocation.LocationUID
	db.DB.Save(&location)

	// return success
	return Success(c, location.LocationName+" set in "+setLocation.LocationName)
}

func LocationEdit(c *fiber.Ctx) error {
	// Initialize variables
	userUID := c.Locals("uid").(string)
	locationUID := c.Query("locationUID")

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, userUID).First(&location)
	code, err := RecordExists("Ownership", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Add new fields
	location.LocationName = c.Query("location_name")
	location.LocationDescription = c.Query("location_description")
	location.LocationTags = c.Query("location_tags")

	db.DB.Save(&location)

	// Ownership successfully updated
	return Success(c, messages.LocationUpdated)
}
