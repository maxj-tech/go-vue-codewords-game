package web

/**
 * a websocket client. One Client per Connection, i.e. Player

 */

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

type ClientReadConnectionConfig struct {
	ReadLimit         int64
	ReadDeadlineDelta time.Duration
}

var DefaultClientReadConnectionConfig = ClientReadConnectionConfig{
	ReadLimit:         256,
	ReadDeadlineDelta: 60 * time.Second,
}

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
	defer c.cleanup()

	c.configureConnection(DefaultClientReadConnectionConfig)

	for {
		_, payload, err := c.readMessage()
		if err != nil {
			break
		}

		gameMessage, err := c.makeGameMessage(payload)
		if err != nil {
			break
		}

		if err := c.hub.route(gameMessage, c); err != nil {
			log.Error("client.readMessages(): Error handling GameMessage: ", err)
		}

	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	ticker := time.NewTicker(c.config.PingInterval) // ticker that will send a ping every pingInterval
	defer func() {
		ticker.Stop()
		c.cleanup()
	}()

	for {
		select {
		case message, ok := <-c.egress: // ok will be false if the egress channel is closed
			if !ok {
				// tell to frontend that we are closing the connection
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Warn("client.writeMessages(): connection closed: ", err)
				}
				return
			}

			msg, err := json.Marshal(message)
			if err != nil {
				log.Fatal("client.writeMessages(): error marshalling", err)
				break
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error("client.writeMessages(): failed to send TextMessage", err)
			}
			log.Debug("client.writeMessages(): marshalled message sent", msg)

		case <-ticker.C:
			log.Debug("client.writeMessages(): ping")
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Error("client.writeMessages(): failed to send PingMessage", err)
				return
			}
		}
	}
}

func (c *Client) readMessage() (messageType int, payload []byte, error error) {
	messageType, payload, err := c.connection.ReadMessage()
	if err == nil {
		log.Debugf("client.readMessages(): MessageType: %d, Payload: %s \n", messageType, string(payload))
	} else {
		log.Error("client.readMessages(): error reading message: %v", err)
	}
	return messageType, payload, err
}

func (c *Client) makeGameMessage(payload []byte) (GameMessage, error) {
	var gameMessage GameMessage
	if err := json.Unmarshal(payload, &gameMessage); err != nil {
		log.Printf("client.readMessages(): error marshalling message: %v", err)
		return GameMessage{}, err
	}
	return gameMessage, nil
}

func (c *Client) cleanup() {
	log.Debug("client.cleanup(): Closing connection")
	c.hub.removeClient(c)
	if err := c.connection.Close(); err != nil {
		log.Error("client.cleanup(): error closing connection: ", err)
	}
}

func (c *Client) configureConnection(config ClientReadConnectionConfig) {
	c.connection.SetReadLimit(config.ReadLimit)
	c.connection.SetPongHandler(func(pongMessage string) error {
		log.Debug("client: pong")
		return c.connection.SetReadDeadline(time.Now().Add(config.ReadDeadlineDelta))
	})

}
