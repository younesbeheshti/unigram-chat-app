package routes

import (
	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/handlers"
	"github.com/younesbeheshti/chatapp-backend/ws"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")

	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")
	router.HandleFunc("/chats/{userid}", handlers.GetChatHandler).Methods("GET")
	router.HandleFunc("/messages/{chatid}", handlers.GetMessagesHandler).Methods("GET")
	router.HandleFunc("/messages/read", handlers.MarkMessagesReadHandler).Methods("POST")
	router.HandleFunc("/ws", ws.HandleWebsocketConnection).Methods("GET")


	


	return router
}