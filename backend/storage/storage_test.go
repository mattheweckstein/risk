package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
)

func newTestGameState() *models.GameState {
	return &models.GameState{
		ID:            "test-123",
		Phase:         models.PhaseAttack,
		Turn:          3,
		CurrentPlayer: "player_0",
		Players: []models.Player{
			{ID: "player_0", Name: "Alice", IsAI: false, Color: "red", Cards: []models.Card{}, IsAlive: true},
			{ID: "player_1", Name: "Bot", IsAI: true, Color: "blue", Cards: []models.Card{{Territory: "alaska", Type: "infantry"}}, IsAlive: true},
		},
		Territories: map[string]models.Territory{
			"alaska": {ID: "alaska", Name: "Alaska", Continent: "north_america", Neighbors: []string{"northwest_territory", "alberta", "kamchatka"}, Owner: "player_0", Troops: 5},
			"brazil": {ID: "brazil", Name: "Brazil", Continent: "south_america", Neighbors: []string{"venezuela", "peru", "argentina", "north_africa"}, Owner: "player_1", Troops: 3},
		},
		Deck:           []models.Card{{Territory: "china", Type: "cavalry"}},
		Log:            []models.LogEntry{{Turn: 1, Player: "player_0", Message: "test"}},
		TroopsToDeploy: 0,
		CardTradeCount: 1,
		FreeFortify:    true,
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "test_games.json")
	store := NewStore(fp)

	games := map[string]*models.GameState{
		"test-123": newTestGameState(),
	}

	err := store.SaveAll(games)
	if err != nil {
		t.Fatalf("SaveAll failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		t.Fatal("expected file to exist after save")
	}

	// Load back
	loaded, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected 1 game, got %d", len(loaded))
	}

	g, ok := loaded["test-123"]
	if !ok {
		t.Fatal("expected game test-123 to exist")
	}

	if g.ID != "test-123" {
		t.Errorf("expected ID test-123, got %s", g.ID)
	}
	if g.Phase != models.PhaseAttack {
		t.Errorf("expected phase attack, got %s", g.Phase)
	}
	if g.Turn != 3 {
		t.Errorf("expected turn 3, got %d", g.Turn)
	}
	if g.CurrentPlayer != "player_0" {
		t.Errorf("expected current player player_0, got %s", g.CurrentPlayer)
	}
	if len(g.Players) != 2 {
		t.Errorf("expected 2 players, got %d", len(g.Players))
	}
	if len(g.Territories) != 2 {
		t.Errorf("expected 2 territories, got %d", len(g.Territories))
	}
	if g.Territories["alaska"].Troops != 5 {
		t.Errorf("expected alaska troops 5, got %d", g.Territories["alaska"].Troops)
	}
	if g.CardTradeCount != 1 {
		t.Errorf("expected card trade count 1, got %d", g.CardTradeCount)
	}
	if !g.FreeFortify {
		t.Error("expected FreeFortify to be true")
	}
	if len(g.Players[1].Cards) != 1 {
		t.Errorf("expected 1 card for player_1, got %d", len(g.Players[1].Cards))
	}
}

func TestLoadNonexistent(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "nonexistent.json")
	store := NewStore(fp)

	games, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll should not error for nonexistent file: %v", err)
	}

	if games == nil {
		t.Fatal("expected non-nil map")
	}

	if len(games) != 0 {
		t.Errorf("expected empty map, got %d entries", len(games))
	}
}

func TestSaveMultipleGames(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "multi_games.json")
	store := NewStore(fp)

	g1 := newTestGameState()
	g2 := newTestGameState()
	g2.ID = "test-456"
	g2.Turn = 10

	games := map[string]*models.GameState{
		"test-123": g1,
		"test-456": g2,
	}

	err := store.SaveAll(games)
	if err != nil {
		t.Fatalf("SaveAll failed: %v", err)
	}

	loaded, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 games, got %d", len(loaded))
	}

	if loaded["test-456"].Turn != 10 {
		t.Errorf("expected turn 10 for test-456, got %d", loaded["test-456"].Turn)
	}
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "atomic_games.json")
	store := NewStore(fp)

	// Save initial state
	games := map[string]*models.GameState{
		"test-123": newTestGameState(),
	}
	if err := store.SaveAll(games); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Save updated state
	games["test-123"].Turn = 99
	if err := store.SaveAll(games); err != nil {
		t.Fatalf("update save failed: %v", err)
	}

	// Temp file should not exist (it should have been renamed)
	tmpPath := fp + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("temp file should not exist after successful save")
	}

	// Load and verify updated data
	loaded, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if loaded["test-123"].Turn != 99 {
		t.Errorf("expected turn 99, got %d", loaded["test-123"].Turn)
	}
}

func TestSaveEmptyMap(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "empty_games.json")
	store := NewStore(fp)

	err := store.SaveAll(map[string]*models.GameState{})
	if err != nil {
		t.Fatalf("SaveAll empty map failed: %v", err)
	}

	loaded, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(loaded) != 0 {
		t.Errorf("expected 0 games, got %d", len(loaded))
	}
}
