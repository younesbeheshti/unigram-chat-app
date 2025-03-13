package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/handlers"
	"github.com/younesbeheshti/chatapp-backend/middleware"
	"github.com/younesbeheshti/chatapp-backend/ws"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	manager := ws.NewManager()

	router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")

	router.HandleFunc("/auth/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/auth/login", handlers.LoginUserHandler).Methods("POST")

	// TODO : ?
	// router.Handle("/user/chats/{userid}", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetChatsHandler))).Methods("GET")
	router.Handle("/user/chat", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetChatsHandler))).Methods("GET")
	router.Handle("/user/chat/{user_id}", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetChatsHandler))).Methods("GET")
	router.Handle("/user/addchat",middleware.ValidateTokenHandler(http.HandlerFunc(handlers.AddChatHandler))).Methods("POST")
	router.Handle("/user/contacts", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetContactHandler))).Methods("GET")
	router.HandleFunc("/user/{userid}", handlers.GetUserHandler).Methods("GET")

	router.HandleFunc("/user/chat/messages/{chatid}", handlers.GetMessagesHandler).Methods("GET")

	router.HandleFunc("/messages/read", handlers.MarkMessagesReadHandler).Methods("POST")

	router.Handle("/ws", middleware.ValidateTokenHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(manager, w, r)
	})))

	return router
}
