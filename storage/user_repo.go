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
func CreatUser(user *models.User) (uint, error) {
	db := config.GetDB()


	// TODO : returting the user id to 

	result := db.Create(&user)
	if err := result.Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}
func GetUserByID(userID uint) *models.User {
	db := config.GetDB()

	user := new(models.User)
	db.Table("users").Find(&user)
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
