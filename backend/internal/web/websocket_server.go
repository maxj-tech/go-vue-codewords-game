package web

/**
 * ServeWebsocket handles incoming websocket requests, upgrades the HTTP connection
 * to a websocket connection, and initializes a new client.
 * It also starts the read and write goroutines for the client.
 */

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	// DefaultPongWait is the default time to wait for a pong message.
	DefaultPongWait = 60 * time.Second
)

var (
	// DefaultClientConfig is the default configuration for a client.
	DefaultClientConfig = ClientConfig{
		PongWait:     DefaultPongWait,
		PingInterval: (DefaultPongWait * 9) / 10, // shorter than pongWait!
	}

	// websocketUpgrader is used to upgrade an HTTP connection to a websocket connection.
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	log.Info("ServeWebsocket: New connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	err = hub.setupOnce(DefaultGameMessageHandlers)
	if err != nil {
		log.Error(err)
	}

	// todo 1st limit number of clients,
	// todo 2nd could we avoid calling newClient here
	client := newClient(conn, hub, DefaultClientConfig)
	hub.addClient(client)

	// fixme feels wrong here
	if err := hub.sendWelcomeMessage(client); err != nil {
		log.Error("serveWS(): failed to send welcome message: ", err)
		return
	}
	client.startReadMessagesRoutine()
	client.startWriteMessagesRoutine()
}
