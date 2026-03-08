package game

import (
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
)

func TestSetupPlaceTroop(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Human", 1, []string{"Bot"})

	// Should be in setup phase
	if state.Phase != models.PhaseSetup {
		t.Fatalf("expected setup phase, got %s", state.Phase)
	}

	// Current player should be player_0
	if state.CurrentPlayer != "player_0" {
		t.Fatalf("expected player_0, got %s", state.CurrentPlayer)
	}

	// Place a troop on an owned territory
	ownedID := findOwnedTerritory(state, "player_0")
	if ownedID == "" {
		t.Fatal("no owned territory found")
	}

	origTroops := state.Territories[ownedID].Troops
	err := e.SetupPlaceTroop(state, ownedID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Territories[ownedID].Troops != origTroops+1 {
		t.Errorf("expected troops %d, got %d", origTroops+1, state.Territories[ownedID].Troops)
	}

	// After placing, current player should advance to player_1
	if state.CurrentPlayer != "player_1" {
		t.Errorf("expected player_1 after placement, got %s", state.CurrentPlayer)
	}
}

func TestSetupPlaceTroopRoundRobin(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Human", 1, []string{"Bot"})

	// Place for player_0
	owned0 := findOwnedTerritory(state, "player_0")
	err := e.SetupPlaceTroop(state, owned0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should now be player_1's turn
	if state.CurrentPlayer != "player_1" {
		t.Fatalf("expected player_1, got %s", state.CurrentPlayer)
	}

	// Place for player_1
	owned1 := findOwnedTerritory(state, "player_1")
	err = e.SetupPlaceTroop(state, owned1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be back to player_0
	if state.CurrentPlayer != "player_0" {
		t.Errorf("expected player_0 after round-robin, got %s", state.CurrentPlayer)
	}
}

func TestSetupPlaceTroopWrongPhase(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	err := e.SetupPlaceTroop(state, "alaska")
	if err == nil {
		t.Error("expected error for setup placement during attack phase")
	}
}

func TestIsSetupComplete(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Human", 1, []string{"Bot"})

	// Right after creation, setup should not be complete (each territory has 1 troop, need 40 each for 2 players)
	if e.IsSetupComplete(state) {
		t.Error("expected setup to not be complete at start")
	}

	// Force all players to have enough troops placed to complete setup
	// For 2 players: 40 starting troops each, territories start with 1 troop
	// Need to give each player enough troops on their territories
	for id, terr := range state.Territories {
		terr.Troops = 4 // boost troops
		state.Territories[id] = terr
	}

	// Now total placed per player should be around 4 * 21 = 84 which exceeds 40
	if !e.IsSetupComplete(state) {
		t.Error("expected setup to be complete when all troops placed")
	}
}

func TestIsCurrentPlayerAI(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Human", 1, []string{"Bot"})

	// player_0 is human
	state.CurrentPlayer = "player_0"
	if e.IsCurrentPlayerAI(state) {
		t.Error("expected player_0 to be human")
	}

	// player_1 is AI
	state.CurrentPlayer = "player_1"
	if !e.IsCurrentPlayerAI(state) {
		t.Error("expected player_1 to be AI")
	}
}

func TestGetSetupTroopsRemaining(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Human", 1, []string{"Bot"})

	// For 2 players: 40 starting troops each
	// Each player has ~21 territories with 1 troop each = 21 placed
	remaining := e.GetSetupTroopsRemaining(state, "player_0")
	if remaining <= 0 {
		t.Errorf("expected remaining troops > 0, got %d", remaining)
	}

	// Count territories owned by player_0
	count := 0
	for _, terr := range state.Territories {
		if terr.Owner == "player_0" {
			count++
		}
	}
	expected := 40 - count // each territory has 1 troop
	if remaining != expected {
		t.Errorf("expected remaining %d, got %d", expected, remaining)
	}
}
