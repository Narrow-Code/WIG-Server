// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/db"
	"WIG-Server/dto"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/utils"
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
func validateToken(c *fiber.Ctx, uid string, token string) (int, error) {
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
	if token == utils.GenerateToken(user.Username, user.Hash) {
		return 200, nil
	} else {
	return 400, errors.New(messages.ErrorToken)
	}

}

/*
RecordExists checks a gorm.DB error message to see if a record existed in the database.

@param field A string representing the field that is getting checked.
@param result The gorm.DB result that may hold the error message.
@return int The HTTP error code to return
@return The error message to return
*/
func recordExists(field string, result *gorm.DB) (int, error) {

	if result.Error == gorm.ErrRecordNotFound {
		return 404, errors.New(field + messages.DoesNotExist)
	}
	if result.Error != nil {
		return 400, errors.New(result.Error.Error())
	}

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

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 400, errors.New(result.Error.Error())
	}
	if result.RowsAffected != 0 {
		return 400, errors.New(field + messages.RecordInUse)
	}

	return 200, nil
}

func createOwnership(uid string, itemUid uint) (models.Ownership, error) {
	// Convert uid to int
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return models.Ownership{}, errors.New(messages.ConversionError)
	}

	ownership := models.Ownership{
		ItemOwner:  uint(uidInt),
		ItemNumber: itemUid,
	}

	db.DB.Create(&ownership)

	return ownership, nil
}

/*
returnSuccess returns a 200 success code, a 'true' success status and a message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param message The success message to return via fiber.
@return error - An error, if any, that occurred during the process.
*/
func Success(c *fiber.Ctx, message string, dtos ...dto.DTO) error {
	responseMap := fiber.Map{
		"message": message,
	}

	for _, dto := range dtos {
		responseMap[dto.Name] = dto.Data
	}

	return c.Status(200).JSON(responseMap)
}

/*
returnError returns the given error code, a 'false' success status and message through fiber to the application.

@param c The fiber context containing the HTTP request and response objects.
@param code The error code to return via fiber
@param message The error message to return via fiber
@return error - An error, if any, that occurred during the process.
*/
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message})
}

func DTO(name string, data interface{}) dto.DTO {
	return dto.DTO{Name: name, Data: data}
}
