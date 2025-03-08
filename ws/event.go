package ws

import (
	"time"

	"github.com/younesbeheshti/chatapp-backend/models"
)


var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait*9) / 10
)

type Event struct {
	*models.MessageRequest `json:"message"`
}


type EventHandler func(event Event, c *Client) error


const (
	EventSendMessage = "send_message"
	EventNewMessage = "new_message"
)
