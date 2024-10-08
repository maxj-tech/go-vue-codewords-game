package web

/**
 * Hub manages clients and their connections todo
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	MaxClients = 4 // fixme just temporary
)

var (
	ErrorMessageTypeNotSupported = errors.New("GameMessage type not supported")

	DefaultGameMessageHandlers = GameMessageHandlers{
		ChatMessage: ChatMessageHandler,
	}
)

type Hub struct {
	sync.RWMutex        // lock state before editing clients
	clients             ClientList
	gameMessageHandlers GameMessageHandlers
}

func NewHub() *Hub {
	m := &Hub{
		clients:             make(ClientList),
		gameMessageHandlers: make(GameMessageHandlers),
	}
	return m
}

// setupGameMessageHandlers configures and adds all gameMessageHandlers
func (h *Hub) setup(handlers GameMessageHandlers) error {
	// todo lock needed here?
	if handlers == nil || len(handlers) == 0 {
		return errors.New("invalid GameMessageHandlers: must not be nil or empty")
	}

	if len(h.gameMessageHandlers) > 0 {
		return errors.New("game message handlers are already set")
	}

	h.gameMessageHandlers = handlers
	return nil
}

// makes sure the correct game message goes into the correct game message handler
func (h *Hub) route(gameMessage GameMessage, c *Client) error {
	if handler, ok := h.gameMessageHandlers[gameMessage.Type]; ok {
		if err := handler(gameMessage, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrorMessageTypeNotSupported
	}
}

// todo tidy up
func (h *Hub) sendWelcomeMessage(client *Client) error {

	welcomeMessageData := []struct {
		Name string
		Team string
		Role string
	}{
		{"Ermittler1", "TeamRot", "Ermittler"},
		{"Chef1", "TeamRot", "Chef"},
		{"Ermittler2", "TeamBlau", "Ermittler"},
		{"Chef2", "TeamBlau", "Chef"},
	}
	if len(h.clients) > MaxClients {
		return fmt.Errorf("max number of %d players reached", MaxClients)
	}

	index := len(h.clients) - 1
	welcomeMessage := welcomeMessageData[index]

	data, err := json.Marshal(welcomeMessage)
	log.Println("manager.sendWelcomeMessage(): Sending welcome message: ", string(data))
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	gameMessage := GameMessage{
		Type:    "welcome-message",
		Payload: data,
	}

	if err := client.connection.WriteJSON(gameMessage); err != nil {
		log.Println("serveWS(): failed to send welcome message: ", err)
	}
	return nil
}

func (h *Hub) addClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	log.Println("manager.addClient(): Adding client")
	h.clients[client] = true
}

func (h *Hub) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clients[client]; ok {
		client.connection.Close()
		delete(h.clients, client)
	}
}
