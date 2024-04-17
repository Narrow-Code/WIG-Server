package models

// Represents a data transmission object to add to response maps.
type DTO struct {
	Name string
	Data interface{}
}

type CheckedOutDTO struct {
	Borrower Borrower `json:"borrower"`
	Ownerships []Ownership `json:"ownerships"`
}

type InventoryDTO struct {
	Parent Location `json:"parent"`
	Locations []InventoryDTO `json:"locations"`
	Ownerships []Ownership `json:"ownerships"`
}
