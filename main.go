package main

import (
	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/routes"
)

func main() {

	router := routes.SetupRoutes()

	server := config.NewServer("127.0.0.1", ":15000", "TCP")
	server.InitServer(router)

}
