package storage

import (
	"log"
	"time"

	"github.com/younesbeheshti/chatapp-backend/config"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func SaveMessage(message *models.MessageRequest, seen bool) error {
	db := config.GetDB()

	msg := models.Message{
		ChatID:     message.ChatID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		Seen:       seen,
		CreatedAt:  time.Now(),
	}

	result := db.Create(&msg)
	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func GetChatHistory(chatID uint) ([]*models.MessageRequest, error) {
	db := config.GetDB()

	var messages []*models.Message

	result := db.Table("messages").Where("chat_id = ?", chatID).Find(&messages)

	if err := result.Error; err != nil {
		return nil, err
	}

	return messageModelToMessageReq(messages), nil

}
func MarkMessageAsRead(chatid uint) error {
	db := config.GetDB()

	result := db.Table("messages").Where("chat_id = ?", chatid).Update("seen", true)

	return result.Error
}

func GetUnseenMessages(receiverId uint) ([]*models.MessageRequest, error) {

	db := config.GetDB()

	var messages []*models.Message

	if err := db.Table("messages").Where("seen = ? and receiver_id = ?", false, receiverId).Find(&messages).Error; err != nil {
		return nil, err
	}

	if len(messages) != 0 {
		MarkMessageAsRead(messages[0].ReceiverID)
	}
	return messageModelToMessageReq(messages), nil
}

func messageModelToMessageReq(messages []*models.Message) []*models.MessageRequest {

	var msgs []*models.MessageRequest

	for _, message := range messages {

		username, err := GetUserNameByID(message.SenderID)
		if err != nil {
			log.Println("error", err)
			username = ""
		}

		msg := &models.MessageRequest{
			ChatID:     message.ChatID,
			SenderName: username,
			SenderID:   message.SenderID,
			ReceiverID: message.ReceiverID,
			Content:    message.Content,
		}

		msgs = append(msgs, msg)
	}

	return msgs
}
