package storage

import (
	"database/sql"
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func CreatChat(user1ID uint, user2ID uint) (*models.Chat, error) {
	db := config.GetDB()

	chat := new(models.Chat)
	chat.User1ID = user1ID
	chat.User2ID = user2ID
	chat.CreatedAt = time.Now()

	result := db.Create(&chat)

	if err := result.Error; err != nil {
		return nil, err
	}

	return chat, nil

}

func GetChatUsersByUserID(userID uint) (*[]models.User, error) {
	db := config.GetDB()

	var users []models.User

	err := db.Raw(`
		SELECT DISTINCT u.* 
		FROM users u
		JOIN chats c ON (u.id = c.user1_id OR u.id = c.user2_id)
		WHERE (c.user1_id = ? OR c.user2_id = ?) AND u.id != ?
	`, userID, userID, userID).Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return &users, nil
}

func GetChatsByUserID(userID uint) (*[]models.Chat, error) {
	db := config.GetDB()

	chats := new([]models.Chat)

	//TODO: if the user has any chat return the contact(*User model)

	result := db.Table("chats").Where("user1_id = ? or user2_id = ?", userID, userID).Find(&chats)

	if err := result.Error; err != nil {
		return nil, err
	}

	return chats, nil
}

func GetChatByUserID(userID2 uint, userID1 uint) (*models.Chat, error) {
	var chat *models.Chat
	db := config.GetDB()
	result := db.Table("chats").Where("(user1_id = @user1 and user2_id = @user2) or (user1_id = @user2 and user2_id = @user1)",
		sql.Named("user1", userID1), sql.Named("user2", userID2)).Find(&chat)
	if err := result.Error; err != nil {
		return nil, err
	}
	return chat, nil

}
