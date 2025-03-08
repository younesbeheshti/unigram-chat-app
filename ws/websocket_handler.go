package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/storage"
	"gorm.io/gorm"
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
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewManager() *Manager {
	m := Manager{
		clients:    make(ClientList),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Event),
	}

	go m.start()
	return &m
}

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

func ServeWS(manager *Manager, w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user_id").(uint)

	user := storage.GetUserByID(uint(id))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, manager, user)
	fmt.Println("creating client: ", user.ID)

	manager.register <- client

	//Ask mehrshad ... 
	manager.SendUnseenMessages(user.ID)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) routeMessage(event *Event, sender *Client) error {

	if event.MessageRequest == nil {
		return fmt.Errorf("msg is nil")
	}

	receiverID := event.MessageRequest.ReceiverID

	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()

	_, err := storage.GetChatByUserID(receiverID, event.MessageRequest.SenderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error ->", err)
		_, err = storage.CreatChat(event.MessageRequest.SenderID, receiverID)
	} else if err != nil {
		fmt.Println("error ->", err)
		return err
	}
	if err != nil {
		return err
	}

	if ok {

		if err := storage.SaveMessage(event.MessageRequest, false); err != nil {
			return err
		}
		receiver.egress <- event

	} else {
		fmt.Println("saved message into db")

		if err := storage.SaveMessage(event.MessageRequest, false); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) SendUnseenMessages(receiverID uint) bool {

	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()

	if ok {

		messages, err := storage.GetUnseenMessages(receiverID)
		if err != nil {

			var evnt *Event
			for _, message := range messages {
				evnt = &Event{
					MessageRequest: message,
				}
				fmt.Println("evnt ->", evnt.MessageRequest.Content)
				receiver.egress <- evnt
			}
		}

		return true
	}

	return false
}
