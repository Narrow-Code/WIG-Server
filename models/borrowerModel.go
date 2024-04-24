// models defines the data models used in the WIG-Server application.
package models

import (
	"github.com/google/uuid"
)

// Borrower represents information about a borrower.
type Borrower struct {
	// BorrowerUID uniquely identifies the borrower
	BorrowerUID uuid.UUID `json:"borrowerUID" gorm:"primary_key;column:borrower_uid;type:varchar(191)"`

	// BorrowerName is the name of the borrower
	BorrowerName string `json:"borrowerName" gorm:"column:borrower_name"`

	// BorrowerOwner is the unique identifier of the owner of the borrower.
	BorrowerOwner uuid.UUID `json:"-" gorm:"column:borrower_owner;type:varchar(191)"`
}
