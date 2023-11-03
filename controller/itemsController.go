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
        if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}

	// Initialize variables
	uid := data["uid"]
	barcode := c.Query("barcode")
	
	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Validate barcode
	if barcode == "" {return returnError(c, 400, "Barcode required")} // TODO add message
	barcodeCheck, err := strconv.Atoi(barcode)
	if err != nil || barcodeCheck < 0 {return returnError(c, 400, "Barcode must be of int value")}

	// Check if item exists in local database
	var item models.Item
        result := db.DB.Where("barcode = ?", barcode).First(&item) 

        // If item isn't found, check api and add to 
        if result.Error == gorm.ErrRecordNotFound {
		upcitemdb.GetBarcode(barcode)
		result = db.DB.Where("barcode = ?", barcode).First(&item)               
        }

	// If there is a connection error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }
	
	// Search Ownership by barcode
	var ownerships []models.Ownership
	result = db.DB.Where("item_barcode = ? AND item_owner = ?", barcode, uid).Find(&ownerships)

	// If no ownership exists, create ownership
	if len(ownerships) == 0 {
		err = createOwnership(uid, barcode)
		if err != nil {return returnError(c, 400, err.Error())}
		return c.Status(200).JSON(
			fiber.Map{
				"success":true,
				"message":"Created new ownership",
				"title":item.Name,
				"barcode":item.Barcode,
				"brand":item.Brand,
				"image":item.Image,
				"owner":uid})
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
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}
  
	// Initialize variables
        userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")
	amountStr := c.Query("amount")
	changeType := c.Params("type")

	// Convert amount to int
	amount, err := strconv.Atoi(amountStr)
	if err != nil {return returnError(c, 400, messages.ConversionError)}
	if amount < 0 {return returnError(c, 400, messages.NegativeError)}

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Valide and retreive the ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {return returnError(c, code, err.Error())}

	// Check type of change
	switch changeType {
	case "increment":
		ownership.ItemQuantity += amount
	case "decrement":
		ownership.ItemQuantity -= amount
		if ownership.ItemQuantity < 0 {ownership.ItemQuantity = 0}
	case "set":
		ownership.ItemQuantity = amount;
	default:
		return returnError(c, 400, "Invalid change type") // TODO message
	}

	// Save new amount to the database and create response
	db.DB.Save(&ownership)
	ownershipResponse := getOwnershipReponse(ownership)

	// Return success
	return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":"Item found",       
                               	"ownership": ownershipResponse})
}

func DeleteOwnership(c *fiber.Ctx) error {
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}
  
	// Initialize variables
        userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")
	
	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {return returnError(c, code, err.Error())}
	
	db.DB.Delete(&ownership)

	// Check for errors after the delete operation
	if result := db.DB.Delete(&ownership); result.Error != nil {
    		return returnError(c, 500, messages.ErrorDeletingOwnership)
	}

	// Ownership successfully deleted
	return returnSuccess(c, "Ownership deleted successfully")
}

func EditOwnership(c *fiber.Ctx) error {
        // Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
	if err != nil {return returnError(c, 400, messages.ErrorParsingRequest)}
  
	// Initialize variables
        userUID := data["uid"]
	ownershipUID := c.Query("ownershipUID")
	changeField := c.Params("field")
	
	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return returnError(c, code, err.Error())}

	// Validate ownership
	var ownership models.Ownership
	result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownershipUID, userUID).First(&ownership)
	code, err = recordExists("Ownership", result)
	if err != nil {return returnError(c, code, err.Error())}
	
	// Add new fields
	if changeField == "name" {ownership.CustomItemName = c.Query("custom_item_name")}
	if changeField == "img" {ownership.CustItemImg = c.Query("custom_item_img")}
	if changeField == "description" {ownership.OwnedCustDesc = c.Query("custom_item_description")}
	if changeField == "tags" {ownership.ItemTags = c.Query("item_tags")}

	db.DB.Save(&ownership)

	// Ownership successfully updated
	return returnSuccess(c, changeField + " updated")
}


