package controllers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/services"
	"github.com/liju-github/user-management/internal/utils"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) Signup(ctx *fiber.Ctx) error {
	var userReq models.UserSignupRequest
	if err := ctx.BodyParser(&userReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input" + err.Error()})
	}

	fmt.Println(userReq)

	if err := c.userService.Signup(&userReq); err != nil {
		if err.Error() == models.UserAlreadyExists {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User signed up successfully!"})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var loginReq models.UserLoginRequest
	if err := ctx.BodyParser(&loginReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := c.userService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	accessToken, accessErr := utils.GenerateJWT(user.Email, user.ID, "user", 1)
	refreshToken, refreshErr := utils.GenerateJWT(user.Email, user.ID, "user", 72)
	if accessErr != nil || refreshErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Login successful",
		"user":          user,
		"token":         accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *UserController) GetRefreshToken(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	userEmail := ctx.Locals("email").(string)
	
	accessToken, accessErr := utils.GenerateJWT(userEmail, userID, "user", 1)
	if accessErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": accessToken})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	if err := c.userService.Logout(userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}

func (c *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is required"})
	}

	if err := c.userService.VerifyEmail(token); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email verified successfully"})
}

func (c *UserController) ResendVerification(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := c.userService.ResendVerification(req.Email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Verification email resent"})
}

func (c *UserController) RequestPasswordReset(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := c.userService.RequestPasswordReset(req.Email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset email sent"})
}

func (c *UserController) ConfirmPasswordReset(ctx *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := c.userService.ConfirmPasswordReset(req.Token, req.NewPassword); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successfully"})
}

func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid user ID"})
	}

	user, err := c.userService.GetProfile(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	profileResponse := models.UserProfileResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Age:        user.Age,
		Gender:     user.Gender,
		Address:    user.Address,
		ImageURL:   user.ImageURL,
		IsVerified: user.IsVerified,
		IsBlocked:  user.IsBlocked,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}

	fmt.Println(profileResponse)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"user": profileResponse})
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {

	email, _ := ctx.Locals("email").(string)
	fmt.Println(email)

	var updateReq models.UserUpdateRequest
	if err := ctx.BodyParser(&updateReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input" + err.Error()})
	}
	userID := ""

	if updateReq.ID != "" {
		userID = updateReq.ID
	} else {
		userID, _ = ctx.Locals("userID").(string)
	}

	if err := c.userService.UpdateProfile(userID, email, &updateReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile updated successfully"})
}

func (c *UserController) UploadProfilePicture(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)

	type ImageUploadRequest struct {
		ImageURL string `json:"image_url"`
	}

	var req ImageUploadRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := url.ParseRequestURI(req.ImageURL)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL format"})
	}

	if err := c.userService.UploadProfilePicture(userID, req.ImageURL); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile picture uploaded successfully", "url": req.ImageURL})
}
