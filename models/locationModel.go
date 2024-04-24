package models

import "github.com/google/uuid"

// Location represents information about a location.
type Location struct {
	// LocationUID uniquely identifies the location
	LocationUID uuid.UUID `json:"locationUID" gorm:"primary_key;column:location_uid;type:varchar(191)"`

	// LocationOwner is the UUID of the user who owns the location
	LocationOwner uuid.UUID `json:"locationOwner" gorm:"column:location_owner;type:varchar(191)"`

	// LocationName is the name of the location
	LocationName string `json:"locationName" gorm:"column:location_name"`
	
	// Parent is the parent location in which the location is located in
	Parent uuid.UUID `json:"locationParent" gorm:"column:location_parent;type:varchar(191)"`
	
	//LocationQR is the QR code associated with the location
	LocationQR string `json:"locationQR" gorm:"column:location_qr"`
	
	// LocationTags are all of the tags representing the location, seperated by commas
	LocationTags string `json:"locationTags" gorm:"column:location_tags"`
	
	// LocationDescription is a custom description for the location
	LocationDescription string `json:"locationDescription" gorm:"column:location_description"`
	
	// User is the user account associated with the location
	User User `json:"user" gorm:"foreignkey:location_owner"`
	
	// Location is the parent location associated with the location
	Location *Location `json:"location" gorm:"foreignkey:location_parent"`
}
