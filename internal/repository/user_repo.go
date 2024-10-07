package repository

import (
	"github.com/liju-github/user-management/internal/models"
	"gorm.io/gorm"
)


// UserRepository struct that implements IUserRepository
type UserRepository struct {
	MySQLDatabase *gorm.DB // Holds the Gorm DB instance
}

// Constructor for UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		MySQLDatabase: db,
	}
}

// CreateUser implements IUserRepository.
func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.MySQLDatabase.Create(user).Error
}

// FindUserByEmail implements IUserRepository.
func (repo *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.MySQLDatabase.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByID implements IUserRepository.
func (repo *UserRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := repo.MySQLDatabase.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements IUserRepository.
func (repo *UserRepository) UpdateUser(user *models.User) error {
	return repo.MySQLDatabase.Save(user).Error
}

// DeleteUser implements IUserRepository.
func (repo *UserRepository) DeleteUser(userID string) error {
	return repo.MySQLDatabase.Delete(&models.User{}, userID).Error
}

// UpdateUserBlockStatus implements IUserRepository.
func (repo *UserRepository) UpdateUserBlockStatus(userID string, isBlocked bool) error {
	return repo.MySQLDatabase.Model(&models.User{}).Where("id = ?", userID).Update("is_verified", isBlocked).Error
}

// FindAllUsers implements IUserRepository.
func (repo *UserRepository) FindAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := repo.MySQLDatabase.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
