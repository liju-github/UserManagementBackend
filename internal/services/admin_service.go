package services

import (
	"errors"

	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/repository"
)


type AdminService struct {
	adminRepo *repository.AdminRepository 
	userRepo  *repository.UserRepository
}


func NewAdminService(adminRepo *repository.AdminRepository, userRepo *repository.UserRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo, 
		userRepo:  userRepo,
	}
}


func (a *AdminService) BlockUser(userID string) error {

	return a.userRepo.BlockUser(userID)
}


func (a *AdminService) UnblockUser(userID string) error {
	return a.userRepo.UnblockUser(userID)
}


func (a *AdminService) DeleteUser(userID string) error {
	return a.userRepo.DeleteUser(userID)
}


func (a *AdminService) GetAllUsers() ([]*models.User, error) {
	return a.userRepo.FindAllUsers()
}


func (a *AdminService) Login(email string, password string) (*models.Admin, error) {
	
	admin, err := a.adminRepo.FindAdminByEmail(email)
	if err != nil {
		return nil, err 
	}

	
	if admin.Password != password {
		return nil, errors.New("invalid password")
	}

	return admin, nil
}


func (a *AdminService) Logout(adminID string) error {
	
	panic("unimplemented")
}
