package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientList

	// Using a syncMutex here to be able to lock state before editing clients (we could also use channels)
	sync.RWMutex

	// a map of message types to their respective game message handler functions
	gameMessageHandlers map[string]GameMessageHandler
}

// NewManager is used to initalize all the values inside the manager
func NewManager() *Manager {
	m := &Manager{
		clients:             make(ClientList),
		gameMessageHandlers: make(map[string]GameMessageHandler),
	}
	m.setupGameMessageHandlers()
	return m
}

// setupGameMessageHandlers configures and adds all gameMessageHandlers
func (m *Manager) setupGameMessageHandlers() {
	m.gameMessageHandlers[ChatMessageType] = ChatMessageHandler
}

// makes sure the correct game message goes into the correct game message handler
func (m *Manager) route(gameMessage GameMessage, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.gameMessageHandlers[gameMessage.Type]; ok {
		// Execute the handler and return any err
		if err := handler(gameMessage, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrorMessageTypeNotSupported
	}
}

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	log.Println("manager.serveWS(): New connection")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create New Client
	client := NewClient(conn, m)
	// Add the newly created client to the manager
	m.addClient(client)

	// Send Welcome Message
	if err := sendWelcomeMessage(client, m); err != nil {
		log.Println("serveWS(): failed to send welcome message: ", err)
		return
	}

	// Start the read / write processes
	go client.readMessages()
	go client.writeMessages()
}

func sendWelcomeMessage(client *Client, manager *Manager) error {

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

// addClient will add clients to our clientList
func (m *Manager) addClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	log.Println("manager.addClient(): Adding client")
	defer m.Unlock()

	// Add Client
	m.clients[client] = true

}

// removeClient will remove the client and clean up
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
	}
}
