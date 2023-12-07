package models

// Represents information about a location.
type Location struct {
	LocationUID         uint      `json:"locationUID" gorm:"primary_key;column:location_uid"`
	LocationOwner       uint      `json:"locationOwner" gorm:"column:location_owner"`
	LocationName        string    `json:"locationName" gorm:"column:location_name"`
	Parent              *uint     `json:"locationParent" gorm:"column:location_parent;default:1"`
	LocationQR          string    `json:"locationQR" gorm:"column:location_qr"`
	LocationTags        string    `json:"locationTags" gorm:"column:location_tags"`
	LocationDescription string    `json:"locationDescription" gorm:"column:location_description"`
	User                User      `json:"user" gorm:"foreignkey:location_owner"`
	Location            *Location `json:"location" gorm:"foreignkey:location_parent"`
}
