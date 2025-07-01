package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/younesbeheshti/chatapp-backend/cmd/config"
	"github.com/younesbeheshti/chatapp-backend/cmd/routes"
	"log"
	"os"
)

func main() {

	router := routes.SetupRoutes()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	port := fmt.Sprintf(":%s",os.Getenv("PORT"))

	if port == "" {
		port = ":8080"
	}
	fmt.Println(port)

	server := config.NewServer("127.0.0.1", port, "TCP")
	server.InitServer(router)

}
