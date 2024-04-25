package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/google/uuid"
)

/*
* preloadLocation preloads the Locations foreignkey structs
*
* @param location The location to preload.
 */
func preloadLocation(location *models.Location) {
	db.DB.Preload("User").Preload("Location").Find(&location)

	// Recursively preload the parent's hierarchy
	if location.Location.LocationUID != uuid.MustParse(db.DefaultLocationUUID) {
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

/*
* ReturnAllInventory returns the entire Inventory of a user.
*
* @param location The default location
* @param user The user getting the inventory
* @return models.InventoryDTO The DTO with all locations and ownerships
*/
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

func createLocation(locationName string, user models.User, locationQR string) models.Location {
	location := models.Location{
		LocationName:  locationName,
		LocationOwner: user.UserUID,
		LocationQR:    locationQR,
		Parent:        uuid.MustParse(db.DefaultLocationUUID),
		LocationUID:   uuid.New(),
	}

	db.DB.Create(&location)
	preloadLocation(&location)

	return location
}
