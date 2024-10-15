package repository

import (
	"github.com/liju-github/user-management/internal/models"
	"gorm.io/gorm"
)




type AdminRepository struct {
	MySQLDatabase *gorm.DB 
}


func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		MySQLDatabase: db,
	}
}


func (repo *AdminRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := repo.MySQLDatabase.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err 
	}
	return &admin, nil 
}

