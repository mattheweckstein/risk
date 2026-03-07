package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/mattheweckstein/risk/backend/ai"
	"github.com/mattheweckstein/risk/backend/game"
	"github.com/mattheweckstein/risk/backend/models"
)

// Server holds the API state.
type Server struct {
	engine *game.GameEngine
	games  map[string]*models.GameState
	mu     sync.RWMutex
}

// NewServer creates a new API server.
func NewServer() *Server {
	return &Server{
		engine: game.NewGameEngine(),
		games:  make(map[string]*models.GameState),
	}
}

// Router returns a configured chi router with all API routes.
func (s *Server) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/game", func(r chi.Router) {
		r.Post("/new", s.handleNewGame)
		r.Get("/{id}", s.handleGetGame)
		r.Post("/{id}/place", s.handlePlace)
		r.Post("/{id}/attack", s.handleAttack)
		r.Post("/{id}/attack/move", s.handleMoveAfterConquest)
		r.Post("/{id}/fortify", s.handleFortify)
		r.Post("/{id}/end-phase", s.handleEndPhase)
		r.Post("/{id}/cards/trade", s.handleTradeCards)
		r.Get("/{id}/ai-turn", s.handleAITurn)
	})

	return r
}

func (s *Server) handleNewGame(w http.ResponseWriter, r *http.Request) {
	var req models.NewGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayerName == "" {
		req.PlayerName = "Player"
	}
	if req.AICount < 1 {
		req.AICount = 1
	}
	if req.AICount > 3 {
		req.AICount = 3
	}

	state := s.engine.NewGame(req.PlayerName, req.AICount, req.AINames)
	state.ID = uuid.New().String()

	s.mu.Lock()
	s.games[state.ID] = state
	s.mu.Unlock()

	log.Printf("New game created: %s with %d players", state.ID, len(state.Players))
	writeJSON(w, http.StatusCreated, state)
}

func (s *Server) handleGetGame(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handlePlace(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req models.PlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.mu.Lock()
	err := s.engine.PlaceTroops(state, req.Territory, req.Troops)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleAttack(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req models.AttackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.mu.Lock()
	_, err := s.engine.Attack(state, req.From, req.To, req.AttackerDice)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleMoveAfterConquest(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req models.MoveAfterConquestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.mu.Lock()
	err := s.engine.MoveAfterConquest(state, req.Troops)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleFortify(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req models.FortifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.mu.Lock()
	err := s.engine.Fortify(state, req.From, req.To, req.Troops)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleEndPhase(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	s.mu.Lock()
	err := s.engine.EndPhase(state)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleTradeCards(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	var req models.CardTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.mu.Lock()
	err := s.engine.TradeCards(state, req.CardIndices)
	s.mu.Unlock()

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func (s *Server) handleAITurn(w http.ResponseWriter, r *http.Request) {
	state, ok := s.getGame(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check that the current player is an AI
	if !s.engine.IsCurrentPlayerAI(state) {
		writeError(w, http.StatusBadRequest, "current player is not AI")
		return
	}

	switch state.Phase {
	case models.PhaseSetup:
		// During setup, keep placing for all consecutive AI players until
		// it's the human's turn (or setup ends). This avoids many round-trips.
		for state.Phase == models.PhaseSetup && s.engine.IsCurrentPlayerAI(state) {
			territoryID := ai.PlaceSetupTroop(state, state.CurrentPlayer)
			if err := s.engine.PlaceTroops(state, territoryID, 1); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

	case models.PhasePlace, models.PhaseAttack, models.PhaseFortify:
		// Full AI turn: the AI package handles placement, attack, and fortify
		// by directly mutating state. We set up troops to deploy first.
		if state.Phase == models.PhasePlace {
			// AI handles its own troop calculation and placement
			if err := ai.ExecuteTurn(state, state.CurrentPlayer); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Reset TroopsToDeploy since AI placed directly
			state.TroopsToDeploy = 0
		} else {
			if err := ai.ExecuteTurn(state, state.CurrentPlayer); err != nil {
				writeError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Check for win condition
		if state.Phase == models.PhaseEnded {
			break
		}

		// Clear any pending conquest from AI (AI resolves its own troop movements)
		state.PendingConquest = nil

		// Advance to next player's turn
		state.Phase = models.PhaseFortify
		if err := s.engine.EndPhase(state); err != nil {
			log.Printf("EndPhase after AI turn failed: %v", err)
		}
	}

	writeJSON(w, http.StatusOK, state)
}

// Helper methods

func (s *Server) getGame(id string) (*models.GameState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.games[id]
	return state, ok
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
