package storage

import (
	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)


func GetUserByUserName(username string) (*models.User, error) {
	db := config.GetDB()

	user := new(models.User)
	result := db.Table("users").Where("username = ?", username).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}
func GetUserByEmail(email string) (*models.User, error) {
	db := config.GetDB()

	user := new(models.User)

	result := db.Table("users").Where("email = ?", email).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}
func CreatUser(user *models.User) (uint, error) {
	db := config.GetDB()

	result := db.Create(&user)
	if err := result.Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}
func GetUserByID(userID uint) *models.User {
	db := config.GetDB()

	user := new(models.User)
	db.Table("users").Where("id = ?", userID).Find(&user)
	return user

}

func GetUsers() (*[]models.User, error) {
	db := config.GetDB()

	users := new([]models.User)
	result := db.Table("users").Find(&users)
	if err := result.Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetContact(userID uint) ([]*models.User, error) {
	db := config.GetDB()

	var contacts []*models.User
	result := db.Table("users").Where("id != ?", userID).Find(&contacts)
	if err := result.Error; err != nil {
		return nil, err
	}

	return contacts, nil
}
