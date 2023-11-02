package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"strconv"
	"github.com/gofiber/fiber/v2"
)


func CreateLocation(c *fiber.Ctx) error {	
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}
  
	// Initialize variables
        userUID := data["uid"]	
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	locationType := c.Params("type")
	
	// Check location type exists
	if locationType != "bin" && locationType != "bag" && locationType != "location" {
		return returnError(c, 400, "Location type is invalid") // TODO make message 
	}

	// convert uid to int
	userUIDInt, err := strconv.ParseUint(userUID, 10, 64)
	if err != nil {return returnError(c, 400, messages.ConversionError)}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Check for empty fields 
	if locationQR == "" {return returnError(c, 400, "QR Location required")} // TODO make message 
	if locationName == "" {return returnError(c, 400, "Location name required")} // TODO make message
	
	// Validate location QR code is not in use
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err = RecordInUse("Location QR", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Valide location name is not in use
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, userUID).First(&location)
	code, err = RecordInUse("Location Name", result)
	if err != nil {return returnError(c, code, err.Error())}

	// create location
	location = models.Location{
		LocationName: locationName,
		LocationOwner: uint(userUIDInt),
		LocationType: locationType,
		LocationQR: locationQR,
	}

	db.DB.Create(&location)

	return returnSuccess(c, "location added successfully") // TODO make message
}

func SetLocation(c *fiber.Ctx) error{
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}
  
	// Initialize variables
        userUID := data["uid"]	
	locationQR := data["location_qr"]
	ownershipUID := c.Query("ownershipUID")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Validate the QR code
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)
	code, err = RecordExists("Location QR", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Validate the ownership
	var ownership models.Ownership
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = RecordExists("Ownership", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Set the location and save
	ownership.ItemLocation = location.LocationUID
	db.DB.Save(&ownership)

	// return success
	return returnSuccess(c, "Ownership set in " + location.LocationName) // TODO make message
}
