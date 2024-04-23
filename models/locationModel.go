package models

import "github.com/google/uuid"

// Represents information about a location.
type Location struct {
	LocationUID         uuid.UUID      `json:"locationUID" gorm:"primary_key;column:location_uid;type:varchar(191)"`
	LocationOwner       uint      `json:"locationOwner" gorm:"column:location_owner"`
	LocationName        string    `json:"locationName" gorm:"column:location_name"`
	Parent              *uuid.UUID     `json:"locationParent" gorm:"column:location_parent;type:varchar(191);default:44444444-4444-4444-4444-444444444444"`
	LocationQR          string    `json:"locationQR" gorm:"column:location_qr"`
	LocationTags        string    `json:"locationTags" gorm:"column:location_tags"`
	LocationDescription string    `json:"locationDescription" gorm:"column:location_description"`
	User                User      `json:"user" gorm:"foreignkey:location_owner"`
	Location            *Location `json:"location" gorm:"foreignkey:location_parent"`
}
