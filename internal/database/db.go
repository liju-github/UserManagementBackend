package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/liju-github/user-management/internal/config"
	models "github.com/liju-github/user-management/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(env config.Env) *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v?parseTime=true", env.DBUSER, env.DBPASSWORD, env.DBNAME)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("database connection failed")
	}

	err = AutoMigrate(DB)
	if err != nil {
		log.Fatal(err.Error())
	}

	return DB

}

func AutoMigrate(DB *gorm.DB) error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.PasswordReset{},
	)

	if err != nil {
		return errors.New("failed to automigrate models" )
	}

	return nil
}
