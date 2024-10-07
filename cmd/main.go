package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/liju-github/user-management/internal/config"
	"github.com/liju-github/user-management/internal/controllers"
	"github.com/liju-github/user-management/internal/database"
	"github.com/liju-github/user-management/internal/repository"
	"github.com/liju-github/user-management/internal/services"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Use the logger middleware
	app.Use(logger.New(logger.Config{
		Output:     os.Stdout, // Log to terminal
		// Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Local",
	}))

	// Load environment configuration
	envConfig := config.EnvConfig()

	// Connect to the database
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
	adminService := services.NewAdminService(adminRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)

	fmt.Println(userController, adminController)

	// Auth group
	authGroup := app.Group("/api/auth")
	authGroup.Post("/signup", userController.Signup)
	authGroup.Post("/login", userController.Login)
	authGroup.Post("/logout", userController.Logout)
	authGroup.Post("/verify-email", userController.VerifyEmail)
	authGroup.Post("/resend-verification", userController.ResendVerification)
	authGroup.Post("/reset-password", userController.RequestPasswordReset)
	authGroup.Post("/confirm-reset-password", userController.ConfirmPasswordReset)

	// User group
	userGroup := app.Group("/api/user")
	userGroup.Get("/profile", userController.GetProfile)
	userGroup.Put("/update", userController.UpdateProfile)
	userGroup.Post("/upload-profile-picture", userController.UploadProfilePicture)

	// Admin group
	adminGroup := app.Group("/api/admin")
	adminGroup.Get("/users", adminController.GetAllUsers)
	adminGroup.Delete("/users/:id", adminController.DeleteUser)
	adminGroup.Post("/users/block/:id", adminController.BlockUser)
	adminGroup.Post("/users/unblock/:id", adminController.UnblockUser)

	// Start the Fiber server
	err := app.Listen(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
