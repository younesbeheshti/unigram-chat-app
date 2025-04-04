package ws

import (
	"time"

	"github.com/younesbeheshti/chatapp-backend/models"
)

// Time allowed to read the next pong message from the peer.
var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait*9) / 10
)


// Event is a message
type Event struct {
	Type string `json:"type"`
	*models.MessageRequest `json:"message"`
}

// Event types
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
