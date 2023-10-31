/*
* Package models defines the data models used in the WIG-Server application.
*/
package models

/*
* Item represents information about an item.
*/
type Item struct {
	Barcode 	string 	`json:"barcode" gorm:"primary_key;type:varchar(255);column:barcode"`
	Name	 	string	`json:"item_name" gorm:"column:item_name"`
	Brand 		string 	`json:"item_brand" gorm:"column:item_brand"`
	Image	 	string 	`json:"item_img" gorm:"column:item_img"`
}
