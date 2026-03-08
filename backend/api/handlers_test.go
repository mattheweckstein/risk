package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
	"github.com/mattheweckstein/risk/backend/storage"
)

// newTestServer creates a test server with an in-memory store.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	dir := t.TempDir()
	store := storage.NewStore(dir + "/test_games.json")
	games := make(map[string]*models.GameState)
	return NewServer(store, games)
}

// createTestGame creates a game via the API and returns the game state.
func createTestGame(t *testing.T, srv *Server) *models.GameState {
	t.Helper()
	body := `{"playerName":"TestPlayer","aiCount":1,"aiNames":["TestBot"]}`
	req := httptest.NewRequest("POST", "/api/game/new", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var state models.GameState
	if err := json.NewDecoder(w.Body).Decode(&state); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return &state
}

func TestNewGameEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	if state.ID == "" {
		t.Error("expected non-empty game ID")
	}
	if len(state.Players) != 2 {
		t.Errorf("expected 2 players, got %d", len(state.Players))
	}
	if state.Players[0].Name != "TestPlayer" {
		t.Errorf("expected player name TestPlayer, got %q", state.Players[0].Name)
	}
	if state.Phase != models.PhaseSetup {
		t.Errorf("expected setup phase, got %s", state.Phase)
	}
	if len(state.Territories) != 42 {
		t.Errorf("expected 42 territories, got %d", len(state.Territories))
	}
}

func TestNewGameEndpointDefaults(t *testing.T) {
	srv := newTestServer(t)

	// Empty body should use defaults
	body := `{}`
	req := httptest.NewRequest("POST", "/api/game/new", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var state models.GameState
	json.NewDecoder(w.Body).Decode(&state)

	if state.Players[0].Name != "Player" {
		t.Errorf("expected default name 'Player', got %q", state.Players[0].Name)
	}
}

func TestGetGameEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	req := httptest.NewRequest("GET", "/api/game/"+state.ID, nil)
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if got.ID != state.ID {
		t.Errorf("expected game ID %s, got %s", state.ID, got.ID)
	}
}

func TestGetGameNotFound(t *testing.T) {
	srv := newTestServer(t)

	req := httptest.NewRequest("GET", "/api/game/nonexistent-id", nil)
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestPlaceEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	// Skip setup to place phase
	srv.mu.Lock()
	g := srv.games[state.ID]
	g.Phase = models.PhasePlace
	g.Turn = 1
	g.TroopsToDeploy = 5
	for id, terr := range g.Territories {
		terr.Troops = 3
		g.Territories[id] = terr
	}
	srv.mu.Unlock()

	// Find a territory owned by player_0
	var ownedID string
	for id, terr := range g.Territories {
		if terr.Owner == "player_0" {
			ownedID = id
			break
		}
	}

	body, _ := json.Marshal(models.PlaceRequest{Territory: ownedID, Troops: 2})
	req := httptest.NewRequest("POST", "/api/game/"+state.ID+"/place", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if got.Territories[ownedID].Troops != 5 {
		t.Errorf("expected 5 troops after placing 2 on 3, got %d", got.Territories[ownedID].Troops)
	}
}

func TestAttackEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	// Set up attack phase with strong attacker
	srv.mu.Lock()
	g := srv.games[state.ID]
	g.Phase = models.PhaseAttack
	g.Turn = 1
	g.TroopsToDeploy = 0
	for id, terr := range g.Territories {
		terr.Troops = 3
		g.Territories[id] = terr
	}

	// Find attack pair
	var fromID, toID string
	for id, terr := range g.Territories {
		if terr.Owner == "player_0" {
			for _, nID := range terr.Neighbors {
				if n, ok := g.Territories[nID]; ok && n.Owner != "player_0" {
					fromID = id
					toID = nID
					break
				}
			}
		}
		if fromID != "" {
			break
		}
	}

	// Give attacker lots of troops
	f := g.Territories[fromID]
	f.Troops = 20
	g.Territories[fromID] = f
	srv.mu.Unlock()

	body, _ := json.Marshal(models.AttackRequest{From: fromID, To: toID, AttackerDice: 3})
	req := httptest.NewRequest("POST", "/api/game/"+state.ID+"/attack", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if got.LastAttackResult == nil {
		t.Error("expected attack result in response")
	}
}

func TestEndPhaseEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	// Set to attack phase
	srv.mu.Lock()
	g := srv.games[state.ID]
	g.Phase = models.PhaseAttack
	g.Turn = 1
	g.TroopsToDeploy = 0
	for id, terr := range g.Territories {
		terr.Troops = 3
		g.Territories[id] = terr
	}
	srv.mu.Unlock()

	req := httptest.NewRequest("POST", "/api/game/"+state.ID+"/end-phase", nil)
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if got.Phase != models.PhaseFortify {
		t.Errorf("expected fortify phase after ending attack, got %s", got.Phase)
	}
}

func TestSurrenderEndpoint(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	// Set to some active phase
	srv.mu.Lock()
	g := srv.games[state.ID]
	g.Phase = models.PhaseAttack
	g.Turn = 1
	for id, terr := range g.Territories {
		terr.Troops = 3
		g.Territories[id] = terr
	}
	srv.mu.Unlock()

	req := httptest.NewRequest("POST", "/api/game/"+state.ID+"/surrender", nil)
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if got.Phase != models.PhaseEnded {
		t.Errorf("expected ended phase after surrender, got %s", got.Phase)
	}
	if got.Winner == "" {
		t.Error("expected a winner after surrender")
	}

	// Human should be dead
	for _, p := range got.Players {
		if !p.IsAI && p.IsAlive {
			t.Error("expected human to be dead after surrender")
		}
	}
}

func TestFreeFortifySetting(t *testing.T) {
	srv := newTestServer(t)

	// Create game with freeFortify enabled
	body := `{"playerName":"TestPlayer","aiCount":1,"freeFortify":true}`
	req := httptest.NewRequest("POST", "/api/game/new", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var state models.GameState
	json.NewDecoder(w.Body).Decode(&state)

	if !state.FreeFortify {
		t.Error("expected FreeFortify to be true")
	}

	// Fetch the game back and verify persistence
	req = httptest.NewRequest("GET", "/api/game/"+state.ID, nil)
	w = httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	var got models.GameState
	json.NewDecoder(w.Body).Decode(&got)

	if !got.FreeFortify {
		t.Error("expected FreeFortify to persist through get")
	}
}

func TestPlaceEndpointNotFound(t *testing.T) {
	srv := newTestServer(t)

	body := `{"territory":"alaska","troops":1}`
	req := httptest.NewRequest("POST", "/api/game/nonexistent/place", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAttackEndpointBadRequest(t *testing.T) {
	srv := newTestServer(t)
	state := createTestGame(t, srv)

	// Try to attack during setup phase
	body := `{"from":"alaska","to":"kamchatka","attackerDice":1}`
	req := httptest.NewRequest("POST", "/api/game/"+state.ID+"/attack", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.Router().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
