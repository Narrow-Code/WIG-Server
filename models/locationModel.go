package models

import "github.com/google/uuid"

// Represents information about a location.
type Location struct {
	LocationUID         uuid.UUID      `json:"locationUID" gorm:"primary_key;column:location_uid;type:varchar(191)"`
	LocationOwner       uuid.UUID      `json:"locationOwner" gorm:"column:location_owner;type:varchar(191)"`
	LocationName        string    `json:"locationName" gorm:"column:location_name"`
	Parent              uuid.UUID     `json:"locationParent" gorm:"column:location_parent;type:varchar(191);default:AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA"`
	LocationQR          string    `json:"locationQR" gorm:"column:location_qr"`
	LocationTags        string    `json:"locationTags" gorm:"column:location_tags"`
	LocationDescription string    `json:"locationDescription" gorm:"column:location_description"`
	User                User      `json:"user" gorm:"foreignkey:location_owner"`
	Location            *Location `json:"location" gorm:"foreignkey:location_parent"`
}
