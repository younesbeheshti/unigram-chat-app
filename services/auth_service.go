package services

import (
	"fmt"

	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/storage"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username string, email string, password string) error {

	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	fmt.Println(username, email, password)

	user := new(models.User)
	user.Username = username
	user.Email = email
	user.PasswordHash = string(encpw)

	if err := storage.CreatUser(user); err != nil {
		return err
	}

	return nil
}

func LoginUser(email string, password string) error {

	user, err := storage.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if !user.ValidatePassword(password) {
		return fmt.Errorf("not registered")
	}

	// handle JWT token latter

	return nil
}