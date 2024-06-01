package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"

	"github.com/google/uuid"
)

/*
* preloadLocation preloads the Locations foreignkey structs
*
* @param location The location to preload.
 */
func preloadLocation(location *models.Location) {
	// Preload the current location
	utils.Log("preloading " + location.LocationName)
	db.DB.Preload("User").Preload("Location").Find(&location)

	// Recursively preload the parent's hierarchy
	if location.Location.LocationUID != uuid.MustParse(db.DefaultLocationUUID) {
		utils.Log("Locations UUID:")
		utils.Log(location.Location.LocationUID.String())
		utils.Log("DEFAULTS:")
		utils.Log(db.DefaultLocationUUID)
		preloadLocation(location.Location)
	}
}

/*
* GetAllFromlocation returns the ownerships and locations inside of a parent location.
*
* @param location The location
* @param user The user making the call
* @return []models.Ownership list of Ownerships contained in the location
* @return []models.Location list of Locations contained in the location
*/
func unpackLocation(location models.Location, user models.User) ([]models.Ownership, []models.Location) {
	// Initialize variables
	utils.Log("unpacking " + location.LocationName)
	var ownerships []models.Ownership
	var locations []models.Location

	// Search all locations and ownerships from location
	db.DB.Where("item_location = ? AND item_owner = ?", location.LocationUID, user.UserUID).Find(&ownerships)	
	db.DB.Where("location_parent = ? AND location_owner = ?", location.LocationUID, user.UserUID).Find(&locations)

	// Preload ownerships and locations, then return
	for i := range ownerships {
		preloadOwnership(&ownerships[i])
	}
	for i := range locations {
		preloadLocation(&locations[i])
	}
	return ownerships, locations
}

/*
* ReturnAllInventory returns the entire Inventory of a user.
*
* @param location The default location
* @param user The user getting the inventory
* @return models.InventoryDTO The DTO with all locations and ownerships
*/
func getInventoryDTO(location models.Location, user models.User) models.InventoryDTO {
	// Initialize variables
	utils.Log("getting inventoryDTO for " + location.LocationName)
	var inventoryDTO models.InventoryDTO
	var inventoryList []models.InventoryDTO

	// Unpack location
	ownerships, locations := unpackLocation(location, user)

	// Recursively append location hierarchy
	for i := range locations {
		inventoryList = append(inventoryList, getInventoryDTO(locations[i], user))
	}

	// Set build inventoryDTO and return
	inventoryDTO.Parent = location
	inventoryDTO.Ownerships = ownerships	
	inventoryDTO.Locations = inventoryList
	return inventoryDTO
}

/*
* createLocation creates a location and adds it to the database
*
* @param locationName the name of the location
* @param user the user associated with the location
* @param locationQR the QR code to associate with the location
*/
func createLocation(locationName string, user models.User, locationQR string) models.Location {
	// Build location
	utils.Log("building Location for " + locationName)
	location := models.Location{
		LocationName:  locationName,
		LocationOwner: user.UserUID,
		LocationQR:    locationQR,
		Parent:        uuid.MustParse(db.DefaultLocationUUID),
		LocationUID:   uuid.New(),
	}

	// Add location to database, preload and return
	db.DB.Create(&location)
	utils.Log(locationName + " added to database")
	preloadLocation(&location)
	return location
}
