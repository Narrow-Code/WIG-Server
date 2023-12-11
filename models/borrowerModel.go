// Defines the data models used in the WIG-Server application.
package models

// Represents information about a borrower.
type Borrower struct {
	BorrowerUID  	uint    `json:"borrowerUID" gorm:"primary_key;column:borrower_uid"`
	BorrowerName 	string  `json:"borrowerName" gorm:"column:borrower_name"`
	BorrowerOwner	uint 	`json:"borrowerOwner" gorm:"column:borrower_owner"`
}
