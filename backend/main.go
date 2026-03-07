package main

import (
	"log"
	"net/http"

	"github.com/mattheweckstein/risk/backend/api"
)

func main() {
	server := api.NewServer()
	router := server.Router()

	log.Println("Risk game server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
