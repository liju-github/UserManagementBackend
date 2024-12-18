package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 string          `gorm:"type:char(36);primaryKey;unique" json:"id"`
	Name               string          `gorm:"type:varchar(100);not null" json:"name"`
	Email              string          `gorm:"type:varchar(255);unique;not null" json:"email"`
	Address            string          `json:"address"`
	ImageURL           string          `json:"imageurl"`
	Age                uint            `json:"age"`
	Gender             string          `json:"gender"`
	PhoneNumber        uint            `json:"phonenumber"`
	PasswordHash       string          `gorm:"type:varchar(255);not null" json:"password_hash"`
	IsVerified         bool            `gorm:"default:false" json:"is_verified"`
	IsBlocked          bool            `gorm:"default:false" json:"is_blocked"`
	VerificationToken  string          `gorm:"type:varchar(255)" json:"verification_token"`
	VerificationExpiry int64           `json:"verification_expiry"`
	PasswordResets     []PasswordReset `gorm:"foreignKey:UserID" json:"password_resets"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return nil
}

type PasswordReset struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID     string    `gorm:"type:char(36);not null" json:"user_id"`
	ResetToken string    `gorm:"type:varchar(255);unique;not null" json:"reset_token"`
	CreatedAt  time.Time `json:"created_at"`
	Expiry     int64     `json:"expiry"`
}

func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) (err error) {
	pr.ID = uuid.New().String()

	return nil
}

type UserSignupRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Age         uint   `json:"age" validate:"required,gte=18,lte=120"` // Validating age between 18 and 120
	Gender      string `json:"gender" validate:"required,oneof=Male Female Other"`
	Address     string `json:"address" validate:"required,min=5,max=100"`
	PhoneNumber uint   `json:"phonenumber" validate:"required,numeric,min=1000000000,max=9999999999"` // Validating as a number
	Password    string `json:"password" validate:"required,min=8"`
	ImageURL    string `json:"image_url" validate:"required,url"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Age         uint   `json:"age"`
	Gender      string `json:"gender"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber uint   `json:"phonenumber"`
	ImageURL    string `json:"image_url,omitempty"`
	IsVerified  bool   `json:"is_verified"`
	IsBlocked   bool   `json:"is_blocked"`
	CreatedAt   string `json:"created_at"`
}

type UserUpdateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      uint   `json:"age"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
	ImageURL string `json:"image_url,omitempty"`
}
