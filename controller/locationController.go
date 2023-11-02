package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"fmt"
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
	locationQR := c.Query("location_qr")
	locationName := c.Query("location_name")
	locationType := c.Params("type")
	
	// Check location type exists
	if locationType != "bin" && locationType != "bag" && locationType != "location" {
		return returnError(c, 400, "Location type is invalid") // TODO make message 
	}

	// convert uid to int
	userUIDInt, err := strconv.ParseUint(userUID, 10, 64)
		if err != nil {
    		return returnError(c, 400, messages.ConversionError) 
	}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {
		return returnError(c, code, err.Error())
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

func SetLocation(c *fiber.Ctx) error{
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
  
        // Error with JSON request
	if err != nil {
		return returnError(c, 400, messages.ErrorParsingRequest)
	}
  
        userUID := data["uid"]	
	locationQR := data["location_qr"]
	ownershipUID := c.Query("ownershipUID")

	fmt.Println(ownershipUID)

	// convert uid to int
	ownershipUIDInt, err := strconv.ParseUint(ownershipUID, 10, 64)
		if err != nil {
    		return returnError(c, 400, messages.ConversionError) 
	}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {
		return returnError(c, code, err.Error())
	}
	// Check if QR code exists in users data
	var location models.Location
	result := db.DB.Where("location_qr = ? AND location_owner = ?", locationQR, userUID).First(&location)

	// If not return error
	if location.LocationQR != locationQR {
		return returnError(c, 400, "Location does not exist") // TODO make message
	}
        if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }

	// Check if ownership exists
	var ownership models.Ownership
	result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)

	// If not return error
	if ownership.OwnershipUID != uint(ownershipUIDInt) {
		return returnError(c, 400, "Ownership does not exist") // TODO make message
	}
	
        // If there is a connection error
        if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }

	// ownership.location = locationUID
	ownership.ItemLocation = location.LocationUID
	
	// save to db
	db.DB.Save(&ownership)

	// return success
	return returnSuccess(c, "Ownership set in " + location.LocationName) // TODO make message
}
