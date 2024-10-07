package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"` // UUID stored as string in MySQL
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Email     string `gorm:"type:varchar(255);unique;not null" json:"email"`
	CreatedAt int64  `json:"created_at"` // Unix timestamp for record creation
}

// BeforeCreate hook to automatically set UUID and Unix timestamp before creating an admin record
func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()   // Generate UUID for the admin
	a.CreatedAt = time.Now().Unix() // Set the created_at field to the current Unix timestamp
	return nil
}
