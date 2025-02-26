package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/models"
)

type ClientList map[uint]*Client

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	user       *models.User
	egress     chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		user:       user,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.unregister <- c
		c.connection.Close()
	}()

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading messages: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			break
		}

		if err := c.manager.routeMessage(request, c); err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("error :", err)
				return
			}

		}
	}

}
