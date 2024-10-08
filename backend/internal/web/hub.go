package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	/**
	websocketUpgrader is used to upgrade incoming HTTP requests into a persistent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ErrorMessageTypeNotSupported = errors.New("GameMessage type not supported")
)

// Hub is used to hold references to all Clients Registered, and Broadcasting etc
type Hub struct {
	sync.RWMutex        // lock state before editing clients
	clients             ClientList
	gameMessageHandlers map[string]GameMessageHandler
}

func NewHub() *Hub {
	m := &Hub{
		clients:             make(ClientList),
		gameMessageHandlers: make(map[string]GameMessageHandler),
	}
	m.setupGameMessageHandlers()
	return m
}

// todo that's pretty much hardcoded, we should make it more configurable
// setupGameMessageHandlers configures and adds all gameMessageHandlers
func (h *Hub) setupGameMessageHandlers() {
	h.gameMessageHandlers[ChatMessageType] = ChatMessageHandler
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
func sendWelcomeMessage(client *Client, manager *Hub) error {

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
	if len(manager.clients) > 4 {
		return fmt.Errorf("max number of 4 players reached")
	}

	index := len(manager.clients) - 1
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
