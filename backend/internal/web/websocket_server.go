package web

import (
	"log"
	"net/http"
	"time"
)

const (
	DefaultPongWait = 60 * time.Second
)

var DefaultClientConfig = ClientConfig{
	PongWait:     DefaultPongWait,
	PingInterval: (DefaultPongWait * 9) / 10,
}

// ServeWS is a HTTP Handler that the has the Hub that allows connections
func ServeWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {

	log.Println("ServeWebsocket: New connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// todo 1st limit number of clients,
	// todo 2nd could we avoid calling newClient here
	client := newClient(conn, hub, DefaultClientConfig)
	hub.addClient(client)

	// fixme feels wrong here
	if err := sendWelcomeMessage(client, hub); err != nil {
		log.Println("serveWS(): failed to send welcome message: ", err)
		return
	}
	go client.readMessages()
	go client.writeMessages()
}
