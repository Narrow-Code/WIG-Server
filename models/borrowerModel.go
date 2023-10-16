/*
* Package models defines the data models used in the WIG-Server application.
*/
package models

/*
* Borrower represents information about a borrower.
*/
type Borrower struct {
	BorrowerUID 	uint 	`json:"borrower_uid" gorm:"primary_key;column:borrower_uid"`
	BorrowerName	string 	`json:"borrower_name" gorm:"column:borrower_name"`
}
