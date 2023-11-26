package models

// Represents information about a location.
type Location struct {
	LocationUID         uint      `json:"location_uid" gorm:"primary_key;column:location_uid"`
	LocationOwner       uint      `json:"location_owner" gorm:"column:location_owner"`
	LocationName        string    `json:"location_name" gorm:"column:location_name"`
	LocationType        string    `json:"location_type" gorm:"column:location_type"`
	LocationLocation    *uint     `json:"location_location" gorm:"column:location_location;default:1"`
	LocationQR          string    `json:"location_qr" gorm:"column:location_qr"`
	LocationTags        string    `json:"location_tags" gorm:"column:location_tags"`
	LocationDescription string    `json:"location_description" gorm:"column:location_description"`
	User                User      `gorm:"foreignkey:location_owner"`
	Location            *Location `gorm:"foreignkey:location_location"`
}
