package main

import (
	"github.com/maxj-tech/go-vue-codewords-game/backend/internal/web"
	"log"
	"net/http"
)

func main() {
	setupAPI()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {
	hub := web.NewHub()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		web.ServeWebsocket(hub, w, r)
	})
}
