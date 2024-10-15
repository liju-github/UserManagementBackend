package repository

import (
	"errors"

	"github.com/liju-github/user-management/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	MySQLDatabase *gorm.DB
}


func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{MySQLDatabase: db}
}


func (repo *UserRepository) CreateUser(user *models.User) error {
	if err := repo.MySQLDatabase.Create(user).Error; err != nil {
		return errors.New("failed to create user: " + err.Error())
	}
	return nil
}


func (repo *UserRepository) FindUser(field string, value interface{}) (*models.User, error) {
	var user models.User
	if err := repo.MySQLDatabase.Where(field+" = ?", value).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to find user: " + err.Error())
	}
	return &user, nil
}


func (repo *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	return repo.FindUser("email", email)
}


func (repo *UserRepository) FindUserByID(userID string) (*models.User, error) {
	return repo.FindUser("id", userID)
}


func (repo *UserRepository) UpdateUser(user *models.User) error {
	if err := repo.MySQLDatabase.Save(user).Error; err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}


func (repo *UserRepository) DeleteUser(userID string) error {
	
	if err := repo.MySQLDatabase.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}
	return nil
}



func (repo *UserRepository) FindUserByVerificationToken(token string) (*models.User, error) {
	return repo.FindUser("verification_token", token)
}


func (repo *UserRepository) CreatePasswordReset(passwordReset *models.PasswordReset) error {
	if err := repo.MySQLDatabase.Create(passwordReset).Error; err != nil {
		return errors.New("failed to create password reset: " + err.Error())
	}
	return nil
}


func (repo *UserRepository) FindPasswordResetByToken(token string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	if err := repo.MySQLDatabase.Where("reset_token = ?", token).First(&passwordReset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("password reset not found")
		}
		return nil, errors.New("failed to find password reset by token: " + err.Error())
	}
	return &passwordReset, nil
}


func (repo *UserRepository) DeletePasswordReset(resetID string) error {
	if err := repo.MySQLDatabase.Delete(&models.PasswordReset{}, resetID).Error; err != nil {
		return errors.New("failed to delete password reset: " + err.Error())
	}
	return nil
}


func (repo *UserRepository) BlockUser(userID string) error {
	return repo.changeUserBlockStatus(userID, true)
}


func (repo *UserRepository) UnblockUser(userID string) error {
	return repo.changeUserBlockStatus(userID, false)
}


func (repo *UserRepository) FindAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := repo.MySQLDatabase.Find(&users).Error; err != nil {
		return nil, errors.New("failed to retrieve all users: " + err.Error())
	}
	return users, nil
}


func (repo *UserRepository) changeUserBlockStatus(userID string, isBlocked bool) error {
	result := repo.MySQLDatabase.Model(&models.User{}).Where("id = ?", userID).Update("is_blocked", isBlocked)
	if result.Error != nil {
		return errors.New("failed to change user block status: " + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found or already in the desired block state")
	}
	return nil
}
