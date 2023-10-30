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
	 // Get parameters
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

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
                return returnError(c, 400, messages.ErrorWithConnection)

        }

	if item.Barcode == "" {
        	return returnError(c, 400, "Item not found")
    	}

	return c.Status(200).JSON(
                        fiber.Map{
                                "success":true,
                                "message":"Barcode added",       
				"title":item.Name,
				"brand":item.Brand,
				"image":item.Image,
				"description":item.ItemDesc})

}

