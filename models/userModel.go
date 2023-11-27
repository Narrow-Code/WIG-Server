package models

// Represents information about User profiles.
type User struct {
	UserUID      uint   `json:"user_uid" gorm:"primary_key;column:user_uid"`
	Username     string `json:"username" gorm:"column:username"`
	Email        string `json:"email" gorm:"column:email"`
	Salt         string `json:"salt" gorm:"column:salt"`
	Hash         string `json:"hash" gorm:"column:hash"`
	EmailConfirm string `json:"email_confirmed" gorm:"column:email_confirm;default:false"`
	TierLevel    string `json:"tier_level" gorm:"column:tier_level"`
	Token        string `json:"-" gorm:"column:token"`
}
