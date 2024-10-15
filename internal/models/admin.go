package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"` 
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Email     string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string `gorm:"type:varchar(100);not null" json:"password"`
	CreatedAt int64  `json:"created_at"` 
}

type AdminRequest struct {
	Email    string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password string `json:"password"`
}


func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()      
	a.CreatedAt = time.Now().Unix() 
	return nil
}
