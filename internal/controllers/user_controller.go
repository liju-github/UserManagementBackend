package controllers

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/liju-github/user-management/internal/models"
	"github.com/liju-github/user-management/internal/services"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

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

	token, err := generateJWT(user.Email, user.ID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

func generateJWT(email string, userID string) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func JWTMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	tokenStr := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// Check if the token is valid and not expired
	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Check expiration
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["exp"] == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Compare the expiration time
	exp := claims["exp"].(float64) // JWT expiration is a float64
	if time.Now().Unix() > int64(exp) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
	}

	ctx.Locals("userID", claims["userID"])
	ctx.Locals("email", claims["email"])
	log.Println("Request from ", claims)

	return ctx.Next()
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

	// Create the profile response
	profileResponse := models.UserProfileResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		ImageURL:   user.ImageURL,
		IsVerified: user.IsVerified,
		IsBlocked:  user.IsBlocked,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"user": profileResponse})
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	userID, _ := ctx.Locals("userID").(string)
	email, ok := ctx.Locals("email").(string)
	fmt.Println(userID, email)
	if !ok || userID == "" || email == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid user ID"})
	}

	var updateReq models.UserUpdateRequest
	if err := ctx.BodyParser(&updateReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := c.userService.UpdateProfile(userID, email, &updateReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile updated successfully"})
}

func (c *UserController) UploadProfilePicture(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)

	cdnURL := ctx.Query("image_url")
	if cdnURL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Image URL not provided"})
	}

	_, err := url.ParseRequestURI(cdnURL)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL format"})
	}

	if err := c.userService.UploadProfilePicture(userID, cdnURL); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile picture uploaded successfully", "url": cdnURL})
}
