package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/google/uuid"
)

func createBorrower(borrowerName string, user models.User) models.Borrower{
	borrower := models.Borrower{
		BorrowerName:  borrowerName,
		BorrowerOwner: user.UserUID,
		BorrowerUID: uuid.New(),
	}

	db.DB.Create(&borrower)

	return borrower
}

func checkoutItems(ownerships []string, borrowerUUID uuid.UUID) []string{
	var successfulOwnerships []string
	for _, ownership := range ownerships {		
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ?", ownership).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = borrowerUUID
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
		}
	}
	return successfulOwnerships
}
