package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager

	// egress is used to avoid concurrent writes on the WebSocket
	egress chan []byte
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

// readMessages will start the client to read messages and handle them appropriately.
// This is supposed to be run as a goroutine
func (c *Client) readMessages() {
	defer func() { // Graceful Close the Connection once this function is done
		c.manager.removeClient(c)
	}()

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will receive an error here
			// We only want to log strange errors, but not simple disconnects
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		log.Printf("MessageType: %d, Payload: %s \n", messageType, string(payload))

		// fixme remove this - just to test that WriteMessages works as intended
		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
		}
	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {

	defer func() { // Graceful close if this triggers a closing
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// ok will be false if the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}
				// Return to close the goroutine
				return
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		}

	}
}
