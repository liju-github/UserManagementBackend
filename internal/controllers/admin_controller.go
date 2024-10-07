package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liju-github/user-management/internal/services"
)

// AdminController struct with adminService interface
type AdminController struct {
	adminService *services.AdminService // No need to pass a pointer to an interface
}

// Constructor for AdminController
func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// Get All Users
func (ac *AdminController) GetAllUsers(c *fiber.Ctx) error {
	// Logic to get all users
	users := []string{"User1", "User2", "User3"} 
	return c.Status(fiber.StatusOK).JSON(users)
}

// Delete User
func (ac *AdminController) DeleteUser(c *fiber.Ctx) error {
	// Logic to delete user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully!",
	})
}

// Block User
func (ac *AdminController) BlockUser(c *fiber.Ctx) error {
	// Logic to block a user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User blocked successfully!",
	})
}

// Unblock User
func (ac *AdminController) UnblockUser(c *fiber.Ctx) error {
	// Logic to unblock a user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User unblocked successfully!",
	})
}
