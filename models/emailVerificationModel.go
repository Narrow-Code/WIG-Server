// models defines the data models used in the WIG-Server application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Borrower represents information about a borrower.
type EmailVerification struct {
    // Token is the unique token generated for email verification and serves as the primary key
    VerificationToken string `json:"verification_token" gorm:"primary_key;column:verification_token;type:varchar(255);not null"`

    // UserID is the unique identifier of the user associated with the verification
    UserID uuid.UUID `json:"userID" gorm:"column:user_id;type:varchar(191);not null"`

    // ExpiresAt is the time when the token will expire
    ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;type:timestamp;not null"`

    // CreatedAt and UpdatedAt are automatically managed by GORM
    CreatedAt time.Time
    UpdatedAt time.Time

    // User is the user account associated with the ownership
    User User `json:"user" gorm:"foreignkey:user_id"`
	
}
