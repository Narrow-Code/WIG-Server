package models

import "github.com/google/uuid"

// Represents information about User profiles.
type User struct {
	UserUID      uuid.UUID   `json:"userUID" gorm:"primary_key;column:user_uid;type:varchar(191)"`
	Username     string `json:"username" gorm:"column:username"`
	Email        string `json:"email" gorm:"column:email"`
	Salt         string `json:"-" gorm:"column:salt"`
	Hash         string `json:"-" gorm:"column:hash"`
	EmailConfirm string `json:"emailConfirmed" gorm:"column:email_confirm;default:false"`
	Tier	     string `json:"tier" gorm:"column:tier"`
	Token        string `json:"-" gorm:"column:token"`
}
