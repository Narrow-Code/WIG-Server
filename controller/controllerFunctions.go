// Provides functions for handling HTTP requests and implementing business logic between the database and application.
package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
* Checks a gorm.DB error message to see if a record existed in the database.
*
* @param field A string representing the field that is getting checked.
* @param result The gorm.DB result to be checked.
*
* @return int The HTTP error code to return
* @return error The error message, if there is one.
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
* Checks a gorm.DB error message to see if a record is in use in the database.
*
* @param field A string representing the field that is getting checked.
* @param result The gorm.DB result to be checked.
*
* @return int The HTTP error code to return
* @return error The error message, if there is one.
 */
func recordNotInUse(field string, result *gorm.DB) (int, error) {
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return 400, errors.New(result.Error.Error())
	}
	if result.RowsAffected != 0 {
		return 400, errors.New(field + " Record is in use in the database")
	}
	log.Printf("controller#recordNotInUse: %s is not in use in the database", field)
	return 200, nil
}

/*
* Creates an ownership relationship between a user and an item.
*
* @param uid The users UID.
* @param itemUid The items UID.
*
* @return models.Ownership The ownership model.
* @return error The error message, if there is one.
 */
func createOwnership(uid uint, item models.Item, qr string, customName string) (models.Ownership, error) {
	if customName == "" {
		customName = item.Name
	}

	ownership := models.Ownership{
		ItemOwner:  uid,
		ItemNumber: item.ItemUid,
		ItemQuantity: 1,
		ItemQR: qr,
		CustomItemName: customName,
		ItemLocation: 1,
	}

	result := db.DB.Create(&ownership)
	if result.Error != nil {
		log.Printf("controller#createOwnership: Error creating ownership record: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("controller#CreateOwnership: No rows were affected, creation may not have been successful")
	}

	log.Printf("controller#createOwnership: Ownership record successfully created between user %d and item %d", uid, item.ItemUid)
	return ownership, nil
}

/*
* Returns a 200 success code, and a message through fiber to the application.
*
* @param c The fiber context containing the HTTP request and response objects.
* @param message The success message to return via fiber.
* @param dtos Any extra fields to be added to the response map.
*
* @return error The c.Status being returned via fiber.
 */
func Success(c *fiber.Ctx, message string, dtos ...models.DTO) error {
	responseMap := fiber.Map{
		"message": message,
		"success": true,}

	for _, dto := range dtos {
		responseMap[dto.Name] = dto.Data
	}

	log.Printf("%s: Status Code: 200, Response: %v", utils.CallerFunctionName(2), responseMap)
	return c.Status(200).JSON(responseMap)
}

/*
* Returns the given error code, and message through fiber to the application.
*
* @param c The fiber context containing the HTTP request and response objects.
* @param code The error code to return via fiber.
* @param message The error message to return via fiber.
*
* @return error The c.Status being returned via fiber.
 */
func Error(c *fiber.Ctx, code int, message string) error {
	log.Printf("%s: Status Code: %d, Response: %v", utils.CallerFunctionName(2), code, fiber.Map{"message": message})
	return c.Status(code).JSON(fiber.Map{
		"message": message,
		"success": false})
}

/*
* Creates a DTO model to pass in a response map from a name and data interface.
*
* @param name The name of the field to add.
* @param data The data to pass in the response map.
*
* @return models.DTO The DTO model.
 */
func DTO(name string, data interface{}) models.DTO {
	return models.DTO{Name: name, Data: data}
}

func CheckedOutDto(borrower models.Borrower, ownerships []models.Ownership) models.CheckedOutDTO {
	return models.CheckedOutDTO{Borrower: borrower, Ownerships: ownerships}
}

/*
* Preloads the Ownerships foreignkey structs
*
* @param ownership The ownership to preload.
 */
func preloadOwnership(ownership *models.Ownership) {
	db.DB.Preload("User").Preload("Item").Preload("Borrower").Preload("Location").Find(ownership)
	preloadLocation(&ownership.Location)
}

/*
* Preloads the Locations foreignkey structs
*
* @param location The location to preload.
 */
func preloadLocation(location *models.Location) {
	db.DB.Preload("User").Preload("Location").Find(&location)

	// Recursively preload the parent's hierarchy
	if location.Parent != nil && location.Location.LocationUID != 1 {
		preloadLocation(location.Location)
	}
}

/*
* Returns the ownerships and locations inside of a parent location.
*
* @param location The location
* @param user The user making the call
*/
func GetAllFromLocation(location models.Location, user models.User) ([]models.Ownership, []models.Location) {
	// search and get all ownerships from location
	var ownerships []models.Ownership
	db.DB.Where("item_location = ? AND item_owner = ?", location.LocationUID, user.UserUID).Find(&ownerships)	

	// search and get all locations from parent location
	var locations []models.Location
	db.DB.Where("location_parent = ? AND location_owner = ?", location.LocationUID, user.UserUID).Find(&locations)

	for i := range ownerships {
		preloadOwnership(&ownerships[i])
	}

	for i := range locations {
		preloadLocation(&locations[i])
	}

	return ownerships, locations
}

func ReturnAllInventory(location models.Location, user models.User) models.InventoryDTO {
	var inventoryDTO models.InventoryDTO
	var inventoryList []models.InventoryDTO

	ownerships, locations := GetAllFromLocation(location, user)

	for i := range locations {
		inventoryList = append(inventoryList, ReturnAllInventory(locations[i], user))
	}
	
	inventoryDTO.Parent = location
	inventoryDTO.Ownerships = ownerships	
	inventoryDTO.Locations = inventoryList

	return inventoryDTO
}
