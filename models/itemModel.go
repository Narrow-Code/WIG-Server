package models

// Item represents information about an item.
type Item struct {
	ItemUid uint   `json:"item_uid" gorm:"primary_key;column:item_uid"`
	Barcode string `json:"barcode" gorm:"type:varchar(255);column:barcode"`
	Name    string `json:"item_name" gorm:"column:item_name"`
	Brand   string `json:"item_brand" gorm:"column:item_brand"`
	Image   string `json:"item_img" gorm:"column:item_img"`
}
