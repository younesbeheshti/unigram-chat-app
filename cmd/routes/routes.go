package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/cmd/handlers"
	"github.com/younesbeheshti/chatapp-backend/cmd/middleware"
	"github.com/younesbeheshti/chatapp-backend/cmd/ws"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	manager := ws.NewManager()

	// get users
	router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")

	// auth
	router.HandleFunc("/auth/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/auth/login", handlers.LoginUserHandler).Methods("POST")

	// user routes with middleware validation token
	router.Handle("/active-users", middleware.ValidateTokenHandler(http.HandlerFunc(manager.GetActiveUsersHandler))).Methods("GET")
	router.Handle("/user/chat", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetChatsHandler))).Methods("GET")
	router.Handle("/user/chat/{user_id}", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetChatsHandler))).Methods("GET")
	router.Handle("/user/addchat", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.AddChatHandler))).Methods("POST")
	router.Handle("/user/contacts", middleware.ValidateTokenHandler(http.HandlerFunc(handlers.GetContactHandler))).Methods("GET")
	router.HandleFunc("/user/{userid}", handlers.GetUserHandler).Methods("GET")

	// chat routes
	router.HandleFunc("/user/chat/messages/{chatid}", handlers.GetMessagesHandler).Methods("GET")

	// message
	router.HandleFunc("/messages/read", handlers.MarkMessagesReadHandler).Methods("POST")

	// ws
	router.Handle("/ws", middleware.ValidateTokenHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(manager, w, r)
	})))

	return router
}
