package storage

import (
	"fmt"
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func CreateChat(user1ID uint, user2ID uint) (uint, error) {
	db := config.GetDB()

	_, err := GetChatByUserID(user1ID, user2ID)
	if err == nil {
		return 0, err
	}

	chat := models.Chat{
		User1ID: user1ID,
		User2ID: user2ID,
		CreatedAt: time.Now(), 
	}

	result := db.Create(&chat)
	if result.Error != nil {
		return 0, result.Error
	}

	return chat.ID, nil
}

func GetChatUsersByUserID(userID uint) ([]*models.User, error) {
	db := config.GetDB()

	var users []*models.User

	err := db.Raw(`
		SELECT DISTINCT u.* 
		FROM users u
		JOIN chats c ON (u.id = c.user1_id OR u.id = c.user2_id)
		WHERE (c.user1_id = ? OR c.user2_id = ?) AND u.id != ?
	`, userID, userID, userID).Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetChatsByUserID(userID uint) ([]*models.Chat, error) {
	db := config.GetDB()

	var chats []*models.Chat

	//TODO: if the user has any chat return the contact(*User model)

	result := db.Table("chats").Where("user1_id = ? or user2_id = ?", userID, userID).Find(&chats)

	if err := result.Error; err != nil {
		fmt.Println("error", err)
		return nil, err
	}

	return chats, nil
}

func GetChatByUserID(userID1, userID2 uint) (*models.Chat, error) {
	db := config.GetDB()

	var chat models.Chat

	result := db.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID1, userID2, userID2, userID1,
	).First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

	return &chat, nil
}

