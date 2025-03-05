package ws

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
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
			delete(m.clients, client.user.ID)
			close(client.egress)
			m.mu.Unlock()
		}
	}
}

func ServeWS(manager *Manager, w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["userid"])
	if err != nil {
		return
	}

	user := storage.GetUserByID(uint(id))

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
	
	receiverID := event.Message.ReceiverID

	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()
	
	_, err := storage.GetChatByUserID(receiverID, event.Message.SenderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_, err = storage.CreatChat(event.Message.SenderID, receiverID)
	} else if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if ok {

		if err := storage.SaveMessage(event.Message, true); err != nil {
			return err
		}

		messages, err := storage.GetUnseenMessages(receiverID)
		if err != nil {
			receiver.egress <- event
		}else {
			for _, message := range *messages {
				evnt := Event{
					Message: &message,
				}
				receiver.egress <-evnt
			}
		}


	} else {
		if err := storage.SaveMessage(event.Message, false); err != nil {
			return err
		}
	}

	return nil
}
