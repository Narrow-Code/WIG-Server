package models

import (
	"github.com/google/uuid"
)

// Item represents information about an item.
type Item struct {
	// ItemUid uniquely identifies the item
	ItemUid uuid.UUID `json:"itemUID" gorm:"primary_key;column:item_uid;type:varchar(191)"`

	// Barcode is the items barcode
	Barcode string `json:"barcode" gorm:"type:varchar(255);column:barcode"`

	// Name is the name of the item
	Name string `json:"itemName" gorm:"column:item_name"`

	// Brand is the items brand
	Brand string `json:"itemBrand" gorm:"column:item_brand"`

	// Image is a link to the image being used for the item
	Image string `json:"itemImage" gorm:"column:item_img"`
}
