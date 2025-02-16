package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/services"
)


var res models.Respnse

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	email := params["email"]
	password := params["password"]
	err := services.LoginUser(email, password)
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


func RegisterUserHandler(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)
	username := params["username"]
	email := params["email"]
	password := params["password"]
	err := services.RegisterUser(username, email, password)
	if err != nil {
		res.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Message = "user successfully registered!"
	w.WriteHeader(http.StatusOK)	
	json.NewEncoder(w).Encode(res)
}


func GetChatHandler(w http.ResponseWriter, r *http.Request) {}
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {}
func MarkMessagesReadHandler(w http.ResponseWriter, r *http.Request) {}
