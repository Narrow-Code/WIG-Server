// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/utils"
	"WIG-Server/db"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/dto"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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
        if !utils.ValidateToken(user.Username, user.Hash, token) {
                return 400, errors.New(messages.ErrorToken)
                }
	
	return 200, nil
}

/*
RecordExists checks a gorm.DB error message to see if a record existed in the database.

@param field A string representing the field that is getting checked.
@param result The gorm.DB result that may hold the error message.
@return int The HTTP error code to return
@return The error message to return
*/
func recordExists(field string, result *gorm.DB) (int, error) {

	if result.Error == gorm.ErrRecordNotFound {return 404, errors.New(field + messages.DoesNotExist)}
	if result.Error != nil {return 400, errors.New(result.Error.Error())}
	
	return 200, nil
}

/*
RecordNotInUse checks a gorm.DB error message to see if a record is in use in the database.

@param field A string representing the field that is getting checked.
@param result The gorm.DB result that may hold the error message.
@return int The HTTP error code to return
@return The error message to return
*/
func recordNotInUse(field string, result *gorm.DB) (int, error) {	
	
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {return 400, errors.New(result.Error.Error())}
 	if result.RowsAffected != 0 {return 400, errors.New(field + messages.RecordInUse)}
	
	return 200, nil
}

func createOwnership(uid string, itemUid uint) (models.Ownership, error){
	// Convert uid to int
	uidInt, err := strconv.Atoi(uid)
	if err != nil {return models.Ownership{}, errors.New(messages.ConversionError)}
	
	ownership := models.Ownership{
               	ItemOwner:uint(uidInt),
		ItemNumber:itemUid,
   	}
		
	db.DB.Create(&ownership)

	return ownership, nil
}
