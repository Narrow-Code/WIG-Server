package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"
	"WIG-Server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// BorrowerCreate creates a borrower and adds it to the database.
func BorrowerCreate(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var borrower models.Borrower
	var data map[string]string
	user := c.Locals("user").(models.User)

	// Parse JSON body
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	borrowerName := data["borrowerName"]


	// Check for empty fields
	utils.UserLog(c, "checking for empty fields")
	if borrowerName == "" {
		return Error(c, 400, "The borrower field is empty")
	}

	// Validate borrowerName is not in use
	utils.UserLog(c, "validating borrowerName is not in use")
	result := db.DB.Where("borrower_name = ? AND borrower_owner = ?", borrowerName, user.UserUID).First(&borrower)
	code, err := recordNotInUse(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Create Borrower and return as DTO
	borrower = createBorrower(borrowerName, user)
	dto := DTO("borrower", borrower)
	utils.UserLog(c, "success")
	return success(c, "Borrower created", dto)
}

// BorrowerCheckout checks out the list of Ownerships to a specified Borrower
func BorrowerCheckout(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var borrower models.Borrower
	var ownerships models.BorrowerRequest
	borrowerUID := c.Params("borrowerUID")

	// Check if borrowerUID is of correct UUID format
	utils.UserLog(c, "validating UUID format")
	borrowerUUID, err := uuid.Parse(borrowerUID)
	if err != nil {
		Error(c, 400, "Borrower UUID not correct format")
	}

	// Check that borrower exists
	utils.UserLog(c, "validating borrower exists")
	result := db.DB.Where("borrower_uid = ?", borrowerUID).First(&borrower)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Parse json body
	utils.UserLog(c, "parsing json body")
	err = c.BodyParser(&ownerships)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Checkout items in list
	utils.UserLog(c, "checking out items in list")
	successfulOwnerships := checkout(ownerships.Ownerships, borrowerUUID)

	// Check if ownerships were successful
	utils.UserLog(c, "checking for successful ownerships")
	if len(successfulOwnerships) == 0 {
		return Error(c, 400, "Failed to checkout ownerships")
	}

	// Return as DTO
	dto := DTO("ownerships", successfulOwnerships)
	utils.UserLog(c, "success")
	return success(c, "Ownerships checked out", dto)
}

// CheckinItems sets returns checked out items to original locations within the list.
func BorrowerCheckin(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var ownerships models.BorrowerRequest

	// Parse json body
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&ownerships)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Checkin items in list
	utils.UserLog(c, "checking in items in list")
	successfulOwnerships := checkin(ownerships.Ownerships)

	// Check if ownerships were successful
	utils.UserLog(c, "checking if ownerships were successful")
	if len(successfulOwnerships) == 0 {
		return Error(c, 400, "Failed to check in ownerships")
	}

	// Return as DTO
	dto := DTO("ownerships", successfulOwnerships)
	utils.UserLog(c, "success")
	return success(c, "Ownerships checked in", dto)
}

// GetBorrower returns all borrowers associated with user.
func BorrowerGetAll(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	user := c.Locals("user").(models.User)

	// Get borrowers
	utils.UserLog(c, "query for borrowers in database")
	var borrowers []models.Borrower
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)

	// Check if borrowers is empty
	utils.UserLog(c, "checking if borrowers were found")
	if len(borrowers) == 0 {
		return success(c, "No borrowers found")
	}

	// Return as DTO
	dto := DTO("borrowers", &borrowers)
	utils.UserLog(c, "success")
	return success(c, "Borrowers returned", dto)
}

// BorrowerGetInventory returns all checked out inventory
func BorrowerGetInventory(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	user := c.Locals("user").(models.User)
	var borrowers []models.Borrower
	var self models.Borrower

	// Get all borrower associated with User and include Self
	utils.UserLog(c, "getting all borrowers from database")
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)
	db.DB.Where("borrower_uid = ?", db.SelfBorrowerUUID).First(&self)
	borrowers = append(borrowers, self)

	// Get checkedOutDTO and return as DTO
	checkedOutDTO := getBorrowerInventory(borrowers)
	dto := DTO("borrowers", checkedOutDTO)
	utils.UserLog(c, "success")
	return success(c, "Checked Out Items returned", dto)
}

// BorrowerDelete deletes a borrower from the database and returns all Ownerships associated to them
func BorrowerDelete(c *fiber.Ctx) error {
	// Initialize variables
	utils.UserLog(c, "began call")
	var borrower models.Borrower
	var ownerships []models.Ownership
	var data map[string]string
	user := c.Locals("user").(models.User)

	// Parse request into data map
	utils.UserLog(c, "parsing json body")
	err := c.BodyParser(&data)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}
	borrowerUID := data["borrowerUID"]

	// Get borrower
	utils.UserLog(c, "rerieving borrower")
	result := db.DB.Where("borrower_uid = ? AND borrower_owner = ?", borrowerUID, user.UserUID).First(&borrower)
	code, err := recordExists(result)
	if err != nil {
		return Error(c, code, err.Error())
	}

	// Return ownerships to Default Borrower and checked in
	utils.UserLog(c, "Returning Ownerships to Default borrower")
	db.DB.Where("item_owner = ? AND item_borrower = ?", user.UserUID, borrower.BorrowerUID).Find(&ownerships)
	for ownership := range ownerships {
		ownerships[ownership].ItemBorrower = uuid.MustParse(db.DefaultBorrowerUUID)
		ownerships[ownership].ItemCheckedOut = "false"
	}
	db.DB.Save(&ownerships)


	utils.UserLog(c, "Deleting Borrower")
	db.DB.Delete(&borrower)
	if result := db.DB.Delete(&borrower); result.Error != nil {
		return Error(c, 500, "There was an error deleting the Borrower")
	}
	utils.UserLog(c, "success")

	return success(c, "Borrower deleted successfully")
}
