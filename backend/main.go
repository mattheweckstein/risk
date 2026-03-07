package main

import (
	"log"
	"net/http"

	"github.com/mattheweckstein/risk/backend/api"
	"github.com/mattheweckstein/risk/backend/storage"
)

func main() {
	store := storage.NewStore("games.json")

	games, err := store.LoadAll()
	if err != nil {
		log.Fatalf("Failed to load saved games: %v", err)
	}
	if len(games) > 0 {
		log.Printf("Loaded %d saved game(s) from disk", len(games))
	}

	server := api.NewServer(store, games)
	router := server.Router()

	log.Println("Risk game server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
