package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/liju-github/user-management/internal/config"
	"github.com/liju-github/user-management/internal/controllers"
	"github.com/liju-github/user-management/internal/database"
	"github.com/liju-github/user-management/internal/repository"
	"github.com/liju-github/user-management/internal/services"
	"github.com/liju-github/user-management/internal/utils"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "server is running"})
	})

	// Use the logger middleware
	app.Use(logger.New(logger.Config{
		Output: os.Stdout,
		// Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}))

	envConfig := config.EnvConfig()

	db := database.ConnectDatabase(envConfig)
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	fmt.Println("Server running successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	adminService := services.NewAdminService(adminRepo, userRepo)
	authService := services.NewAuthService(adminRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)
	authController := controllers.NewAuthController(authService)

	fmt.Println(userController, adminController, authController)

	// Auth group
	authGroup := app.Group("/api/auth")
	authGroup.Post("/signup", userController.Signup)
	authGroup.Post("/login", userController.Login)
	authGroup.Post("/admin/login", adminController.Login)
	authGroup.Get("/verify-email/:token", userController.VerifyEmail)
	authGroup.Post("/resend-verification", userController.ResendVerification)
	authGroup.Post("/reset-password", userController.RequestPasswordReset)
	authGroup.Post("/confirm-reset-password", userController.ConfirmPasswordReset)
	authGroup.Get("/refresh", utils.JWTMiddleware("", userRepo), authController.GetRefreshToken)

	// User group
	userGroup := app.Group("/api/user")
	userGroup.Use(utils.JWTMiddleware("user", userRepo))
	userGroup.Get("/profile", userController.GetProfile)
	userGroup.Put("/update", userController.UpdateProfile)
	userGroup.Post("/upload-profile-picture", userController.UploadProfilePicture)

	// Admin group
	adminGroup := app.Group("/api/admin")
	adminGroup.Use(utils.JWTMiddleware("admin", userRepo))
	adminGroup.Get("/users", adminController.GetAllUsers)
	adminGroup.Delete("/users/", adminController.DeleteUser)
	adminGroup.Put("/users/block/", adminController.BlockUser)
	adminGroup.Put("/users/unblock/", adminController.UnblockUser)

	// Start the Fiber server
	err := app.Listen(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
