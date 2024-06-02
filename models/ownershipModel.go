package models

import(
	"github.com/google/uuid"
)

// Ownership represents information about ownership.
type Ownership struct {
	// OwnershipUID uniquely identifies the ownership
	OwnershipUID uuid.UUID `json:"ownershipUID" gorm:"primary_key;column:ownership_uid;type:varchar(191)"`

	// ItemOwner is the UUID of the user who owns the ownership  
	ItemOwner uuid.UUID `json:"itemOwner" gorm:"column:item_owner;type:varchar(191)"`

	// ItemNumber is the UUID of the item
	ItemNumber uuid.UUID `json:"itemNumber" gorm:"column:item_number;type:varchar(191)"`

	// CustomItemName is the users custom name for the ownership
	CustomItemName string `json:"customItemName" gorm:"column:custom_item_name"`

	// CustItemImg is a link to the image to be used to represent the ownership item
	CustItemImg string `json:"customItemImage" gorm:"column:custom_item_img"`

	// OwnedCustDesc is the users custom description for the ownership
	OwnedCustDesc string `json:"customItemDescription" gorm:"column:custom_item_description"`
	
	// ItemLocation is the UUID of the location in which the ownership is stored
	ItemLocation uuid.UUID `json:"itemLocation" gorm:"column:item_location;type:varchar(191)"`
	
	// ItemQR is the QR code representing the ownership
	ItemQR string `json:"itemQR" gorm:"column:item_qr"`
	
	// ItemTags are all of the tags representing the ownership, seperated by commas 
	ItemTags string `json:"itemTags" gorm:"column:item_tags"`
	
	// ItemQuantity is the inventory quantity of the ownership
	ItemQuantity int `json:"itemQuantity" gorm:"column:item_quantity;default:1"`
	
	// ItemCheckedOut is a boolean expression detailing if the item is checked out or borrowed
	ItemCheckedOut string `json:"itemCheckedOut" gorm:"column:item_checked_out;default:'false'"`
	
	// ItemBorrower is the UUID of the borrower who the ownership is checked out to
	ItemBorrower uuid.UUID `json:"itemBorrower" gorm:"column:item_borrower;type:varchar(191)"`
	
	// User is the user account associated with the ownership
	User User `json:"user" gorm:"foreignkey:item_owner"`
	
	// Location is the location associated with the ownership
	Location Location `json:"location" gorm:"foreignkey:item_location"`
	
	// Item is the item associated with the ownership
	Item Item `json:"item" gorm:"foreignkey:item_number"`
	
	// Borrower is the borrower associated with the ownership
	Borrower Borrower `json:"borrower" gorm:"foreignkey:item_borrower"`
}
