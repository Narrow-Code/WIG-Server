package models

import(
	"github.com/google/uuid"
)

// Represents information about ownership.
type Ownership struct {
	OwnershipUID   uuid.UUID     `json:"ownershipUID" gorm:"primary_key;column:ownership_uid;type:varchar(191)"`
	ItemOwner      uuid.UUID     `json:"itemOwner" gorm:"column:item_owner;type:varchar(191)"`
	ItemNumber     uuid.UUID     `json:"itemNumber" gorm:"column:item_number;type:varchar(191)"`
	CustomItemName string   `json:"customItemName" gorm:"column:custom_item_name"`
	CustItemImg    string   `json:"customItemImage" gorm:"column:custom_item_img"`
	OwnedCustDesc  string   `json:"customItemDescription" gorm:"column:custom_item_description"`
	ItemLocation   uuid.UUID     `json:"itemLocation" gorm:"column:item_location;default:AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA;type:varchar(191)"`
	ItemQR         string   `json:"itemQR" gorm:"column:item_qr"`
	ItemTags       string   `json:"itemTags" gorm:"column:item_tags"`
	ItemQuantity   int      `json:"itemQuantity" gorm:"column:item_quantity;default:1"`
	ItemCheckedOut string   `json:"itemCheckedOut" gorm:"column:item_checked_out"`
	ItemBorrower   uuid.UUID     `json:"itemBorrower" gorm:"column:item_borrower;type:varchar(191);default:11111111-1111-1111-1111-111111111111"`
	User           User     `json:"user" gorm:"foreignkey:item_owner"`
	Location       Location `json:"location" gorm:"foreignkey:item_location"`
	Item           Item     `json:"item" gorm:"foreignkey:item_number"`
	Borrower       Borrower `json:"borrower" gorm:"foreignkey:item_borrower"`
}
