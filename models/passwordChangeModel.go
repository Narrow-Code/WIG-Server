// models defines the data models used in the WIG-Server application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Borrower represents information about a borrower.
type PasswordChange struct {
    // Token is the unique token generated for email verification and serves as the primary key
    PasswordChangeToken string `json:"password_change_token" gorm:"primary_key;column:password_change_token;type:varchar(255);not null"`

    // UserID is the unique identifier of the user associated with the verification
    PasswordUserID uuid.UUID `json:"password_userID" gorm:"column:password_user_id;type:varchar(191);not null"`

    // ExpiresAt is the time when the token will expire
    PasswordExpiresAt time.Time `json:"password_expires_at" gorm:"column:password_expires_at;type:timestamp;not null"`

    // CreatedAt and UpdatedAt are automatically managed by GORM
    CreatedAt time.Time
    UpdatedAt time.Time

    // User is the user account associated with the ownership
    PasswordUser User `json:"password_user_id" gorm:"foreignkey:password_user_id"`	
}
