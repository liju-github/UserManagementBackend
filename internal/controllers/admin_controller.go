package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/services"
	"github.com/liju-github/user-management/internal/utils"
)


type AdminController struct {
	adminService *services.AdminService
}


func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

func (ac *AdminController) Login(c *fiber.Ctx) error {
	var admin models.AdminRequest
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	
	authAdmin, err := ac.adminService.Login(admin.Email, admin.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	
	token, err := utils.GenerateJWT(authAdmin.Email, authAdmin.ID, "admin",72)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"admin":   authAdmin,
		"token":   token,
	})
}



func (ac *AdminController) GetAllUsers(c *fiber.Ctx) error {
	
	users, err := ac.adminService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users: " + err.Error(),
		})
	}

	
	safeUsers := make([]models.UserProfileResponse, len(users))

	
	for i, user := range users {
		safeUsers[i] = models.UserProfileResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			Age:        user.Age,
			Gender:     user.Gender,
			Address:    user.Address,
			IsVerified: user.IsVerified,
			IsBlocked:  user.IsBlocked,
			CreatedAt: user.CreatedAt.String(),
		}
	}

	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users":safeUsers,
	})
}


func (ac *AdminController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Query("id") 
	if err := ac.adminService.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully!",
	})
}


func (ac *AdminController) BlockUser(c *fiber.Ctx) error {
	userID := c.Query("id") 
	if err := ac.adminService.BlockUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to block user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User blocked successfully!",
	})
}


func (ac *AdminController) UnblockUser(c *fiber.Ctx) error {
	userID := c.Query("id") 
	if err := ac.adminService.UnblockUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unblock user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User unblocked successfully!",
	})
}
