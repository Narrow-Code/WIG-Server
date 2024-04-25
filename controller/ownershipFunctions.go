package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
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
		log.Printf("controller#createOwnership: Error creating ownership record: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("controller#CreateOwnership: No rows were affected, creation may not have been successful")
	}
	log.Printf("controller#createOwnership: Ownership record successfully created between user %d and item %d", uid, item.ItemUid)
	return ownership, nil
}

/*
* preloadOwnership preloads the Ownerships foreignkey structs
*
* @param ownership The ownership to preload.
 */
func preloadOwnership(ownership *models.Ownership) {
	db.DB.Preload("User").Preload("Item").Preload("Borrower").Preload("Location").Find(ownership)
	preloadLocation(&ownership.Location)
}

