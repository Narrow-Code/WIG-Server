// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/components"
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/structs"
	"fmt"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
returnError returns the given error code, success status and message through fiber to the application.

@param c c *fiber.Ctx - The Fiber context containing the HTTP request and response objects.
@param code The error code to return via fiber
@param The error message to retrun via fiber

@return error - An error, if any, that occurred during the registration process.
*/
func returnError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"success":false,
		"message":message})
}

func returnSuccess (c *fiber.Ctx, message string) error {
	return c.Status(200).JSON(fiber.Map{
		"success":true,
		"message":message})
}

/*
validateToken checks if a users UID and token match and are valid.

@param c *fiber.Ctx - The fier context containing the HTTP request and response objects.
@param UID The users UID
@param token The users authentication token

@return error - An error that occured during the process or if the token does not match
*/
func validateToken(c *fiber.Ctx, uid string, token string) (int, error){
	// Check if UID and token exist
        if uid == "" {
                return 400, errors.New(messages.UIDEmpty)
        }

        if token == "" {
                return 400, errors.New(messages.TokenEmpty)
        }

        // Query for UID
        var user models.User
        result := db.DB.Where("user_uid = ?", uid).First(&user)

        // Check if UID was found
        if result.Error == gorm.ErrRecordNotFound {
                return 404, errors.New("UID " + messages.RecordNotFound)

        } else if result.Error != nil {
                return 400, errors.New(messages.ErrorWithConnection)
	}

        // Validate token
        if !components.ValidateToken(user.Username, user.Hash, token) {
                return 400, errors.New(messages.ErrorToken)
                }
	
	return 200, nil
}

func getOwnershipReponse(ownership models.Ownership) structs.OwnershipResponse {
	return structs.OwnershipResponse{
                        OwnershipUID: ownership.OwnershipUID,                                    
                        ItemBarcode: ownership.ItemBarcode,
                        CustomItemName: ownership.CustomItemName,
                        CustItemImg: ownership.CustItemImg,
                        OwnedCustDesc: ownership.OwnedCustDesc,
                        ItemLocation: ownership.ItemLocation,
                        ItemQR: ownership.ItemQR,
                        ItemTags: ownership.ItemTags,
                        ItemQuantity: ownership.ItemQuantity,
                        ItemCheckedOut: ownership.ItemCheckedOut,
                        ItemBorrower: ownership.ItemBorrower,} 	
}

func CheckQR(c *fiber.Ctx) error {
	// Parse request into data map
        var data map[string]string
        err := c.BodyParser(&data)

        // Error with JSON request
        if err != nil {
                return returnError(c, 400, messages.ErrorParsingRequest)
        }

        uid := data["uid"]
        qr := c.Query("qr")

	// Validate Token
	code, err := validateToken(c, data["uid"], data["token"])	
	if err != nil {
		return returnError(c, code, err.Error())
	}

	if qr == "" {
		return returnError(c, 400, messages.QRMissing)
	}
  
        // Check if qr exists as location
        var location models.Location
        result := db.DB.Where("location_qr = ? AND location_owner = ?", qr, uid).First(&location)

	fmt.Println(location.LocationUID)

	if location.LocationUID != 0 {
		return c.Status(200).JSON(
			fiber.Map{
				"success":true,
 				"message":messages.Location}) 
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return returnError(c, 400, messages.ErrorWithConnection)
	}

	// Check if qr exists as ownership
	var ownership models.Ownership
	result = db.DB.Where("item_qr = ? AND item_owner = ?", qr, uid).First(&ownership)
	
	fmt.Println(ownership.OwnershipUID)

	if ownership.OwnershipUID != 0 {
		return c.Status(200).JSON(
			fiber.Map{
				"success":true,
				"message":messages.Ownership})
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return returnError(c, 400, messages.ErrorWithConnection)
	}

	return c.Status(200).JSON(
		fiber.Map{
			"success":true,
			"message":messages.New})
}
