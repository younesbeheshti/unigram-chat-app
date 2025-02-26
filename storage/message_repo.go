package storage

import (
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func SaveMessage(message *models.MessageRequest) error{
	db := config.GetDB()

	msg := models.Message{
		ChatID: message.ChatID,
		SenderID: message.SenderID,
		ReceiverID: message.ReceiverID,
		Content: message.Content,
		CreatedAt: time.Now(),
	}

	result := db.Create(&msg)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}
func GetChatHistory(chatID uint) (*[]models.Message, error) {
	db := config.GetDB()

	messages := new([]models.Message)

	result := db.Table("messages").Where("chat_id = ?", chatID).Order("created_at desc").Find(&messages)

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
