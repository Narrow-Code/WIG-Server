package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func CreateLocation(c *fiber.Ctx) error {	
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
  
        // Error with JSON request
	if err != nil {
		return returnError(c, 400, messages.ErrorParsingRequest)
	}
  
        userUID := data["uid"]	
	locationQR := data["location_qr"]
	locationName := data["location_name"]
	locationType := c.Params("type")

	// convert uid to int
	userUIDInt, err := strconv.ParseUint(userUID, 10, 64)
		if err != nil {
    		return returnError(c, 400, messages.ConversionError) 
	}

        // Validate Token
        err = validateToken(c, userUID, data["token"])      
        if err == nil {
                return validateToken(c, userUID, data["token"])
        }

	// Check for valid entries
	if locationQR == "" {
		return returnError(c, 400, "QR Location required") // TODO MAKE MESSAGE
	} 
	
	if locationName == "" {
		return returnError(c, 400, "Location name required") // TODO MAKE MESSAGE
	} 
	
	// Check if QR code exists in users data
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)

	if location.LocationQR == locationQR {
		return returnError(c, 400, "Location already exists") // TODO make message
	}
	
        // If there is a connection error
        if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }

	// check if location name exists
	result = db.DB.Where("location_name = ? AND location_owner = ?", locationName, userUID).First(&location)
	
	if location.LocationName == locationName {
		return returnError(c, 400, "Location already exists") // TODO make message
	}

        // If there is a connection error
        if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
	}

	// create location
	location = models.Location{
		LocationName: locationName,
		LocationOwner: uint(userUIDInt),
		LocationType: locationType,
		LocationQR: locationQR,
	}

	// add to database
	db.DB.Create(&location)	

	// fix returnSuccess
	return returnSuccess(c, "location added successfully") // TODO make message
}
