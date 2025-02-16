package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/services"
	"github.com/younesbeheshti/chatapp-backend/storage"
)

var res models.Respnse

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var req models.LoginRequst
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res.Message = "user not found"
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	res.Message = "user found!"
	json.NewEncoder(w).Encode(res)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := services.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Message = "user successfully registered!"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := storage.GetUsers()
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
func GetChatHandler(w http.ResponseWriter, r *http.Request) {	
	w.Header().Set("Content-Type", "application/json")
	userID,err := strconv.Atoi((mux.Vars(r)["chatid"]))
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	chats, err := storage.GetChatsByUserID(uint(userID))
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chats)

}
func GetMessagesHandler(w http.ResponseWriter, r *http.Request)      {
	w.Header().Set("Content-Type", "application/json")
	chatid, err := strconv.Atoi(mux.Vars(r)["chatid"])
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := storage.GetChatHistory(uint(chatid))
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req models.MessageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.SaveMessage(req.ChatID, req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func MarkMessagesReadHandler(w http.ResponseWriter, r *http.Request) {
}


