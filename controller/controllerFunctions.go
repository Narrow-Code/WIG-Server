// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/db"
	"WIG-Server/dto"
	"WIG-Server/messages"
	"WIG-Server/models"
	"WIG-Server/utils"
	"errors"
	"fmt"
	"log"
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
	if uid == "" || token == "" {
		log.Println("controller#validateToken: UID or Token is empty, returning error")
		return 400, errors.New("UID or Token is empty")
	}

	log.Printf("controller#validateToken: Retrieving user from the database for UID: %s", uid)
	var user models.User
	result := db.DB.Where("user_uid = ?", uid).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("controller#validateToken: UID %s was not found in the database, returning error", uid)
		return 404, errors.New(fmt.Sprintf("UID %s was not found in the database", uid))

	} else if result.Error != nil {
		log.Printf("controller#validateToken: There was an error with the connection, returning error %v", result.Error)
		return 500, errors.New("Internal server error")
	}

	if token == utils.GenerateToken(user.Username, user.Hash) {
		log.Printf("controller#ValidateToken: UID %s has been authenticated", uid)
		return 200, nil
	} else {
		log.Printf("controller#ValidateToken: UID %s is unauthorized by token", uid)
		return 401, errors.New(fmt.Sprintf("UID %s is unauthorized by token", uid))
	}
}

/*
RecordExists checks a gorm.DB error message to see if a record existed in the database.

@param field A string representing the field that is getting checked.
@param result The gorm.DB result that may hold the error message.
@return int The HTTP error code to return
@return The error message to return
*/
func RecordExists(field string, result *gorm.DB) (int, error) {
	if result.Error == gorm.ErrRecordNotFound {
		log.Printf("controller#recordExists: %s was not found in the database, returning error", field)
		return 404, errors.New(fmt.Sprintf("%s was not found in the database", field))
	}
	if result.Error != nil {
		log.Printf("controller#recordExists: Error checking %s in the database: %v", field, result.Error.Error())
		return 400, errors.New(result.Error.Error())
	}
	log.Printf("controller#recordExists: %s was successfully found in the database", field)
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
		log.Printf("controller#recordNotInUse: Error checking %s usage in the database: %v", field, result.Error)
		return 400, errors.New(result.Error.Error())
	}
	if result.RowsAffected != 0 {
		log.Printf("controller#recordNotInUse: %s is in use in the database, returning error", field)
		return 400, errors.New(field + messages.RecordInUse)
	}
	log.Printf("controller#recordNotInUse: %s is not in use in the database", field)
	return 200, nil
}

func createOwnership(uid string, itemUid uint) (models.Ownership, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		log.Printf("controller#createOwnership: Error converting UID %s to uint", uid)
		return models.Ownership{}, errors.New(messages.ConversionError)
	}

	ownership := models.Ownership{
		ItemOwner:  uint(uidInt),
		ItemNumber: itemUid,
	}

	result := db.DB.Create(&ownership)
	if result.Error != nil {
		log.Printf("controller#createOwnership: Error creating ownership record: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("controller#CreateOwnership: No rows were affected, creation may not have been successful")
	}

	log.Printf("controller#createOwnership: Ownership record successfully created between user %s and item %d", uid, itemUid)
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
	
	log.Printf("controller#Success: Status Code: 200, Response: %v", responseMap)
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
	log.Printf("controller#Error: Status Code: %d, Response: %v", code, fiber.Map{"message": message})
	return c.Status(code).JSON(fiber.Map{
		"message": message})
}

func DTO(name string, data interface{}) dto.DTO {
	log.Printf("controller#DTO: DTO was created for %s", name)
	return dto.DTO{Name: name, Data: data}
}
