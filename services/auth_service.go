package services

import "github.com/younesbeheshti/chatapp-backend/models"

func RegisterUser(username string, email string, password string) error {
	//hassedPassword := password
	user := models.User{
		Username: username,
		Email: email,
		PasswordHash: password,
	}
	return nil
}

func LoginUser(email string, password string) error {
	return nil
}