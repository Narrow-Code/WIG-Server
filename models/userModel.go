package models

import "github.com/google/uuid"

// User represents information about User profiles.
type User struct {

	// UserUID uniquely identifies the user
	UserUID uuid.UUID `json:"userUID" gorm:"primary_key;column:user_uid;type:varchar(191)"`
	
	// Username is the username associated with the users account
	Username string `json:"username" gorm:"column:username"`
	
	// Email is the email associated with the users account
	Email string `json:"email" gorm:"column:email"`
	
	// Salt is the salt associated with the users account for password security
	Salt string `json:"-" gorm:"column:salt"`
	
	// Hash is the hashed password associated with the user account
	Hash string `json:"-" gorm:"column:hash"`
	
	// EmailConfirm is a boolean expression representing if the users email has been confirmed
	EmailConfirm string `json:"emailConfirmed" gorm:"column:email_confirm;default:false"`
	
	// Tier represents the level of access the user may have
	Tier string `json:"tier" gorm:"column:tier"`
	
	// Token is the users login token
	Token string `json:"-" gorm:"column:token"`
}
