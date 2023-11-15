// structs package holds all response structs for api calls
package structs

type OwnershipResponse struct {
	OwnershipUID    uint    `json:"ownership_uid"`
        CustomItemName  string  `json:"custom_item_name"`
        CustItemImg     string  `json:"custom_item_img"`                                  
        OwnedCustDesc   string  `json:"custom_item_description"`
        ItemLocation    string  `json:"item_location"`
        ItemQR          string  `json:"item_qr"`
        ItemTags        string  `json:"item_tags"`
        ItemQuantity    int     `json:"item_quantity"`
        ItemCheckedOut  string  `json:"item_checked_out"`
        ItemBorrower    uint    `json:"item_borrower"`
}
