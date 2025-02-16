package config

import (
	"log"
	"os"

	"github.com/younesbeheshti/chatapp-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var db *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("POSTGRES_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}


func Init() {
	db = ConnectDB()

	db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})
}


func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("no db found")
	}
	return db
}
