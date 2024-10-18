package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/liju-github/user-management/internal/utils"
	"github.com/liju-github/user-management/internal/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) GetRefreshToken(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	ID := ctx.Locals("ID").(string)
	role := ctx.Locals("role").(string)

	accessToken, accessErr := utils.GenerateJWT(email, ID, role, 1)
	if accessErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": accessToken})
}
