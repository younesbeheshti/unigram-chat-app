package storage

import (
	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	db := config.GetDB()

	user := new(models.User)

	result := db.Table("users").Where("email = ?", email).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}
func CreatUser(user *models.User) {
	db := config.GetDB()

	db.Create(&user)
}
func GetUserByID(userID uint) {}

