package lib

import (
	"encoding/json"
	"fmt"
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
func SendMessageHandler(event Event, c *Client) error {
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatEvent.Message
	broadMessage.From = chatEvent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	c.manager.broadcastMessage(outgoingEvent)

	return nil
}
