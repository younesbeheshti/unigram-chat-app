package storage

import (
	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)



// GetUserNameByID from db
func GetUserNameByID(ID uint) (string, error) {
	db := config.GetDB()

	var user *models.User
	if err := db.Table("users").Where("id = ?", ID).First(&user).Error; err != nil {
		return "", err
	}

	return user.Username, nil
}

// GetUserByUserName from db
func GetUserByUserName(username string) (*models.User, error) {
	db := config.GetDB()

	user := new(models.User)
	result := db.Table("users").Where("username = ?", username).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail from db
func GetUserByEmail(email string) (*models.User, error) {
	db := config.GetDB()

	user := new(models.User)

	result := db.Table("users").Where("email = ?", email).Find(&user)
	if err := result.Error; err != nil {
		return nil, err
	}

	return user, nil
}

// CreatUser and save it to the db
func CreatUser(user *models.User) (uint, error) {
	db := config.GetDB()

	_, err := GetUserByUserName(user.Username)
	if err == nil {
		return 0, err
	}
	result := db.Create(&user)
	if err := result.Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// GetUserByID from db
func GetUserByID(userID uint) *models.User {
	db := config.GetDB()

	user := new(models.User)
	db.Table("users").Where("id = ?", userID).Find(&user)
	return user

}



// GetUsers from db
func GetUsers() (*[]models.User, error) {
	db := config.GetDB()

	users := new([]models.User)
	result := db.Table("users").Find(&users)
	if err := result.Error; err != nil {
		return nil, err
	}

	return users, nil
}


// GetContact from db
func GetContact(userID uint) ([]*models.User, error) {
	db := config.GetDB()

	var contacts []*models.User
	result := db.Table("users").Where("id != ?", userID).Find(&contacts)
	if err := result.Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

