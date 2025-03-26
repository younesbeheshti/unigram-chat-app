package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/models"
)

type ClientList map[uint]*Client
type PublicChannel map[*Client]bool

type Client struct {
    connection *websocket.Conn
    manager    *Manager
    user       *models.User
    egress     chan *Event
}

func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
    return &Client{
        connection: conn,
        manager:    manager,
        user:       user,
        egress:     make(chan *Event),
    }
}


func (c *Client) readMessages() {
	defer func() {
		c.manager.unregister <- c
		c.connection.Close()
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.connection.SetReadLimit(4096)
	c.connection.SetPongHandler(c.pongHandler)

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

		if request.Type == EventJoinChannel {
			c.manager.pbjoin <- c
			request.Type = EventServerMessage
			request.Content = fmt.Sprintf("%v joined the chat room", request.SenderName)
			event := Event{
				Type: EventServerMessage,
				MessageRequest: &models.MessageRequest{
					SenderName: request.SenderName,
					Content:    fmt.Sprintf("hi %v, welcome to the chat room", request.SenderName),
				},
			}

			c.egress <- &event
		}

		if request.Type == EventLeaveChannel {
			request.Type = EventServerMessage
			request.Content = fmt.Sprintf("%v left the chat room", request.SenderName)
			c.manager.pbleave <- c
		}

		if err := c.manager.routeMessage(&request, c); err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.unregister <- c
	}()

	ticker := time.NewTicker(pingInterval)

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
				fmt.Println("error marshalin", err)
				log.Println(err)
				continue
			}

			fmt.Println(msg)

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("error :", err)
				return
			}
		case <-ticker.C:
			log.Println("ping")
			
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("write msg err:", err)
				return 
			}
		}
	}

}

func (c *Client) pongHandler(pongmsg  string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
