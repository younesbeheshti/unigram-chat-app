package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/models"
)

func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request) {}
func ReadMessages(conn *websocket.Conn, userID uint)  {
	
}
func BroadcastMessage(message *models.Message) {}
func DisconnectUser(userID uint) {}
