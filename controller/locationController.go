package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

/*
* Creates a location using a QR code, a location name and the type of location.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func LocationCreate(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	log.Printf("controller#LocationCreate: User %d called LocationCreate", user.UserUID)

	// Check for empty fields
	if locationQR == "" || locationName == "" {
		return Error(c, 400, "The locationQR or locationName field is empty")
	}

	// Validate location QR code is not in use
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, user.UserUID).First(&location)
	code, err := recordNotInUse("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Valide location name is not in use
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, user.UserUID).First(&location)
	code, err = recordNotInUse("Location Name", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// create location
	location = models.Location{
		LocationName:  locationName,
		LocationOwner: user.UserUID,
		LocationQR:    locationQR,
	}

//	preloadLocation(&location)
	db.DB.Create(&location)
	preloadLocation(&location)
	locationDTO := DTO("location", &location)

	return Success(c, "Location has been added successfully", locationDTO)
}

/*
* Sets the location of a specific location.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
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
	code, err := RecordExists("Location QR", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Validate the ownership
	var setLocation models.Location
	result = db.DB.Where("location_uid = ? AND location_owner = ?", setLocationUID, user.UserUID).First(&setLocation)
	code, err = RecordExists("Location", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Set the location and save
	location.Parent = &setLocation.LocationUID
	db.DB.Save(&location)

	// return success
	return Success(c, location.LocationName+" set in "+setLocation.LocationName)
}

/*
* Edits the fields of the location.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func LocationEdit(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
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
	return Success(c, "Location updated successfully")
}
