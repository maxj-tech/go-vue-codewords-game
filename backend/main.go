package main

import (
	"fmt"
	"github.com/maxj-tech/go-vue-codewords-game/backend/internal/web"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const DEFAULT_LOG_LEVEL = log.InfoLevel

func setupLogger() error {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = DEFAULT_LOG_LEVEL.String()
		log.Infof("Using default log level %s. Set LOG_LEVEL in the environment, e.g., "+
			"\"LOG_LEVEL=trace go run .\"", logLevel)
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}

	log.Debugf("Setting log level to %s", level)
	log.SetLevel(level)

	if level == log.DebugLevel || level == log.TraceLevel {
		log.SetReportCaller(true)
	}

	return nil
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
