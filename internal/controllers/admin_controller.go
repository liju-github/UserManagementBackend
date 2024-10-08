package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/services"
)

// AdminController struct with adminService interface
type AdminController struct {
	adminService *services.AdminService
}

// NewAdminController constructs a new AdminController
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

	// Call the AdminService to perform login
	authAdmin, err := ac.adminService.Login(admin.Email, admin.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate JWT token
	token, err := generateJWT(authAdmin.Email,authAdmin.ID)
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

// GetAllUsers retrieves all users and responds with JSON
func (ac *AdminController) GetAllUsers(c *fiber.Ctx) error {
	users, err := ac.adminService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// DeleteUser deletes a user by userID and responds with a message
func (ac *AdminController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Query("id") // Extract userID from URL parameters
	if err := ac.adminService.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully!",
	})
}

// BlockUser blocks a user by userID and responds with a message
func (ac *AdminController) BlockUser(c *fiber.Ctx) error {
	userID := c.Query("id") // Extract userID from URL parameters
	if err := ac.adminService.BlockUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to block user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User blocked successfully!",
	})
}

// UnblockUser unblocks a user by userID and responds with a message
func (ac *AdminController) UnblockUser(c *fiber.Ctx) error {
	userID := c.Query("id") // Extract userID from URL parameters
	if err := ac.adminService.UnblockUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unblock user: " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User unblocked successfully!",
	})
}
