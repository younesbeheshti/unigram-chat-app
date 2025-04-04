package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/younesbeheshti/chatapp-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var db *gorm.DB

// function to connect to the database
func ConnectDB() *gorm.DB {
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	
	dsn := os.Getenv("POSTGRES_URL")

	fmt.Println(dsn)
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
