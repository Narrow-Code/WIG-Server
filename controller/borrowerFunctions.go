package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/google/uuid"
)

// createBorrower creates a models.Borrower and adds it to the database
func createBorrower(borrowerName string, user models.User) models.Borrower{
	borrower := models.Borrower{
		BorrowerName:  borrowerName,
		BorrowerOwner: user.UserUID,
		BorrowerUID: uuid.New(),
	}

	db.DB.Create(&borrower)

	return borrower
}

// checkoutItems takes a list of ownership UUID's and checks them out to a single borrower
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

// checckinItems takes a list of ownership UUID's and returns them to the original location
func checkinItems(ownerships []string) []string {
	var successfulOwnerships []string
	for _, ownership := range ownerships{		
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = uuid.MustParse("11111111-1111-1111-1111-111111111111")
			item.ItemCheckedOut = "false"
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
		}
	}
	return successfulOwnerships
}

// getCheckedOutDto returns a CheckedOutDTO model with all borrowed Ownerships
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
