/* Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
 */
package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/structs"
	"WIG-Server/upcitemdb"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
GetBarcode handles the functionality of returning any ownerships and items back after scanning a barcode.

@param c *fiber.Ctx
*/
func GetBarcode(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)

        // Error with JSON request
        if err != nil {
                return returnError(c, 400, messages.ErrorParsingRequest)
        }

	uid := data["uid"]
	barcode := c.Query("barcode")

	// Validate Token
	err = validateToken(c, uid, data["token"])	
	if err == nil {
		return validateToken(c, uid, data["token"])
	}

	// Check if item exists in cache
	var item models.Item
        result := db.DB.Where("barcode = ?", barcode).First(&item)
        
        // If item isn't found, check api and add to cache
        if result.Error == gorm.ErrRecordNotFound {
		upcitemdb.GetBarcode(barcode)
		result = db.DB.Where("barcode = ?", barcode).First(&item)               
        }

	// If there is a connection error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }
	
	// If the barcode is empty retrun error
	if item.Barcode == "" {
        	return returnError(c, 404, messages.ItemNotFound)
    	}

	// Search Ownership by barcode
	var ownerships []models.Ownership
	result = db.DB.Where("item_barcode = ? AND item_owner = ?", barcode, uid).Find(&ownerships)

	// Convert uid to int
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return returnError(c, 400, "Error converting UID") // TODO make message
	}

	// If no ownership exists, create ownership
	if len(ownerships) == 0 {
		ownership := models.Ownership{
               		ItemOwner:uint(uidInt),
			ItemBarcode:barcode,
			ItemQuantity:1,
   		}
		db.DB.Create(&ownership)

		return c.Status(200).JSON(
			fiber.Map{
				"success":true,
				"message":"Created new ownership",
				"title":item.Name,
				"barcode":item.Barcode,
				"brand":item.Brand,
				"image":item.Image,
				"owner":ownership.ItemOwner})
	}

	// If ownerships exist, return as slice
	var ownershipResponses []structs.OwnershipResponse
	for _, ownership := range ownerships {
		ownershipResponse := getOwnershipReponse(ownership)
		ownershipResponses = append(ownershipResponses, ownershipResponse)	
	}

	return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":"Item found",       
				"item":item.Name,
				"brand":item.Brand,
				"image":item.Image,
				"owner":uid,
				"ownership":ownershipResponses})
}

/*
IncrementOwnership increases the ownerships quantity by the designated value

@param c *fiber.Ctx
*/
func ChangeQuantity(c *fiber.Ctx) error {
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
  
        // Error with JSON request
	if err != nil {
		return returnError(c, 400, messages.ErrorParsingRequest)
	}
  
        userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")
	amountStr := c.Query("amount")
	changeType := c.Params("type")

	amount, err := strconv.Atoi(amountStr)
		if err != nil {
    		return returnError(c, 400, "Error converting int") // TODO make message
	}

        // Validate Token
        err = validateToken(c, userUID, data["token"])      
        if err == nil {
                return validateToken(c, userUID, data["token"])
        }

	// Search for ownership UID and pair with user UID to make sure they match
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)

	// If item is not found
	if result.Error == gorm.ErrRecordNotFound {
                return returnError(c, 404, messages.ItemNotFound)              
        }

        // If there is a connection error
        if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }

	if changeType == "increment" {
		// Add the incremental value to quantity
		ownership.ItemQuantity = ownership.ItemQuantity + amount 
	} else if changeType == "decrement" {
		ownership.ItemQuantity = ownership.ItemQuantity - amount
		if ownership.ItemQuantity < 0 {
			ownership.ItemQuantity = 0
		}
	} else if changeType == "set" {
		if amount < 0 {
			return returnError(c, 400, "Cannot set negative numbers")
		}
		ownership.ItemQuantity = amount
	}
	

	// Save new amount to the database
	db.DB.Save(&ownership)

	// Create ownership response
	ownershipResponse := getOwnershipReponse(ownership)

	// Return success
	return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":"Item found",       
                               	"ownership": ownershipResponse})
}
