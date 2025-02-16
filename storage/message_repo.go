package storage

import (
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func SaveMessage(chatID uint, senderID uint, receiverID uint, content string) {
	db := config.GetDB()

	message := models.Message{
		ChatID:     chatID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	db.Create(&message)
}
func GetChatHistory(chatID uint) (*[]models.Message, error) {
	db := config.GetDB()

	messages := new([]models.Message)

	result := db.Table("messages").Where("chat_id = ?", chatID).Find(&messages)

	if err := result.Error; err != nil {
		return nil, err
	}

	return messages, nil

}
func MarkMessageAsRead(chatID uint) error{
	db := config.GetDB()

	result := db.Table("messages").Where("chat_id = ?", chatID).Set("seen", true)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}
