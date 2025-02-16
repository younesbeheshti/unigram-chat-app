package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/routes"
)



func main() {


	router := routes.SetupRoutes()
	config.Init()

	fmt.Println("server is up on port: 8080")

	log.Fatal(http.ListenAndServe(":8080", router))

}