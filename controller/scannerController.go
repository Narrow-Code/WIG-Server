package controller

import (
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/structs"
	"WIG-Server/upcitemdb"
	"WIG-Server/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
GetBarcode handles the functionality of returning any ownerships and items back after scanning a barcode.

@param c *fiber.Ctx
*/
func ScanBarcode(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
        if err != nil {return utils.NewError(c, 400, messages.ErrorParsingRequest)}

	// Initialize variables
	uid := data["uid"]
	barcode := c.Query("barcode")
	
	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return utils.NewError(c, code, err.Error())}

	// Validate barcode
	if barcode == "" {return utils.NewError(c, 400, messages.BarcodeMissing)}
	barcodeCheck, err := strconv.Atoi(barcode)
	if err != nil || barcodeCheck < 0 {return utils.NewError(c, 400, messages.BarcodeIntError)}

	// Check if item exists in local database
	var item models.Item
        result := db.DB.Where("barcode = ?", barcode).First(&item) 

        // If item isn't found, check api and add to 
        if result.Error == gorm.ErrRecordNotFound {
		upcitemdb.GetBarcode(barcode)
		result = db.DB.Where("barcode = ?", barcode).First(&item)
		if result.Error == gorm.ErrRecordNotFound {
			return utils.NewError(c, 400, messages.ItemNotFound)
		}
        }

	// If there is a connection error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return utils.NewError(c, 400, messages.ErrorWithConnection)
        }
	
	// Search Ownership by barcode
	var ownerships []models.Ownership
	result = db.DB.Where("item_barcode = ? AND item_owner = ?", barcode, uid).Find(&ownerships)

	// If no ownership exists, create ownership
	if len(ownerships) == 0 {
		ownership, err := createOwnership(uid, item.ItemUid)
		if err != nil {return utils.NewError(c, 400, err.Error())}
		
		var ownershipResponses []structs.OwnershipResponse
		ownershipResponses = append(ownershipResponses, getOwnershipReponse(ownership))
		return c.Status(200).JSON(
			fiber.Map{
				"success":true,
				"message":"Created new ownership",
				"item":item.Name,
				"barcode":item.Barcode,
				"brand":item.Brand,
				"image":item.Image,
				"owner":uid,
				"ownership":ownershipResponses})
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
				"barcode":item.Barcode,
				"brand":item.Brand,
				"image":item.Image,
				"owner":uid,
				"ownership":ownershipResponses})
}

/*
CheckQR takes a QR code as parameter, and checks whether it is an item, location or a unused QR.

@param c *fiber.Ctx - The fier context containing the HTTP request and response objects.
@return error - An error that occured during the process or if the token does not match
*/
func ScanCheckQR(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)
        if err != nil {return utils.NewError(c, 400, messages.ErrorParsingRequest)}

	// Initialize variables
        uid := data["uid"]
        qr := c.Query("qr")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {return utils.NewError(c, code, err.Error())}

	// Check for empty fields
	if qr == "" {return utils.NewError(c, 400, messages.QRMissing)}
  
        // Check if qr exists as location
        var location models.Location
        result := db.DB.Where("location_qr = ? AND location_owner = ?", qr, uid).First(&location)
	if location.LocationUID != 0 {return returnSuccess(c, messages.Location)
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {return utils.NewError(c, 400, messages.ErrorWithConnection)}

	// Check if qr exists as ownership
	var ownership models.Ownership
	result = db.DB.Where("item_qr = ? AND item_owner = ?", qr, uid).First(&ownership)
	if ownership.OwnershipUID != 0 {return returnSuccess(c, messages.Ownership)}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {return utils.NewError(c, 400, messages.ErrorWithConnection)}

	return returnSuccess(c, messages.New)
}


