package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/services"
	"github.com/younesbeheshti/chatapp-backend/storage"
	"github.com/younesbeheshti/chatapp-backend/utils"
)

var res models.Respnse
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.LoginRequst
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userid, err := services.LoginUser(req.Username, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		res.Message = "Unauthorized"
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Token, err = utils.GenerateJWT(userid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res.Message = "Unauthorized"
		json.NewEncoder(w).Encode(res)
		return
	}
	res.UserID = userid
	res.Message = "login successful"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "applcation/json")
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userid, err := services.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Token, err = utils.GenerateJWT(userid)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res.UserID = userid
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(mux.Vars(r)["userid"])
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := storage.GetUserByID(uint(id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetContactHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["userid"])
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp models.ContactsRespose
	resp.Contacts, err = storage.GetContact(uint(id))
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func GetChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := strconv.Atoi((mux.Vars(r)["chatid"]))
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

	var resp models.ChatResponse
	resp.Chats = chats

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
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

	// err = storage.SaveMessage(req.ChatID, req.SenderID, req.ReceiverID, req.Content)
	// if err != nil {
	// 	res.Message = err.Error()
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
}

func MarkMessagesReadHandler(w http.ResponseWriter, r *http.Request) {
}
