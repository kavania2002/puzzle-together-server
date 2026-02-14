package lib

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	// To know which game room id the client belongs to
	gameRoom string

	// egress used to avoid concurrent writes to websocket
	egress chan Event

	done chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	closeOnce sync.Once
}

// NewClient creates a Client that wraps the provided WebSocket connection and manager.
// The returned Client has an internal buffered egress channel (capacity 256) for outgoing events.
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event, 256),
		done:       make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (c *Client) readMessages() {
	defer c.close()

	for {
		select {
		case <-c.done:
			return
		default:
			messageType, payload, err := c.connection.ReadMessage()

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				return
			}

			if messageType != websocket.TextMessage {
				log.Printf("Unexpected message type: %d", messageType)
				continue
			}

			var request Event
			if err := json.Unmarshal(payload, &request); err != nil {
				log.Printf("Erro unmarshalling message %v", err)
				return
			}

			if err := c.manager.routeEvent(request, c); err != nil {
				log.Println("Error handling message: ", err)
			}
		}
	}
}

func (c *Client) writeMessage() {
	defer c.close()

	for {
		select {
		case <-c.done:

			// Send close message and exit
			if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
				log.Println("Error sending close message: ", err)
			}
			return

		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection Closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// close handles cleanup of client resources in a thread-safe, idempotent manner.
// It uses sync.Once to ensure cleanup happens exactly once, even if called from
// multiple goroutines.
func (c *Client) close() {
	c.closeOnce.Do(func() {
		c.cancel()
		close(c.done)

		c.manager.removeClient(c)
	})
}
