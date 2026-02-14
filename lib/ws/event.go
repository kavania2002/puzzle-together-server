package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// SendMessageHandler handles a send_message event by broadcasting a new_message event to all connected clients.
// It unmarshals the incoming payload into a SendMessageEvent, constructs a NewMessageEvent with the current time,
// marshals it to JSON, and sends the resulting Event to each client's egress channel. Returns an error if payload
// unmarshal or message marshal fails.
// Uses RLock for reading the client list and non-blocking sends to prevent deadlocks.
func SendMessageHandler(event Event, c *Client) error {
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bay payload in request: %v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatEvent.Message
	broadMessage.From = chatEvent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadchat message %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	// Use RLock to safely iterate over clients
	c.manager.RLock()
	defer c.manager.RUnlock()

	for client := range c.manager.clients {
		// Non-blocking send to prevent deadlocks
		// If the channel is full or closed, skip this client
		select {
		case client.egress <- outgoingEvent:
			// Message sent successfully
		default:
			// Channel full or client is shutting down, skip
			log.Printf("Failed to send message to client, channel full or closed")
		}
	}

	return nil
}