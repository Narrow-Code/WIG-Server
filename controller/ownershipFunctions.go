package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"
	"log"

	"github.com/google/uuid"
)

/*
* createOwnership creates an ownership relationship between a user and an item.
*
* @param uid The users UID.
* @param item The item associate with ownership
* @param qr The qr code to associate with ownership
* @param customName User generated custom name for ownership
* @return models.Ownership The ownership model.
* @return error The error message, if there is one.
 */
func createOwnership(uid uuid.UUID, item models.Item, qr string, customName string) (models.Ownership, error) {
	// Give blank variables value
	if customName == "" {
		customName = item.Name
	}
	if item.Name == "" {
		item.ItemUid = uuid.MustParse(db.DefaultItemUUID)
	}

	utils.Log("building ownership for " + customName)
	// Build ownership
	ownership := models.Ownership{
		OwnershipUID: uuid.New(),
		ItemOwner:  uid,
		ItemNumber: item.ItemUid,
		ItemLocation: uuid.MustParse(db.DefaultLocationUUID),
		ItemBorrower: uuid.MustParse(db.DefaultBorrowerUUID),
		ItemQR: qr,
		CustomItemName: customName,
	}

	// Create ownership in database and return
	result := db.DB.Create(&ownership)
	if result.Error != nil {
		utils.Log("Error creating ownership record: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		utils.Log("No rows were affected, creation may not have been successful")
	}
	utils.Log("Ownership record successfully created")
	return ownership, nil
}

/*
* preloadOwnership preloads the Ownerships foreignkey structs
*
* @param ownership The ownership to preload.
 */
func preloadOwnership(ownership *models.Ownership) {
	utils.Log("preloading: " + ownership.CustItemImg)
	db.DB.Preload("User").Preload("Item").Preload("Borrower").Preload("Location").Find(ownership)
	preloadLocation(&ownership.Location)
}

