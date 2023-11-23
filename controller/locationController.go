package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func LocationCreate(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return utils.Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	locationType := c.Params("type")

	// Check location type exists
	if locationType != "bin" && locationType != "bag" && locationType != "location" {
		return utils.Error(c, 400, messages.LocationTypeInvalid)
	}

	// convert uid to uint
	userUIDInt, err := strconv.ParseUint(userUID, 10, 64)
	if err != nil {
		return utils.Error(c, 400, messages.ConversionError)
	}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Check for empty fields
	if locationQR == "" {
		return utils.Error(c, 400, messages.LocationQRRequired)
	}
	if locationName == "" {
		return utils.Error(c, 400, messages.LocationNameRequired)
	}

	// Validate location QR code is not in use
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err = recordNotInUse("Location QR", result)
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Valide location name is not in use
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, userUID).First(&location)
	code, err = recordNotInUse("Location Name", result)
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// create location
	location = models.Location{
		LocationName:  locationName,
		LocationOwner: uint(userUIDInt),
		LocationType:  locationType,
		LocationQR:    locationQR,
	}

	db.DB.Create(&location)

	return utils.Success(c, messages.LocationAdded)
}

func LocationSetLocation(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return utils.Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	locationUID := c.Query("location_uid")
	setLocationUID := c.Query("set_location_uid")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Verify locations are not the same
	if locationUID == setLocationUID {
		return utils.Error(c, 400, messages.LocationSelfError)
	}

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, userUID).First(&location)
	code, err = recordExists("Location QR", result)
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Validate the ownership
	var setLocation models.Location
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, userUID).First(&setLocation)
	code, err = recordExists("Location", result)
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Set the location and save
	location.LocationLocation = &setLocation.LocationUID
	db.DB.Save(&location)

	// return success
	return utils.Success(c, location.LocationName+" set in "+setLocation.LocationName)
}

func LocationEdit(c *fiber.Ctx) error {
	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return utils.Error(c, 400, messages.ErrorParsingRequest)
	}

	// Initialize variables
	userUID := data["uid"]
	locationUID := c.Query("locationUID")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, userUID).First(&location)
	code, err = recordExists("Ownership", result)
	if err != nil {
		return utils.Error(c, code, err.Error())
	}

	// Add new fields
	location.LocationName = c.Query("location_name")
	location.LocationDescription = c.Query("location_description")
	location.LocationTags = c.Query("location_tags")

	db.DB.Save(&location)

	// Ownership successfully updated
	return utils.Success(c, messages.LocationUpdated)
}
