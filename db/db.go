package db

import (
	"app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func SetupDB() {
	database, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Product{})
	// database.AutoMigrate(&model.Transaction{})

	db = database

}
