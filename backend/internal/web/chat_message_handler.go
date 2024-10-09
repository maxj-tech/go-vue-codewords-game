package web

/**
 * ChatMessageHandler routes an incoming chat message to all connected clients in the hub.
 * It parses the message payload, creates a new GameMessage with the updated timestamp,
 * and broadcasts it to all clients.
 */

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChatMessagePayload struct {
	From string    `json:"from"`
	Text string    `json:"text"`
	Sent time.Time `json:"sent"`
}

// ChatMessageHandler will send out a message to all other participants in the chat
func ChatMessageHandler(gameMessage GameMessage, c *Client) error {
	incomingPayload, err := parseChatMessage(gameMessage)
	if err != nil {
		return err
	}

	outgoingMessage, err := createGameMessage(incomingPayload)
	if err != nil {
		return err
	}

	for client := range c.hub.clients {
		client.egress <- outgoingMessage
	}
	return nil
}

func createGameMessage(incomingPayload ChatMessagePayload) (GameMessage, error) {
	chatMessagePayload := ChatMessagePayload{
		From: incomingPayload.From,
		Text: incomingPayload.Text,
		Sent: time.Now(),
	}

	marshalledChatMessagePayload, err := json.Marshal(chatMessagePayload)
	if err != nil {
		return GameMessage{},
			fmt.Errorf("createGameMessage(): failed to marshal chat message payload: %v", err)
	}

	gameMessage := GameMessage{
		Type:    ChatMessage,
		Payload: marshalledChatMessagePayload,
	}
	return gameMessage, nil
}

func parseChatMessage(gameMessage GameMessage) (ChatMessagePayload, error) {
	var incomingPayload ChatMessagePayload
	if err := json.Unmarshal(gameMessage.Payload, &incomingPayload); err != nil {
		return ChatMessagePayload{}, fmt.Errorf("ChatMessageHandler(): bad payload in request: %v", err)
	}
	return incomingPayload, nil
}
