package models

type Borrower struct {
	BorrowerUID 	uint 	`json:"borrower_uid" gorm:"primary_key;column:borrower_uid"`
	BorrowerName	string 	`json:"borrower_name" gorm:"column:borrower_name"`
}
