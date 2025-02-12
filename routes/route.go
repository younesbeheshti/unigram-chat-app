package routes

import (
	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/handlers"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/login", handlers.HandlerLogin).Methods("POST", "OPTIONS")

	return router
}