/*
* Package models defines the data models used in the WIG-Server application.
*/
package models

/*
* Ownership represents information about ownership.
*/
type Ownership struct {
	OwnershipUID 	uint 	`json:"ownership_uid" gorm:"primary_key;column:ownership_uid"`
	ItemOwner 	uint 	`json:"item_owner" gorm="column:item_owner"`
	ItemBarcode 	uint	`json:"item_barcode" gorm="column:item_barcode"`
	OwnedCustName	string 	`json:"owned_custom_name" gorm="column:owned_custom_name"`
	OwnedCustImg	string	`json:"owned_custom_img" gorm="column:owned_custom_img"`
	OwnedCustDesc	string 	`json:"owned_custom_description" gorm="column:owned_custom_description"`
	OwnedLocation 	uint 	`json:"owned_location" gorm="column:owned_location"`
	OwnedQR 	string 	`json:"owned_qr" gorm="column:owned_qr"`
	OwnedTags	string 	`json:"owned_tags" gorm="column:owned_tags"`
	OwnedQuantity	int 	`json:"owned_quantity" gorm="column:owned_quantity"`
	OwnedCheckedOut	string	`json:"owned_checked_out" gorm="column:owned_checked_out"`
	OwnedBorrower	uint 	`json:"owned_borrower" gorm="column:owned_borrower"`
	User            User    `gorm:"foreignkey:item_owner"`
	Location	Location `gorm:"foreignkey:owned_location"`
	Item		Item	`gorm:"foreignkey:item_barcode"`
	Borrower	Borrower `gorm:"foreignkey:owned_borrower"`
}
