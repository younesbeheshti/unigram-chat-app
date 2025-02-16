package storage

import (
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func CreatChat(user1ID uint, user2ID uint) error{
	db := config.GetDB()

	chat := new(models.Chat)
	chat.User1ID = user1ID
	chat.User2ID = user2ID
	chat.CreatedAt = time.Now()

	result := db.Create(&chat)

	if err := result.Error; err != nil {
		return err
	}

	return nil

}

func GetChatsByUserID(userID uint) (*[]models.Chat, error){
	db := config.GetDB()

	chats := new([]models.Chat)

	result := db.Table("chats").Where("user1_id = ? or user2_id = ?", userID, userID).Find(&chats)

	if err := result.Error; err != nil {
		return nil, err
	}

	return chats, nil
}