package config

import (
	"os"
	"webapp/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.ProductCategory{})
	DB.AutoMigrate(&models.User{})
}
