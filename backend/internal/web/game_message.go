package web

import (
	"encoding/json"
)

// GameMessage is sent / received over the websocket
type GameMessage struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"` // todo narrow down the type later
}

// GameMessageHandler is a function signature to affect messages on the socket and triggered depending on the type
type GameMessageHandler func(message GameMessage, c *Client) error

const (
	ChatMessageType = "chat-message"
)
