package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/handlers"
	"github.com/younesbeheshti/chatapp-backend/ws"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	manager := ws.NewManager()

	router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")

	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")
	router.HandleFunc("/chats/{userid}", handlers.GetChatHandler).Methods("GET")
	router.HandleFunc("/messages/{chatid}", handlers.GetMessagesHandler).Methods("GET")
	router.HandleFunc("/messages/read", handlers.MarkMessagesReadHandler).Methods("POST")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(manager, w, r)
	})

	return router
}
