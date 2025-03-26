package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/younesbeheshti/chatapp-backend/storage"
)

type Manager struct {
	clients    ClientList
	pbChannel  PublicChannel
	pbjoin chan *Client
	pbleave chan *Client
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewManager() *Manager {
	m := Manager{
		clients:    make(ClientList),
		pbChannel:  make(PublicChannel),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go m.start()
	go m.handlePublicChannel()
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

func (m *Manager) handlePublicChannel() {
	for {
		select {
			case client :=<- m.pbjoin:
				m.mu.Lock()
				m.pbChannel[client] = true
				m.mu.Unlock()

			case client :=<- m.pbleave:
				m.mu.Lock()
				m.pbChannel[client] = false
				delete(m.pbChannel, client)
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

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) routeMessage(event *Event, sender *Client) error {

	if event.MessageRequest == nil {
		return fmt.Errorf("msg is nil")
	}

	if event.MessageRequest.ReceiverID == 0 {
		return m.sendChannelMessage(event, sender)
	} else {
		return m.sendPrivateMessage(event)
	}

}

func (m *Manager) sendChannelMessage(event *Event, sender *Client) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for client := range m.pbChannel {
		if client != sender {
			client.egress <- event
		}
	}

	return storage.SaveMessage(event.MessageRequest, true)
}

func (m *Manager) sendPrivateMessage(event *Event) error {
	receiverID := event.MessageRequest.ReceiverID

	m.mu.Lock()
	receiver, ok := m.clients[receiverID]
	m.mu.Unlock()

	if ok {
		if err := storage.SaveMessage(event.MessageRequest, true); err != nil {
			return err
		}
		receiver.egress <- event
	} else {
		fmt.Println("saved message for offline user:", receiverID)
		if err := storage.SaveMessage(event.MessageRequest, false); err != nil {
			return err
		}
	}

	return nil
}
