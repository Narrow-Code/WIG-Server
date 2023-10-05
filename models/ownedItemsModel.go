package Models

type OwnedItem struct {
	ownedUID 	string 	'json:"owned_uid"'
	OwnedOwner 	string 	'json:"owned_owner"'
	OwnedName 	string	'json:"owned_name"'
	OwnedCustName	string 	'json:"owned_custom_name"'
	OwnedCustImg	string	'json:"owned_custom_img"'
	OwnedCustDesc	string 	'json:"owned_custom_description"'
	OwnedLocation 	string 	'json:"owned_location"'
	OwnedQR 	string 	'json:"owned_qr"'
	OwnedTags	string 	'json:"owned_tags"'
	OwnedQuantity	int 	'json:"owned_quantity"'
	OwnedCheckedOut	string	'json:"owned_checked_out"'
	OwnedBorrower	string 	'json:"owned_borrower"'
}
