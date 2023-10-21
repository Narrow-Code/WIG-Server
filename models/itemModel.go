/*
* Package models defines the data models used in the WIG-Server application.
*/
package models

/*
* Item represents information about an item.
*/
type Item struct {
	ItemUID 	uint 	`json:"item_uid" gorm:"primary_key;column:item_uid"`
	ItemBarcode 	string 	`json:"item_barcode" gorm:"column:item_barcode"`
	ItemName 	string	`json:"item_name" gorm:"column:item_name"`
	ItemBrand 	string 	`json:"item_brand" gorm:"column:item_brand"`
	ItemImg 	string 	`json:"item_img" gorm:"column:item_img"`
	ItemDesc 	string 	`json:"item_description" gorm:"column:item_description"`
}
