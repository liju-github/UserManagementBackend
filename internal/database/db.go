package database

import (
	"fmt"
	"log"

	"github.com/liju-github/user-management/internal/config"
	models "github.com/liju-github/user-management/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


// ConnectDatabase initializes the connection to the database
func ConnectDatabase(env config.Env) *gorm.DB {
	// Database connection string (DSN)
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?parseTime=true", env.DBUSER, env.DBPASSWORD, env.DBNAME)

	// Open MySQL database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
		return nil
	}

	// Automigrate the models (creates tables if not exists)
	err = AutoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to automigrate models: %v", err)
		return nil
	}

	return db
}

// AutoMigrate handles the automatic migration of models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.PasswordReset{},
	)
}
