package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/google/uuid"
)

/*
* createBorrower creates a models.Borrower and adds it to the database
* 
* borrowerName the name of the Borrower to create
* user the User creating the Borrower
*/
func createBorrower(borrowerName string, user models.User) models.Borrower{
	borrower := models.Borrower{
		BorrowerName:  borrowerName,
		BorrowerOwner: user.UserUID,
		BorrowerUID: uuid.New(),
	}

	db.DB.Create(&borrower)

	return borrower
}

/* 
* checkoutItems takes a list of ownership UUID's and checks them out to a single borrower
* 
* @param ownerships the list of ownerships UID's to be checked out
* @param borrowerUUID the UUID of the Borrower who the Ownerships are being checked out to
* @return []string list of successfully checked out Ownership UID's
*/
func checkoutItems(ownerships []string, borrowerUUID uuid.UUID) []string{
	var successfulOwnerships []string
	for _, ownership := range ownerships {		
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = borrowerUUID
			item.ItemCheckedOut = "true"
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
		}
	}
	return successfulOwnerships
}

/* 
* checkinItems takes a list of ownership UUID's and returns them to the original location
*
* @param ownerships the list of Ownership UID's in which to return to original location
* @return []string list of successfully checked in Ownership UID's
*/
func checkinItems(ownerships []string) []string {
	var successfulOwnerships []string
	for _, ownership := range ownerships{		
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = uuid.MustParse(db.DefaultBorrowerUUID)
			item.ItemCheckedOut = "false"
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
		}
	}
	return successfulOwnerships
}

/* getCheckedOutDto returns a CheckedOutDTO model with all borrowed Ownerships
*
* @param list of Borrowers to retrieve checked out items for
* @return []models.CheckedOutDTO the CheckedOutDTO list of inventory
*/
func getCheckedOutDto(borrowers []models.Borrower) []models.CheckedOutDTO {
	var ownerships []models.Ownership
	var checkedOutDTO []models.CheckedOutDTO

	for b := range borrowers{
		query := db.DB.Where("item_borrower = ?", borrowers[b].BorrowerUID)
		
		if err := query.Find(&ownerships).Error; err != nil{
			continue
		}	
		for o := range ownerships {
			preloadOwnership(&ownerships[o])
		}
		borrower := CheckedOutDto(borrowers[b], ownerships)
		if len(ownerships) != 0 {
			checkedOutDTO = append(checkedOutDTO, borrower)
		}
	}
	return checkedOutDTO
}
