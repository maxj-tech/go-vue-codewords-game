package web

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type ClientConfig struct {
	PongWait     time.Duration
	PingInterval time.Duration
}

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	connection *websocket.Conn
	hub        *Hub
	egress     chan GameMessage // todo "egress"  suggests outgoing messages only. check if that's the case
	config     ClientConfig
}

func newClient(conn *websocket.Conn, hub *Hub, config ClientConfig) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
		egress:     make(chan GameMessage),
		config:     config,
	}
}

// todo if this is supposed to be run as a goroutine, why not internalize the goroutine here?
func (c *Client) readMessages() {
	defer func() { // Graceful Close the Connection once this function is done
		log.Println("client.readMessages(): Closing connection")
		c.hub.removeClient(c)
	}()

	c.connection.SetReadLimit(256) // maximum message size is 256 bytes

	// Configure wait time for Pong respons: current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(c.config.PongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	for {
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

		// Route the GameMessage
		if err := c.hub.route(request, c); err != nil {
			log.Println("client.readMessages(): Error handling GameMessage: ", err)
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	//log.Println("client.pongHandler(): pong")	// todo  set log levels
	return c.connection.SetReadDeadline(time.Now().Add(c.config.PongWait)) // Current time + Pong Wait time
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {

	ticker := time.NewTicker(c.config.PingInterval) // ticker that will send a ping every pingInterval

	defer func() { // Graceful close if this triggers a closing
		ticker.Stop()
		log.Println("client.writeMessages(): Closing connection")
		c.hub.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress: // ok will be false if the egress channel is closed
			if !ok {
				// tell to frontend that we are closing the connection
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("client.writeMessages(): connection closed: ", err)
				}
				return
			}

			msg, err := json.Marshal(message)
			if err != nil {
				log.Println("client.writeMessages(): error marshalling", err)
				return // closes the connection, should we really fail here
			}

			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("client.writeMessages(): failed to send TextMessage", err)
			}
			log.Println("client.writeMessages(): marshalled message sent", msg)

		case <-ticker.C:
			//log.Println("client.writeMessages(): ping")
			// Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("client.writeMessages(): failed to send PingMessage", err)
				return
			}
		}

	}
}
