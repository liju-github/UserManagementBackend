package services

import (
	"errors"

	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/repository"
)

// AdminService struct that implements IAdminService
type AdminService struct {
	adminRepo *repository.AdminRepository // Admin repository interface
	userRepo  *repository.UserRepository
}

// NewAdminService constructs a new AdminService
func NewAdminService(adminRepo *repository.AdminRepository, userRepo *repository.UserRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo, // Passing the interface
		userRepo:  userRepo,
	}
}

// BlockUser implements IAdminService.
func (a *AdminService) BlockUser(userID string) error {
	
	return a.userRepo.BlockUser(userID)
}

// UnblockUser implements IAdminService.
func (a *AdminService) UnblockUser(userID string) error {
	return a.userRepo.UnblockUser(userID)
}

// DeleteUser implements IAdminService.
func (a *AdminService) DeleteUser(userID string) error {
	return a.userRepo.DeleteUser(userID)
}

// GetAllUsers implements IAdminService.
func (a *AdminService) GetAllUsers() ([]*models.User, error) {
	return a.userRepo.FindAllUsers()
}

// Login implements IAdminService.
func (a *AdminService) Login(email string, password string) (*models.Admin, error) {
	// Find the admin by email
	admin, err := a.adminRepo.FindAdminByEmail(email)
	if err != nil {
		return nil, err // Return the error if admin not found
	}

	// Check if the password matches (no hashing for simplicity)
	if admin.Password != password {
		return nil, errors.New("invalid password")
	}

	return admin, nil
}

// Logout implements IAdminService.
func (a *AdminService) Logout(adminID string) error {
	// Implement logout logic (e.g., invalidate session)
	panic("unimplemented")
}
