package repository

import (
	"github.com/liju-github/user-management/internal/models"
	"gorm.io/gorm"
)



// AdminRepository concrete implementation of IAdminRepository
type AdminRepository struct {
	MySQLDatabase *gorm.DB // Holds the Gorm DB instance
}

// Constructor for AdminRepository
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		MySQLDatabase: db,
	}
}

// Implement the FindAdminByEmail method from the interface
func (repo *AdminRepository) FindAdminByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := repo.MySQLDatabase.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err // Return error if admin not found
	}
	return &admin, nil // Return the found admin
}

