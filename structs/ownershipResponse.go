// structs package holds all response structs for api calls
package structs

type OwnershipResponse struct {
	OwnershipUID    uint    `json:"ownership_uid" gorm:"primary_key;column:ownership_uid"`
        ItemBarcode     string  `json:"item_barcode" gorm:"column:item_barcode"`
        CustomItemName  string  `json:"custom_item_name" gorm:"column:custom_item_name"`
        CustItemImg     string  `json:"custom_item_img" gorm:"column:custom_item_img"`                                  
        OwnedCustDesc   string  `json:"custom_item_description" gorm:"column:custom_item_description"`
        ItemLocation    uint    `json:"item_location" gorm:"column:item_location;"`
        ItemQR          string  `json:"item_qr" gorm:"column:item_qr"`
        ItemTags        string  `json:"item_tags" gorm:"column:item_tags"`
        ItemQuantity    int     `json:"item_quantity" gorm:"column:item_quantity"`
        ItemCheckedOut  string  `json:"item_checked_out" gorm:"column:item_checked_out"`
        ItemBorrower    uint    `json:"item_borrower" gorm:"column:item_borrower;"`
}