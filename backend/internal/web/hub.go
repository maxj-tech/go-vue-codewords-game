package web

/**
 * Hub manages the active clients and their connections, allowing
 * broadcasting and routing of messages between them.
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
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

// setupGameMessageHandlers allows to set up the game message handlers. Only Once!
func (h *Hub) setupOnce(handlers GameMessageHandlers) error {
	h.Lock()
	defer h.Unlock()

	if handlers == nil || len(handlers) == 0 {
		return errors.New("invalid GameMessageHandlers: must not be nil or empty")
	}

	if len(h.gameMessageHandlers) > 0 {
		log.Info("hub.setupOnce(): GameMessageHandlers already set up")
		return nil
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

func (h *Hub) addClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	log.Debug("hub.addClient(): Adding client")
	h.clients[client] = true
}

func (h *Hub) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clients[client]; ok {
		log.Debug("hub.removeClient(): Removing client")
		close(client.egress)
		delete(h.clients, client)
	} else {
		log.Warn("hub.removeClient(): Client not found")
	}
}

// todo not sure if this is the right place for this function
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
	log.Debug("hub.sendWelcomeMessage(): Sending welcome message: ", string(data))
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	gameMessage := GameMessage{
		Type:    "welcome-message",
		Payload: data,
	}

	if err := client.connection.WriteJSON(gameMessage); err != nil {
		log.Error("serveWS(): failed to send welcome message: ", err)
	}
	return nil
}
