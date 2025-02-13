package routes

import (
	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/handlers"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")
	router.HandleFunc("/chats", handlers.GetChatHandler).Methods("GET")
	router.HandleFunc("/messages/{chatid}", handlers.GetChatHandler).Methods("GET")
	router.HandleFunc("/messages/read", handlers.MarkMessagesReadHandler).Methods("POST")
	router.HandleFunc("/ws", handlers.HandleWebsocketConnection).Methods("GET")


	


	return router
}