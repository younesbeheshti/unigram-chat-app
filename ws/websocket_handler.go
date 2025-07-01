package ws

import (
	"encoding/json"
	"fmt"
	"github.com/younesbeheshti/chatapp-backend/rabbitmq"
	"github.com/younesbeheshti/chatapp-backend/utils"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/models"
	"github.com/younesbeheshti/chatapp-backend/storage"
)

// Manager maintains the set of active clients and broadcasts messages to the clients
type Manager struct {
	clients    ClientList
	pbChannel  PublicChannel
	pbjoin     chan *Client
	pbleave    chan *Client
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	rabbit     *rabbitmq.Service
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewManager creates a new Manager
func NewManager() *Manager {
	m := Manager{
		clients:    make(ClientList),
		pbChannel:  make(PublicChannel),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		pbjoin:     make(chan *Client),
		pbleave:    make(chan *Client),
		rabbit:     rabbitmq.NewRabbitService(),
	}

	// Start the manager
	go m.start()
	// Start the public channel
	go m.handlePublicChannel()
	return &m
}

// Start the manager
func (m *Manager) start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()

			m.clients[client.user.ID] = client

			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()

			client.connection.Close()
			delete(m.clients, client.user.ID)

			m.mu.Unlock()
		}
	}
}

// Start the public channel
func (m *Manager) handlePublicChannel() {
	for {
		select {
		case client := <-m.pbjoin:
			m.mu.Lock()
			m.pbChannel[client] = true
			m.mu.Unlock()

		case client := <-m.pbleave:
			m.mu.Lock()
			m.pbChannel[client] = false
			delete(m.pbChannel, client)
			m.mu.Unlock()
		}
	}
}

// ServeWS handles websocket requests from the peer
func ServeWS(manager *Manager, w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(uint)
	user := storage.GetUserByID(uint(id))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client
	client := NewClient(conn, manager, user)
	fmt.Println("creating client: ", user.ID)

	//kill the connection if the user does exist and make the new connection
	manager.CheckIfUserConnectedBefore(client)

	// Register client
	manager.register <- client

	// Start processing messages
	go client.readMessages()
	go client.writeMessages()
}

// routeMessage routes the message
func (m *Manager) routeMessage(event *utils.Event, sender *Client) error {

	if event.MessageRequest == nil {
		return fmt.Errorf("msg is nil")
	}

	// If receiverID is 0, then it's a channel message
	if event.MessageRequest.ReceiverID == 0 {
		return m.sendChannelMessage(event, sender)
	} else {
		return m.sendPrivateMessage(event)
	}

}

// sendChannelMessage
func (m *Manager) sendChannelMessage(event *utils.Event, sender *Client) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for client := range m.pbChannel {
		if client != sender {
			client.egress <- event
		}
	}

	// Save message
	return storage.SaveMessage(event.MessageRequest, true)
}

// sendPrivateMessage
func (m *Manager) sendPrivateMessage(event *utils.Event) error {
	receiverID := event.MessageRequest.ReceiverID

	err := m.rabbit.PublishPrivateMessages(event)
	if err != nil {
		return err
	}
	// Check if receiver is online
	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()

	// if receiver is online send message
	if ok {
		if err := storage.SaveMessage(event.MessageRequest, true); err != nil {
			return err
		}
		receiver.egress <- event

		// if receiver is offline save message
	} else {
		fmt.Println("saved message for offline user:", receiverID)
		if err := storage.SaveMessage(event.MessageRequest, false); err != nil {
			return err
		}
	}

	return nil
}

// GetActiveUsersHandler
func (m *Manager) GetActiveUsersHandler(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	activeUsers := []*models.User{}
	for client := range m.pbChannel {
		activeUsers = append(activeUsers, client.user)
		fmt.Println("active user: ", client.user.Username)
	}

	response := models.OnlineUsers{Users: activeUsers}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// function to check if the user was connected before or not
func (m *Manager) CheckIfUserConnectedBefore(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.clients {
		if c == client {
			m.unregister <- c
			return
		}
	}

	fmt.Println("no user found with client:", client.user.Username)
}
