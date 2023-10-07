package models

type User struct {
	UserUID		uint	`json:"uid" gorm:"primary_key;column:user_uid"` 
	Username 	string 	`json:"username" gorm"column:username"`
	UserEmail 	string 	`json:"email" gorm"column:email"`
	UserSalt 	string	`json:"salt" gorm"column:salt"`
	UserHash 	string 	`json:"hash" gorm"column:hash"`
	EmailConfirm 	string 	`json:"email_confirmed" gorm"column:email_confirm"`
	TierLevel 	string 	`json:"tier_level" gorm"column:tier_level"`
}
