package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"

	"github.com/google/uuid"
)

/*
* createBorrower creates a models.Borrower and adds it to the database
*
* borrowerName the name of the Borrower to create
* user the User creating the Borrower
 */
func createBorrower(borrowerName string, user models.User) models.Borrower {
	utils.Log("building borrower")
	borrower := models.Borrower{
		BorrowerName:  borrowerName,
		BorrowerOwner: user.UserUID,
		BorrowerUID:   uuid.New(),
	}

	db.DB.Create(&borrower)
	utils.Log(borrower.BorrowerName + " successfully created")
	return borrower
}

/*
* checkout takes a list of ownership UUID's and checks them out to a single borrower
*
* @param ownerships the list of ownerships UID's to be checked out
* @param borrowerUUID the UUID of the Borrower who the Ownerships are being checked out to
* @return []string list of successfully checked out Ownership UID's
 */
func checkout(ownerships []string, borrowerUUID uuid.UUID) []string {
	utils.Log("began call")
	var successfulOwnerships []string
	for _, ownership := range ownerships {
		utils.Log("checking out " + ownership)
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)

		_, err := recordExists(result)
		if err == nil {
			item.ItemBorrower = borrowerUUID
			item.ItemCheckedOut = "true"
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
			utils.Log(item.CustomItemName + " checked out")
		}
	}
	utils.Log("success")
	return successfulOwnerships
}

/*
* checkin takes a list of ownership UUID's and returns them to the original location
*
* @param ownerships the list of Ownership UID's in which to return to original location
* @return []string list of successfully checked in Ownership UID's
 */
func checkin(ownerships []string) []string {
	utils.Log("began call")
	var successfulOwnerships []string
	for _, ownership := range ownerships {
		utils.Log("checking out " + ownership)
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)

		_, err := recordExists(result)
		if err == nil {
			item.ItemBorrower = uuid.MustParse(db.DefaultBorrowerUUID)
			item.ItemCheckedOut = "false"
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
			utils.Log(item.CustomItemName + " checked in")
		}
	}
	utils.Log("success")
	return successfulOwnerships
}

/*
* getBorrowerInventory returns a BorrowerInventory model with all borrowed Ownerships
*
* @param list of Borrowers to retrieve checked out items for
* @return []models.CheckedOutDTO the CheckedOutDTO list of inventory
 */
func getBorrowerInventory(borrowers []models.Borrower) []models.BorrowerInventory {
	var ownerships []models.Ownership
	var inventory []models.BorrowerInventory

	for b := range borrowers {
		query := db.DB.Where("item_borrower = ?", borrowers[b].BorrowerUID)

		if err := query.Find(&ownerships).Error; err != nil {
			continue
		}
		for o := range ownerships {
			preloadOwnership(&ownerships[o])
		}
		borrower := borrowerInventory(borrowers[b], ownerships)
		if len(ownerships) != 0 {
			inventory = append(inventory, borrower)
		}
	}
	return inventory
}

/*
* borrowerInventory creates a models.BorrowerInventory
*
* @param borrower The Borrower to create CheckedOutDTO for
* @param ownerships All of the Ownerships to add with Borrower
* @return models.CheckedOutDTO The CheckedOutDTO
 */
func borrowerInventory(borrower models.Borrower, ownerships []models.Ownership) models.BorrowerInventory {
	return models.BorrowerInventory{Borrower: borrower, Ownerships: ownerships}
}
