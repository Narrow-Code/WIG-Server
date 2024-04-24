package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		Parent:	       uuid.MustParse(db.DefaultLocationUUID),
		LocationUID: uuid.New(),
	}

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
	location.Parent = setLocation.Location.LocationUID
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

	// Parse request into data map
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := RecordExists("Ownership", result)
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
	return Success(c, "Location updated successfully")
}

/*
* Returns all ownerships and locations stored in a location.
*
* @param c The fiber context containing the HTTP request and esponse objects.
* @return error The error message, if there is any.
*/
func UnpackLocation( c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	locationUID := c.Query("locationUID")

	// Validate ownership
	var location models.Location
	result := db.DB.Where("location_uid = ? AND location_owner = ?", locationUID, user.UserUID).First(&location)
	code, err := RecordExists("Location", result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	ownerships, locations := GetAllFromLocation(location, user)

	ownershipDTO := DTO("ownerships", ownerships)
	locationDTO := DTO("locations", locations)

	return Success(c, "Unpacked", ownershipDTO, locationDTO)
}

/* 
* Searches for locations based on users query.
*
* @param c The FIber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
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

	if err := query.Find(&locations).Error; err != nil{
		return Error(c, 404, "Not found")
	}

	for i := range locations {
		preloadLocation(&locations[i])
	}

	locationDTO := DTO("locations", locations)
	return Success(c, "Items found", locationDTO)
}


/* 
* Returns the entire inventory for a user.
*
* @param c The FIber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func ReturnInventory(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)

	var locations models.Location
	db.DB.Where("location_uid = ?", db.DefaultLocationUUID).First(&locations)

	inventory := ReturnAllInventory(locations, user)
	inventoryDTO := DTO("inventory", inventory)
	

	return Success(c, "Inventory returned", inventoryDTO)
}
