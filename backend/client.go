package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	pongWait = 10 * time.Second // how long we will await a pong response from client

	// pingInterval has to be less than pongWait, We cant multiply by 0.9 to get 90% of time
	// Because that can make decimals, so instead *9 / 10 to get 90%
	// The reason why it has to be less than PingRequency is becuase otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager

	// egress is used to avoid concurrent writes on the WebSocket connection
	egress chan GameMessage
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan GameMessage),
	}
}

// readMessages will start the client to read messages and handle them appropriately.
// This is supposed to be run as a goroutine
func (c *Client) readMessages() {
	defer func() { // Graceful Close the Connection once this function is done
		log.Println("client.readMessages(): Closing connection")
		c.manager.removeClient(c)
	}()

	c.connection.SetReadLimit(256) // maximum message size is 256 bytes

	// Configure wait time for Pong respons: current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will receive an error here
			// We only want to log strange errors, but not simple disconnects
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("client.readMessages(): error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}

		log.Printf("client.readMessages(): MessageType: %d, Payload: %s \n", messageType, string(payload))

		// Marshal incoming data into a GameMessage
		var request GameMessage
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("client.readMessages(): error marshalling message: %v", err)
			break // fixme better avoid Breaking the connection here
		}

		// fixme remove this (just to test if the message is broadcasted to all clients)
		//for client := range c.manager.clients {
		//	client.egress <- request
		//}

		// Route the GameMessage
		if err := c.manager.route(request, c); err != nil {
			log.Println("client.readMessages(): Error handling GameMessage: ", err)
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	log.Println("client.pongHandler(): pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait)) // Current time + Pong Wait time
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {

	ticker := time.NewTicker(pingInterval) // ticker that will send a ping every pingInterval

	defer func() { // Graceful close if this triggers a closing
		ticker.Stop()
		log.Println("client.writeMessages(): Closing connection")
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
					log.Println("client.writeMessages(): connection closed: ", err)
				}
				return
			}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Println("client.writeMessages(): error marshalling", err)
				return // closes the connection, should we really
			}
			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("client.writeMessages(): failed to send TextMessage", err)
			}
			log.Println("client.writeMessages(): marshalled message sent", msg)

		case <-ticker.C:
			log.Println("client.writeMessages(): ping")
			// Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("client.writeMessages(): failed to send PingMessage", err)
				return // return to break this goroutine triggering cleanup
			}
		}

	}
}
