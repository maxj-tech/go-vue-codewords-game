package web

import (
	"encoding/json"
)

// GameMessageType is a string that represents the type of message sent over the websocket
type GameMessageType string

const (
	ChatMessage GameMessageType = "chat-message"
)

// GameMessageHandler is a function signature to affect messages on the socket and triggered depending on the type
type GameMessageHandler func(message GameMessage, c *Client) error

// GameMessageHandlers is a map that associates message types with their handlers
type GameMessageHandlers map[GameMessageType]GameMessageHandler

// GameMessage is sent / received over the websocket
type GameMessage struct {
	// Type is the message type sent
	Type GameMessageType `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"` // todo narrow down the type later
}
