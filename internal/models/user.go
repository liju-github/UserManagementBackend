// user.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 string          `gorm:"type:char(36);primaryKey;unique" json:"id"` // UUID stored as string in MySQL
	Name               string          `gorm:"type:varchar(100);not null" json:"name"`
	Email              string          `gorm:"type:varchar(255);unique;not null" json:"email"`
	ImageURL           string          `json:"imageurl"`
	Age                uint            `json:"age"`
	Gender             string          `json:"gender"`
	PasswordHash       string          `gorm:"type:varchar(255);not null" json:"password_hash"`
	IsVerified         bool            `gorm:"default:false" json:"is_verified"`
	IsBlocked          bool            `gorm:"default:false" json:"is_blocked"`
	VerificationToken  string          `gorm:"type:varchar(255)" json:"verification_token"`
	VerificationExpiry int64           `json:"verification_expiry"`                             
	PasswordResets     []PasswordReset `gorm:"foreignKey:UserID" json:"password_resets"` // One-to-many relationship with PasswordReset
}

// BeforeCreate hook to automatically set UUID and Unix timestamps before creating a user record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String() // Generate UUID for the user
	// u.VerificationExpiry = time.Now().Add(5 * time.Minute).Unix() // Set token expiry to 5 minutes from now
	return nil
}

type PasswordReset struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`    // UUID stored as string in MySQL
	UserID     string    `gorm:"type:char(36);not null" json:"user_id"` // Foreign key from User
	ResetToken string    `gorm:"type:varchar(255);unique;not null" json:"reset_token"`
	CreatedAt  time.Time `json:"created_at"` // Unix timestamp for record creation
	Expiry     int64     `json:"expiry"`     // Unix timestamp for token expiration
}

// BeforeCreate hook to automatically set UUID, Unix timestamps, and token expiration for PasswordReset
func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) (err error) {
	pr.ID = uuid.New().String() // Generate UUID for the password reset record
	// pr.Expiry = time.Now().Add(5 * time.Minute).Unix() // Set token expiry to 5 minutes from now
	return nil
}

type UserSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Age        string `json:"age"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	ImageURL   string `json:"image_url,omitempty"`
	IsVerified bool   `json:"is_verified"`
	IsBlocked  bool   `json:"is_blocked"`
	CreatedAt  string `json:"created_at"`
}

type UserUpdateRequest struct {
	Name   string `json:"name"`
	Age    uint   `json:"age"`
	Gender string `json:"gender"`
}
