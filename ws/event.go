package ws

import "github.com/younesbeheshti/chatapp-backend/models"

type Event struct {
	Message *models.MessageRequest `json:"message"`
}


type EventHandler func(event Event, c *Client) error


const (
	EventSendMessage = "send_message"
	EventNewMessage = "new_message"
)
