package controller

import (
	"WIG-Server/db"
	"WIG-Server/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateBorrower creates a borrower and adds it to the database.
func CreateBorrower(c *fiber.Ctx) error {
	// Initialize variables
	var borrower models.Borrower
	user := c.Locals("user").(models.User)
	borrowerName := c.Query("borrower")

	// Check for empty fields
	if borrowerName == "" {
		return Error(c, 400, "The borrower field is empty")
	}

	// Validate borrowerName is not in use
	result := db.DB.Where("borrower_name = ? AND borrower_owner = ?", borrowerName, user.UserUID).First(&borrower)
	code, err := recordNotInUse("Borrower Name", result)
	if err != nil {
		return Error(c, code, err.Error())
	}
	
	// Create Borrower and return as DTO
	borrower = createBorrower(borrowerName, user)
	dto := DTO("borrower", borrower)
	return Success(c, "Borrower created", dto)
}

// CheckoutItems checks out the list of Ownerships to a specified Borrower 
func CheckoutItems(c *fiber.Ctx) error {
	// Initialize variables
	var borrower models.Borrower
	var ownerships []string
	borrowerUID := c.Query("borrowerUID")

	// Check if borrowerUID is of correct UUID format
	borrowerUUID, err := uuid.Parse(borrowerUID)
	if err != nil {
		Error(c, 400, "Borrower UUID not correct format")	
	}

	// Check that borrower exists
	result := db.DB.Where("borrower_uid = ?", borrowerUID).First(&borrower)
	code, err := RecordExists("Borrower UID", result)
	if err != nil{
		return Error(c, code, err.Error())
	}
		
	// Parse json body
	err = c.BodyParser(&ownerships)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Checkout items in list
	successfulOwnerships := checkoutItems(ownerships, borrowerUUID)

	// Check if ownerships were successful
	if len(successfulOwnerships) == 0 {
		return Error(c, 400, "Failed to checkout ownerships")
	}

	// Return as DTO
	dto := DTO("ownerships", successfulOwnerships)	
	return Success(c, "Ownerships checked out", dto)
}

// CheckinItems sets returns checked out items to original locations within the list.
func CheckinItem(c *fiber.Ctx) error {
	// Initialize variables
	var ownerships []string

	// Parse json body
	err := c.BodyParser(&ownerships)
	if err != nil {
		return Error(c, 400, "There was an error parsing JSON")
	}

	// Checkin items in list
	successfulOwnerships := checkinItems(ownerships)

	// Check if ownerships were successful
	if len(successfulOwnerships) == 0 {
		return Error(c, 400, "Failed to checkout ownerships")
	}

	// Return as DTO
	dto := DTO("ownerships", successfulOwnerships)	
	return Success(c, "Ownerships checked in", dto)
}

// GetBorrower returns all borrowers associated with user.
func GetBorrowers(c *fiber.Ctx) error{
	// Initialize variables
	user := c.Locals("user").(models.User)

	// Get borrowers	
	var borrowers []models.Borrower
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)

	// Check if borrowers is empty
	if len(borrowers) == 0 {
		return Success(c, "No borrowers found")
	}

	// Return as DTO
	dto := DTO("borrowers", &borrowers)
	return Success(c, "Borrowers returned", dto)
}

// CheckedOutInventory returns all checked out inventory
func CheckedOutInventory(c *fiber.Ctx) error{
	// Initialize variables
	user := c.Locals("user").(models.User)
	var borrowers []models.Borrower
	var self models.Borrower
	
	// Get all borrower associated with User and include Self	
	db.DB.Where("borrower_owner = ?", user.UserUID).Find(&borrowers)
	db.DB.Where("borrower_uid = ?", "22222222-2222-2222-2222-222222222222").First(&self)
	borrowers = append(borrowers, self)

	// Get checkedOutDTO and return as DTO
	checkedOutDTO := getCheckedOutDto(borrowers)
	dto := DTO("borrowers", checkedOutDTO)
	return Success(c, "Checked Out Items returned", dto)
}
