package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liju-github/user-management/internal/config"
	"github.com/liju-github/user-management/internal/database"
)

func main() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	fmt.Println("server running successfully", database.ConnectDatabase(config.EnvConfig()))
	router.Run(":8080")

}


//  // Auth group
//  authGroup := app.Group("/api/auth")
//  authGroup.Post("/signup", Signup)
//  authGroup.Post("/login", Login)
//  authGroup.Post("/logout", Logout)
//  authGroup.Post("/verify-email", VerifyEmail)
//  authGroup.Post("/resend-verification", ResendVerification)
//  authGroup.Post("/reset-password", RequestPasswordReset)
//  authGroup.Post("/confirm-reset-password", ConfirmPasswordReset)

//  // User group
//  userGroup := app.Group("/api/user")
//  userGroup.Get("/profile", GetProfile)
//  userGroup.Put("/update", UpdateProfile)
//  userGroup.Post("/upload-profile-picture", UploadProfilePicture)

//  // Admin group
//  adminGroup := app.Group("/api/admin")
//  adminGroup.Get("/users", GetAllUsers)
//  adminGroup.Delete("/users/:id", DeleteUser)

//  // Start server
//  app.Listen(":3000")