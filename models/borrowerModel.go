// Defines the data models used in the WIG-Server application.
package models

import (
	"github.com/google/uuid"
)

// Represents information about a borrower.
type Borrower struct {
	BorrowerUID  	uuid.UUID    `json:"borrowerUID" gorm:"primary_key;column:borrower_uid;type:varchar(191)"`
	BorrowerName 	string  `json:"borrowerName" gorm:"column:borrower_name"`
	BorrowerOwner	uint 	`json:"-" gorm:"column:borrower_owner"`
}
