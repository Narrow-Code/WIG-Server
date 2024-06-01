package models

// Represents a data transmission object to add to response maps.
type DTO struct {
	Name string
	Data interface{}
}

// BorrowerInventory is the data transmission object for getting Checked Out items
type BorrowerInventory struct {
	Borrower   Borrower    `json:"borrower"`
	Ownerships []Ownership `json:"ownerships"`
}

// InventoryDTO is the data transmission object for getting Inventory
type InventoryDTO struct {
	Parent     Location       `json:"parent"`
	Locations  []InventoryDTO `json:"locations"`
	Ownerships []Ownership    `json:"ownerships"`
}

type BorrowerRequest struct {
	Ownerships []int `json:"ownerships"` 
}
