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
	Type string `json:"type"`
	*models.MessageRequest `json:"message"`
}


type EventHandler func(event Event, c *Client) error


const (
	EventSendMessage = "send_message"
	EventNewMessage = "new_message"
	EventJoinChannel = "join_channel"
	EventLeaveChannel = "leave_channel"
	EventServerMessage = "server_message"
	EventFileMessage = "file_message"
	EventPublicMessage = "public_message"
	EventPrivateMessage = "private_message"
)
