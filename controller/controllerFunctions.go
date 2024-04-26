// Provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/models"
	"WIG-Server/utils"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
* RecordExists checks a gorm.DB error message to see if a record existed in the database.
*
* @param result The gorm.DB result to be checked.
* @return int The HTTP error code to return
* @return error The error message, if there is one.
 */
func recordExists(result *gorm.DB) (int, error) {
	if result.Error == gorm.ErrRecordNotFound {
		return 404, fmt.Errorf("not found in the database")
	}
	if result.Error != nil {
		return 400, errors.New(result.Error.Error())
	}
	return 200, nil
}

/*
* Checks a gorm.DB error message to see if a record is not in use in the database.
*
* @param field A string representing the field that is getting checked.
* @param result The gorm.DB result to be checked.
* @return int The HTTP error code to return
* @return error The error message, if there is one.
 */
func recordNotInUse(result *gorm.DB) (int, error) {
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 400, errors.New(result.Error.Error())
	}
	if result.RowsAffected != 0 {
		return 400, errors.New("record is in use in the database")
	}
	return 200, nil
}

/*
* success returns a 200 success code, and a message through fiber to the application.
*
* @param c The fiber context containing the HTTP request and response objects.
* @param message The success message to return via fiber.
* @param dtos Any extra fields to be added to the response map.
* @return error The c.Status being returned via fiber.
 */
func success(c *fiber.Ctx, message string, dtos ...models.DTO) error {
	user := c.Locals("user").(models.User)
	responseMap := fiber.Map{
		"message": message,
		"success": true}

	for _, dto := range dtos {
		responseMap[dto.Name] = dto.Data
	}

	log.Printf("%s#%s: %v \n\n", user.Username, utils.CallerFunctionName(2), responseMap)
	return c.Status(200).JSON(responseMap)
}

/*
* Error returns the given error code, and message through fiber to the application.
*
* @param c The fiber context containing the HTTP request and response objects.
* @param code The error code to return via fiber.
* @param message The error message to return via fiber.
* @return error The c.Status being returned via fiber.
 */
func Error(c *fiber.Ctx, code int, message string) error {
	user := c.Locals("user").(models.User)
	log.Printf("%s#%s: %v \n\n", user.Username, utils.CallerFunctionName(2), fiber.Map{"message": message})
	return c.Status(code).JSON(fiber.Map{
		"message": message,
		"success": false})
}

/*
* DTO creates a DTO model to pass in a response map from a name and data interface.
*
* @param name The name of the field to add.
* @param data The data to pass in the response map.
* @return models.DTO The DTO model.
 */
func DTO(name string, data interface{}) models.DTO {
	return models.DTO{Name: name, Data: data}
}
