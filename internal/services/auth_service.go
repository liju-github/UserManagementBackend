package services

import "github.com/liju-github/user-management/internal/repository"

type AuthService struct {
	adminRepo *repository.AdminRepository
	userRepo  *repository.UserRepository
}

func NewAuthService(adminRepo *repository.AdminRepository, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}
