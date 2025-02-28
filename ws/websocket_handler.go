package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/storage"
	"github.com/younesbeheshti/chatapp-backend/utils"
)

type Manager struct {
	clients    ClientList
	register   chan *Client
	unregister chan *Client
	broadcast  chan Event
	mu         sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


func NewManager() *Manager {
	m := Manager{
		clients: make(ClientList),
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan Event),
	}

	go m.start()
	return &m
}

func (m *Manager) start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.user.ID]= client
			m.mu.Unlock()

		case client := <-m.unregister:
			m.mu.Lock()
			delete(m.clients, client.user.ID)
			close(client.egress)
			m.mu.Unlock()
		}
	}
}


func ServeWS(manager *Manager, w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := utils.ValidateToket(tokenString)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, manager, user)

	manager.register <- client

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) routeMessage(event Event, sender *Client) error {
	if event.Message != nil {
		return fmt.Errorf("msg is nil")
	}
	if err := storage.SaveMessage(event.Message); err != nil {
		return err
	}
	receiverID := event.Message.ReceiverID

	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()

	if ok {
		receiver.egress <- event
	} else {
		return fmt.Errorf("user not found %v", receiverID)
	}

	return nil
}
