package main

import (
	"final-project/pkg/api"
	"final-project/pkg/db"
	"log"
	"net/http"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("[ERROR] Database init error: %v", err)
	}

	api.Init()

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := ":7540"
	log.Printf("[Server] Starting on port %s...", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("[ERROR] Server failed: %v", err)
	}
}
