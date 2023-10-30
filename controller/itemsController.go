/* Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
*/
package controller

import (
        "github.com/gofiber/fiber/v2"
        "WIG-Server/db"
	"WIG-Server/models"
        "WIG-Server/messages"
	"WIG-Server/upcitemdb"
        "gorm.io/gorm"
)

func GetBarcode(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)

        // Error with JSON request
        if err != nil {
                return returnError(c, 400, messages.ErrorParsingRequest)
        }

	// Validate Token
	err = validateToken(c, data["uid"], data["token"])	
	if err == nil {
		return validateToken(c, data["uid"], data["token"])
	}

	// Get barcode parameter
        barcode:= c.Query("barcode")

	// Check if item exists in cache
	var item models.Item
        result := db.DB.Where("barcode = ?", barcode).First(&item)
        
        // If item isn't found, check api and add to cache
        if result.Error == gorm.ErrRecordNotFound {
		// If item isn't found check API and add to cache
		upcitemdb.GetBarcode(barcode)
		result = db.DB.Where("barcode = ?", barcode).First(&item)               
        }

	// If there is a connection error
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)
        }

	if item.Barcode == "" {
        	return returnError(c, 404, messages.ItemNotFound)
    	}

	return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":"Item found",       
				"title":item.Name,
				"brand":item.Brand,
				"image":item.Image,
				"description":item.ItemDesc})
}

