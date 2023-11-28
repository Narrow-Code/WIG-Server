package models

// Represents information about User profiles.
type User struct {
	UserUID      uint   `json:"userUID" gorm:"primary_key;column:user_uid"`
	Username     string `json:"username" gorm:"column:username"`
	Email        string `json:"email" gorm:"column:email"`
	Salt         string `json:"salt" gorm:"column:salt"`
	Hash         string `json:"hash" gorm:"column:hash"`
	EmailConfirm string `json:"emailConfirmed" gorm:"column:email_confirm;default:false"`
	Tier	     string `json:"tier" gorm:"column:tier"`
	Token        string `json:"-" gorm:"column:token"`
}
