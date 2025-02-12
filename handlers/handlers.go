package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/models"
)


var users []models.User

func Init()  {
	users = append(users, models.User{
		ID: 1,
		Name: "yones",
		Password: "1234",
	})

	users = append(users, models.User{
		ID: 2,
		Name: "yalda",
		Password: "1234",
	})
}

func GetUsers() []models.User {
	return users
}


func HandlerLogin(w http.ResponseWriter, r *http.Request) {

	user := mux.Vars(r)

	for _, u := range users {
		if u.Name == user["name"] {
			if u.Password == user["password"] {
				res := models.Respnse{
					Message: fmt.Sprintf("user with id: %v successfuly loged in", u.ID),
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(res)
				return
			}
		}
	}

	res := models.Respnse{
		Message: "user not found!",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(res)

}