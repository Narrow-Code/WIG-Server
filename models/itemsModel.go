package Models

type Item struct {
	ItemUID 	string 	'json:"item_uid"'
	ItemBarcode 	string 	'json:"item_barcode"'
	ItemName 	string	'json:"item_name"'
	ItemBrand 	string 	'json:"item_brand"'
	ItemImg 	string 	'json:"item_img"'
	ItemDesc 	string 	'json:"item_description"'
}
