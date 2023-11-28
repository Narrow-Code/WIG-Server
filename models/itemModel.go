package models

// Item represents information about an item.
type Item struct {
	ItemUid uint   `json:"itemUID" gorm:"primary_key;column:item_uid"`
	Barcode string `json:"barcode" gorm:"type:varchar(255);column:barcode"`
	Name    string `json:"itemName" gorm:"column:item_name"`
	Brand   string `json:"itemBrand" gorm:"column:item_brand"`
	Image   string `json:"itemImage" gorm:"column:item_img"`
}
