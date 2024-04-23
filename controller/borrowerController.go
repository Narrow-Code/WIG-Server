package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/**
* Creates a borrower and adds it to the database.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
 */
func CreateBorrower(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	borrowerName := c.Query("borrower")

	// Check for empty fields
	if borrowerName == "" {return Error(c, 400, "The borrower field is empty")}

	// Validate location QR code is not in use
	var borrower models.Borrower
	result := db.DB.Where("borrower_name = ? AND borrower_owner = ?", borrowerName, user.UserUID).First(&borrower)
	code, err := recordNotInUse("Borrower Name", result)
	if err != nil {return Error(c, code, err.Error())}

	// create location
	borrower = models.Borrower{
		BorrowerName:  borrowerName,
		BorrowerOwner: user.UserUID,
		BorrowerUID: uuid.New(),
	}

	db.DB.Create(&borrower)

	borrowerDTO := DTO("borrower", borrower)

	return Success(c, "Borrower created", borrowerDTO)
}

/*
* Sets a list of ownerships to be checked out to a borrower.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func CheckoutItem(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	var request BorrowerRequest
	borrowerUID := c.Query("borrowerUID")
	borrowerUUID, err := uuid.Parse(borrowerUID)
	if err != nil {
		Error(c, 400, "Borrower UUID not correct format")	
	}
	var borrower models.Borrower
	result := db.DB.Where("borrower_uid = ?", borrowerUID).First(&borrower)

	err = c.BodyParser(&request)
	if err != nil {return Error(c, 400, "There was an error parsing JSON")}

	code, err := RecordExists("Borrower UID", result)
	if err != nil {return Error(c, code, err.Error())}

	success := 0
	var successfulOwnerships []string

	for _, ownership := range request.Ownerships {		
		var item models.Ownership
		result = db.DB.Where("ownership_uid = ? AND item_owner = ?", ownership, user.UserUID).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = borrowerUUID
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
			success++
		}
	}

	if success == 0 {
		return Error(c, 400, "Failed to checkout ownerships")
	}

	ownershipsDTO := DTO("ownerships", successfulOwnerships)	
	return Success(c, "Checked out", ownershipsDTO)
}

/*
* Sets returns checked out items to original owners within the list.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/

func CheckinItem(c *fiber.Ctx) error {
	// Initialize variables
	user := c.Locals("user").(models.User)
	var request BorrowerRequest

	err := c.BodyParser(&request)
	if err != nil {return Error(c, 400, "There was an error parsing JSON")}

	success := 0
	var successfulOwnerships []string

	for _, ownership := range request.Ownerships {		
		var item models.Ownership
		result := db.DB.Where("ownership_uid = ? AND item_owner = ?", ownership, user.UserUID).First(&item)
		
		_, err := RecordExists("Ownership", result)
		if err == nil {
			item.ItemBorrower = uuid.MustParse("11111111-1111-1111-1111-111111111111")
			db.DB.Save(&item)
			preloadOwnership(&item)
			successfulOwnerships = append(successfulOwnerships, ownership)
			success++
		}
	}

	if success == 0 {
		return Error(c, 400, "Failed to checkout ownerships")
	}

	ownershipsDTO := DTO("ownerships", successfulOwnerships)	
	return Success(c, "Checked in", ownershipsDTO)

}

/*
* Returns all borrowers associated with user.
*
* @param c The Fiber context containing the HTTP request and response objects.
*
* @return error The error message, if there is any.
*/
func GetBorrowers(c *fiber.Ctx) error{
	// Initialize variables
	user := c.Locals("user").(models.User)

	// Get borrower	
	var borrowers []models.Borrower
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)

	if len(borrowers) == 0 {
		return Success(c, "No borrowers found")
	}

	borrowersDTO := DTO("borrowers", &borrowers)

	return Success(c, "Borrowers returned", borrowersDTO)
}

func GetCheckedOutItems(c *fiber.Ctx) error{
	log.Print("GetCheckedOutItems: Started")
	// Initialize variables
	user := c.Locals("user").(models.User)
	var ownerships []models.Ownership
	var checkedOut []models.CheckedOutDTO
	
	// Get borrower	
	var borrowers []models.Borrower
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)
	log.Print("Borrowers searched")
	log.Print(borrowers)
	var self models.Borrower
	db.DB.Where("borrower_uid = ?", "22222222-2222-2222-2222-222222222222").First(&self)

	borrowers = append(borrowers, self)

	if len(borrowers) == 0 {
		return Success(c, "No borrowers found")
	}

	for b := range borrowers{
		ownerships = nil
		query := db.DB.Where("item_owner = ? AND item_borrower = ?", user.UserUID, borrowers[b].BorrowerUID)
		
		if err := query.Find(&ownerships).Error; err != nil{
			continue
		}	
		for o := range ownerships {
			preloadOwnership(&ownerships[o])
		}
		borrower := CheckedOutDto(borrowers[b], ownerships)
		if len(ownerships) != 0 {
			checkedOut = append(checkedOut, borrower)
		}
	} 
	checkedOutItems := DTO("borrowers", checkedOut)

	return Success(c, "Checked Out Items returned", checkedOutItems)
}

type BorrowerRequest struct {
	Ownerships []string `json:"ownerships"` 
}
