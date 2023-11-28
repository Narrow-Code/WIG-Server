package models

// Represents information about ownership.
type Ownership struct {
	OwnershipUID   uint     `json:"ownershipUID" gorm:"primary_key;column:ownership_uid"`
	ItemOwner      uint     `json:"itemOwner" gorm:"column:item_owner"`
	ItemNumber     uint     `json:"itemNumber" gorm:"column:item_number"`
	CustomItemName string   `json:"customItemName" gorm:"column:custom_item_name"`
	CustItemImg    string   `json:"customItemImage" gorm:"column:custom_item_img"`
	OwnedCustDesc  string   `json:"customItemDescription" gorm:"column:custom_item_description"`
	ItemLocation   uint     `json:"itemLocation" gorm:"column:item_location;default:1"`
	ItemQR         string   `json:"itemQR" gorm:"column:item_qr"`
	ItemTags       string   `json:"itemTags" gorm:"column:item_tags"`
	ItemQuantity   int      `json:"itemQuantity" gorm:"column:item_quantity;"`
	ItemCheckedOut string   `json:"itemCheckedOut" gorm:"column:item_checked_out"`
	ItemBorrower   uint     `json:"itemBorrower" gorm:"column:item_borrower;default:1"`
	User           User     `json:"user" gorm:"foreignkey:item_owner"`
	Location       Location `json:"location" gorm:"foreignkey:item_location"`
	Item           Item     `json:"item" gorm:"foreignkey:item_number"`
	Borrower       Borrower `json:"borrower" gorm:"foreignkey:item_borrower"`
}
