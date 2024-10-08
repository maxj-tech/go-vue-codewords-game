package main

import (
	"github.com/maxj-tech/go-vue-codewords-game/backend/internal/web"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const DefaultLogLevel = log.InfoLevel

func setupLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = DefaultLogLevel.String()
		log.Infof("Using default log level %s. Set LOG_LEVEL in the environment, e.g., "+
			"\"LOG_LEVEL=trace go run .\"", logLevel)
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("Invalid log level %s: %v", logLevel, err)
		level = DefaultLogLevel
	}

	log.Debugf("Setting log level to %s", level)
	log.SetLevel(level)

	if level == log.DebugLevel || level == log.TraceLevel {
		log.SetReportCaller(true)
	}
}

func main() {
	setupLogger()

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
