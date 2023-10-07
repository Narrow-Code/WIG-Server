package models

type Ownership struct {
	OwnedUID 	uint 	`json:"owned_uid" gorm:"primary_key;column:owned_uid"`
	OwnedOwner 	uint 	`json:"owned_owner" gorm="column:owned_owner"`
	OwnedItem 	uint	`json:"owned_item" gorm="column:owned_item"`
	OwnedCustName	string 	`json:"owned_custom_name" gorm="column:owned_custom_name"`
	OwnedCustImg	string	`json:"owned_custom_img" gorm="column:owned_custom_img"`
	OwnedCustDesc	string 	`json:"owned_custom_description" gorm="column:owned_custom_description"`
	OwnedLocation 	uint 	`json:"owned_location" gorm="column:owned_location"`
	OwnedQR 	string 	`json:"owned_qr" gorm="column:owned_qr"`
	OwnedTags	string 	`json:"owned_tags" gorm="column:owned_tags"`
	OwnedQuantity	int 	`json:"owned_quantity" gorm="column:owned_quantity"`
	OwnedCheckedOut	string	`json:"owned_checked_out" gorm="column:owned_checked_out"`
	OwnedBorrower	uint 	`json:"owned_borrower" gorm="column:owned_borrower"`
	User            User    `gorm:"foreignkey:owned_owner"`
	Location	Location `gorm:"foreignkey:owned_location"`
	Item		Item	`gorm:"foreignkey:owned_item"`
	Borrower	Borrower `gorm:"foreignkey:owned_borrower"`
}
