package main

import (
	"net/http"

	"github.com/younesbeheshti/chatapp-backend/handlers"
	"github.com/younesbeheshti/chatapp-backend/routes"
)



func main() {

	handlers.Init()


	router := routes.NewRouter()

	http.ListenAndServe(":8080", router)

}