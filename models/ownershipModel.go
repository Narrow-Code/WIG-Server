package models

/*
* Ownership represents information about ownership.
 */
type Ownership struct {
	OwnershipUID   uint     `json:"ownership_uid" gorm:"primary_key;column:ownership_uid"`
	ItemOwner      uint     `json:"item_owner" gorm:"column:item_owner"`
	ItemNumber     uint     `json:"item_number" gorm:"column:item_number"`
	CustomItemName string   `json:"custom_item_name" gorm:"column:custom_item_name"`
	CustItemImg    string   `json:"custom_item_img" gorm:"column:custom_item_img"`
	OwnedCustDesc  string   `json:"custom_item_description" gorm:"column:custom_item_description"`
	ItemLocation   uint     `json:"item_location" gorm:"column:item_location;default:1"`
	ItemQR         string   `json:"item_qr" gorm:"column:item_qr"`
	ItemTags       string   `json:"item_tags" gorm:"column:item_tags"`
	ItemQuantity   int      `json:"item_quantity" gorm:"column:item_quantity;"`
	ItemCheckedOut string   `json:"item_checked_out" gorm:"column:item_checked_out"`
	ItemBorrower   uint     `json:"item_borrower" gorm:"column:item_borrower;default:1"`
	User           User     `gorm:"foreignkey:item_owner"`
	Location       Location `gorm:"foreignkey:item_location"`
	Item           Item     `gorm:"foreignkey:item_number"`
	Borrower       Borrower `gorm:"foreignkey:item_borrower"`
}
