package services

import (
	"fmt"

	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/storage"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username string, email string, password string) (uint, error) {

	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	user := new(models.User)
	user.Username = username
	user.Email = email
	user.PasswordHash = string(encpw)

	userID, err := storage.CreatUser(user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func LoginUser(username string, password string) (uint, error) {

	user, err := storage.GetUserByUserName(username)
	if err != nil {
		return 0, err
	}

	if !user.ValidatePassword(password) {
		return 0, fmt.Errorf("not registered")
	}

	return user.ID, nil
}
