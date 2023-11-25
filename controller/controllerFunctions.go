// Package controller provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/db"
	"WIG-Server/dto"
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
RecordExists checks a gorm.DB error message to see if a record existed in the database.

@param field A string representing the field that is getting checked.
@param result The gorm.DB result that may hold the error message.
@return int The HTTP error code to return
@return The error message to return
*/
func RecordExists(field string, result *gorm.DB) (int, error) {
	if result.Error == gorm.ErrRecordNotFound {
		return 404, fmt.Errorf("%s was not found in the database", field)
	}
	if result.Error != nil {
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
		return 400, errors.New(result.Error.Error())
	}
	if result.RowsAffected != 0 {
		return 400, errors.New(field + "Record is in use in the database")
	}
	log.Printf("controller#recordNotInUse: %s is not in use in the database", field)
	return 200, nil
}

func createOwnership(uid string, itemUid uint) (models.Ownership, error) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return models.Ownership{}, errors.New("There was an error converting uid to Int")
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

	log.Printf("%s: Status Code: 200, Response: %v", utils.CallerFunctionName(), responseMap)
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
	log.Printf("%s: Status Code: %d, Response: %v", utils.CallerFunctionName(), code, fiber.Map{"message": message})
	return c.Status(code).JSON(fiber.Map{
		"message": message})
}

func DTO(name string, data interface{}) dto.DTO {
	return dto.DTO{Name: name, Data: data}
}
