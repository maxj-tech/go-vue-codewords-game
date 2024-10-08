package web

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChatMessagePayload struct {
	Text string `json:"text"`
	From string `json:"from"`
}

// ChatMessageWithSentTimeStamp is returned when responding to chat-message
type ChatMessageWithSentTimeStamp struct {
	ChatMessagePayload
	Sent time.Time `json:"sent"`
}

// ChatMessageHandler will send out a message to all other participants in the chat
func ChatMessageHandler(gameMessage GameMessage, c *Client) error {
	// Marshal Payload into wanted format
	var chatMessagePayload ChatMessagePayload
	if err := json.Unmarshal(gameMessage.Payload, &chatMessagePayload); err != nil {
		return fmt.Errorf("ChatMessageHandler(): bad payload in request: %v", err)
	}

	// Prepare an Outgoing Message to others
	var broadcastChatMessage ChatMessageWithSentTimeStamp

	broadcastChatMessage.Sent = time.Now()
	broadcastChatMessage.Text = chatMessagePayload.Text
	broadcastChatMessage.From = chatMessagePayload.From

	data, err := json.Marshal(broadcastChatMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent GameMessage
	outgoingEvent.Payload = data
	outgoingEvent.Type = ChatMessageType
	// Broadcast to all other Clients
	for client := range c.hub.clients {
		client.egress <- outgoingEvent
	}

	return nil

}
