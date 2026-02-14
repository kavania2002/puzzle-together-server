package lib

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

// NewManager creates and returns a Manager with an initialized client list and handler map.
// The returned Manager has its default event handlers registered.
func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		return handler(event, c)
	}

	return ErrEventNotSupported
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust as needed
	},
}

func (m *Manager) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to Websocket:", err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessage()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		delete(m.clients, client)
		client.connection.Close()
		close(client.egress)
	}
}

func (m *Manager) broadcastMessage(event Event) {
	m.RLock()
	defer m.RUnlock()

	for client := range m.clients {
		client.egress <- event
	}
}
