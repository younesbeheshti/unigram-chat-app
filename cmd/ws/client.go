package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/cmd/models"
	"github.com/younesbeheshti/chatapp-backend/cmd/utils"
)

// ClientList is a map of clients
type ClientList map[uint]*Client

// PublicChannel is a map of clients
type PublicChannel map[*Client]bool

// Client is a middleman between the websocket connection and the hub
type Client struct {
	connection *websocket.Conn
	manager    *Manager
	user       *models.User
	egress     chan *utils.Event
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
	client := &Client{
		connection: conn,
		manager:    manager,
		user:       user,
		egress:     make(chan *utils.Event),
	}

	go client.ConsumePrivateMessages()

	return client
}

func (c *Client) ConsumePrivateMessages() {
	err := c.manager.rabbit.ConsumePrivateMessages(c.user.ID, func(event *utils.Event) {
		c.egress <- event
	})
	if err != nil {
		log.Println("failed to consume private msg", err)
	}
}

// readMessages reads messages from the websocket connection
func (c *Client) readMessages() {
	defer func() {
		c.manager.unregister <- c
		c.manager.pbleave <- c
		c.connection.Close()
	}()

	// Set read deadline
	if err := c.connection.SetReadDeadline(time.Now().Add(utils.PongWait)); err != nil {
		log.Println(err)
		return
	}
	// Set maximum message size
	c.connection.SetReadLimit(4096)
	// Set ping handler
	c.connection.SetPongHandler(c.pongHandler)

	for {

		// Read message
		_, payload, err := c.connection.ReadMessage()

		// Handle error
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading messages: %v", err)
			}
			break
		}

		// decrypt message
		//message, err := c.DecryptMessage(payload)
		//if err != nil {
		//	log.Fatal(err)
		//	return
		//}

		fmt.Println(string(payload))
		// Unmarshal message
		var request utils.Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			break
		}

		go printRequest(request)

		// Handle message to see if the user join or leave the channel
		if request.Type == utils.EventJoinChannel {
			c.manager.pbjoin <- c
			request.Type = utils.EventServerMessage
			request.Content = fmt.Sprintf("%v joined the chat room", request.SenderName)
			event := utils.Event{
				Type: utils.EventServerMessage,
				MessageRequest: &models.MessageRequest{
					SenderName: request.SenderName,
					Content:    fmt.Sprintf("hi %v, welcome to the chat room", request.SenderName),
				},
			}

			c.egress <- &event
		}

		// Handle message to see if the user join or leave the channel
		if request.Type == utils.EventLeaveChannel {
			request.Type = utils.EventServerMessage
			request.Content = fmt.Sprintf("%v left the chat room", request.SenderName)
			c.manager.pbleave <- c
		}

		// send to routemessage to send the message to the client
		if err := c.manager.routeMessage(&request, c); err != nil {
			log.Println(err)
		}
	}
}

// writeMessages writes messages to the websocket connection
func (c *Client) writeMessages() {
	defer func() {
		c.manager.unregister <- c
		c.manager.pbleave <- c

	}()

	// Set write deadline
	ticker := time.NewTicker(utils.PingInterval)

	for {

		// Write message
		select {
		case msg, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				return
			}

			// Marshal message
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("error marshalin", err)
				log.Println(err)
				continue
			}

			// encrypt message
			//message := c.EncryptMessage(data)

			go printRequest(*msg)

			// Write message
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("error :", err)
				return
			}

		// Write ping
		case <-ticker.C:
			log.Println("ping")

			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("write msg err:", err)
				return
			}
		}
	}

}

// pongHandler
func (c *Client) pongHandler(pongmsg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(utils.PongWait))
}

func printRequest(event utils.Event) {
	fmt.Printf("%v, length = %v: %v\n", event.Type, len(event.Content), event.Content)
}

func (c *Client) EncryptMessage(message []byte) []byte {
	return utils.Base64Encode(message)
}

func (c *Client) DecryptMessage(message []byte) ([]byte, error) {
	return utils.Base64Decode(message)
}
