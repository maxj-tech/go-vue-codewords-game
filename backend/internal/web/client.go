package web

/**
 * Client represents a websocket connection tied to a player and their interaction with the game server.
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

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	hub        *Hub
	egress     chan GameMessage
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

func (c *Client) startReadMessagesRoutine() {
	go func() {
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
				log.WithError(err).Error("client.startReadMessagesRoutine(): Error handling GameMessage")
			}
		}
	}()
}

func (c *Client) startWriteMessagesRoutine() {
	go func() {
		ticker := time.NewTicker(c.config.PingInterval) // ticker that will send a ping every pingInterval
		defer func() {
			ticker.Stop()
			c.cleanup()
		}()

		for {
			select {
			case gameMessage, ok := <-c.egress: // ok will be false if the egress channel is closed
				if !ok {
					// tell to frontend that we are closing the connection
					if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
						log.WithError(err).Warn("client.startWriteMessagesRoutine(): connection closed")
					}
					return
				}

				if ok := c.send(gameMessage); !ok {
					return
				}

			case <-ticker.C:
				if ok := c.sendPing(); !ok {
					return
				}
			}
		}
	}()
}

func (c *Client) send(gameMessage GameMessage) (ok bool) {
	msg, err := json.Marshal(gameMessage)
	if err != nil {
		log.WithError(err).Errorf("client.send(): error marshalling GameMessage %s", gameMessage)
		return false
	}

	if err := c.connection.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.WithError(err).Error("client.send(): failed to send TextMessage")
		return false
	}
	log.Debug("client.startWriteMessagesRoutine(): marshalled gameMessage sent", msg)
	return true
}

func (c *Client) sendPing() (ok bool) {
	log.Debug("client: ping")
	if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		log.WithError(err).Error("client.startWriteMessagesRoutine(): failed to send PingMessage")
		return false
	}
	return true
}

func (c *Client) readMessage() (messageType int, payload []byte, err error) {
	messageType, payload, err = c.connection.ReadMessage()
	if err != nil {
		log.WithError(err).Error("client.startReadMessagesRoutine(): error reading message")
	}
	log.Debugf("client.startReadMessagesRoutine(): MessageType: %d, Payload: %s \n", messageType, string(payload))
	return messageType, payload, err
}

func (c *Client) makeGameMessage(payload []byte) (GameMessage, error) {
	var gameMessage GameMessage
	if err := json.Unmarshal(payload, &gameMessage); err != nil {
		log.WithError(err).Error("client.startReadMessagesRoutine(): error unmarshalling message")
		return GameMessage{}, err
	}
	return gameMessage, nil
}

func (c *Client) cleanup() {
	log.Debug("client.cleanup(): Closing connection")
	c.hub.removeClient(c) // todo error handling?
	if err := c.connection.Close(); err != nil {
		log.WithError(err).Error("client.cleanup(): error closing connection")
	}
}

func (c *Client) configureConnection(config ClientReadConnectionConfig) {
	c.connection.SetReadLimit(config.ReadLimit)
	c.connection.SetPongHandler(func(pongMessage string) error {
		log.Debug("client: pong")
		return c.connection.SetReadDeadline(time.Now().Add(config.ReadDeadlineDelta))
	})

}
