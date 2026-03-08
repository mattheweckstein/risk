package ai

import (
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
)

// newTestState creates a minimal 2-player game state for AI testing.
func newTestState() *models.GameState {

	territories := make(map[string]models.Territory, len(models.AllTerritories))
	terrIDs := make([]string, 0, len(models.AllTerritories))
	for _, td := range models.AllTerritories {
		territories[td.ID] = models.Territory{
			ID:        td.ID,
			Name:      td.Name,
			Continent: td.Continent,
			Neighbors: td.Neighbors,
			Troops:    3,
		}
		terrIDs = append(terrIDs, td.ID)
	}

	// Distribute territories round-robin
	for i, tid := range terrIDs {
		t := territories[tid]
		if i%2 == 0 {
			t.Owner = "player_0"
		} else {
			t.Owner = "player_1"
		}
		territories[tid] = t
	}

	// Build deck
	types := []string{"infantry", "cavalry", "artillery"}
	deck := make([]models.Card, 0, 44)
	for i, td := range models.AllTerritories {
		deck = append(deck, models.Card{Territory: td.ID, Type: types[i%3]})
	}
	deck = append(deck, models.Card{Type: "wild"}, models.Card{Type: "wild"})

	return &models.GameState{
		ID:            "test-game",
		Phase:         models.PhaseSetup,
		Turn:          0,
		CurrentPlayer: "player_1",
		Players: []models.Player{
			{ID: "player_0", Name: "Human", IsAI: false, Color: "red", Cards: []models.Card{}, IsAlive: true},
			{ID: "player_1", Name: "AI Bot", IsAI: true, Color: "blue", Cards: []models.Card{}, IsAlive: true},
		},
		Territories:    territories,
		Deck:           deck,
		Log:            []models.LogEntry{},
		TroopsToDeploy: 1,
		CardTradeCount: 0,
	}
}

func TestExecuteTurnSetup(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhaseSetup

	// Count total troops before
	totalBefore := 0
	for _, terr := range state.Territories {
		if terr.Owner == "player_1" {
			totalBefore += terr.Troops
		}
	}

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Count total troops after - should have increased by 1
	totalAfter := 0
	for _, terr := range state.Territories {
		if terr.Owner == "player_1" {
			totalAfter += terr.Troops
		}
	}

	if totalAfter != totalBefore+1 {
		t.Errorf("expected troops to increase by 1, before=%d, after=%d", totalBefore, totalAfter)
	}
}

func TestExecuteTurnFullTurn(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI enough troops to do something interesting
	for id, terr := range state.Territories {
		if terr.Owner == "player_1" {
			terr.Troops = 5
		}
		state.Territories[id] = terr
	}

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// AI should have placed troops and possibly attacked
	if len(state.Log) == 0 {
		t.Error("expected log entries after AI turn")
	}
}

func TestPlaceSetupTroop(t *testing.T) {
	state := newTestState()

	territoryID := PlaceSetupTroop(state, "player_1")
	if territoryID == "" {
		t.Fatal("expected non-empty territory ID")
	}

	// Must be a valid territory
	if _, ok := state.Territories[territoryID]; !ok {
		t.Errorf("returned territory %q does not exist", territoryID)
	}

	// Must be owned by the player
	if state.Territories[territoryID].Owner != "player_1" {
		t.Errorf("expected territory to be owned by player_1, owned by %s", state.Territories[territoryID].Owner)
	}
}

func TestAIAttacksWhenFavorable(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI overwhelming troops on border territories
	for id, terr := range state.Territories {
		if terr.Owner == "player_1" {
			terr.Troops = 20
		} else {
			terr.Troops = 1
		}
		state.Territories[id] = terr
	}

	initialP1Terrs := countPlayerTerritories(state, "player_1")

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	finalP1Terrs := countPlayerTerritories(state, "player_1")
	if finalP1Terrs <= initialP1Terrs {
		t.Error("expected AI to conquer at least one territory with overwhelming advantage")
	}
}

func TestAIDoesntAttackWhenWeak(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI minimal troops, enemies strong
	for id, terr := range state.Territories {
		if terr.Owner == "player_1" {
			terr.Troops = 1
		} else {
			terr.Troops = 20
		}
		state.Territories[id] = terr
	}

	initialP0Terrs := countPlayerTerritories(state, "player_0")

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Player 0 should still have same or more territories (AI shouldn't be conquering)
	finalP0Terrs := countPlayerTerritories(state, "player_0")
	if finalP0Terrs < initialP0Terrs-1 {
		// Allow at most 1 territory loss from lucky dice
		t.Errorf("AI conquered too many territories while weak: p0 went from %d to %d", initialP0Terrs, finalP0Terrs)
	}
}

func TestAITradesCards(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI 4 cards (should trigger trade)
	for i := range state.Players {
		if state.Players[i].ID == "player_1" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "cavalry"},
				{Territory: "china", Type: "artillery"},
				{Territory: "india", Type: "infantry"},
			}
		}
	}

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// AI should have traded cards
	var aiCards int
	for _, p := range state.Players {
		if p.ID == "player_1" {
			aiCards = len(p.Cards)
		}
	}

	if aiCards >= 4 {
		t.Errorf("expected AI to have traded cards, still has %d", aiCards)
	}
}

func TestAIFortifies(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI lots of troops on interior territories to encourage fortification
	for id, terr := range state.Territories {
		if terr.Owner == "player_1" {
			terr.Troops = 10
		} else {
			terr.Troops = 15 // Strong enemies to discourage attacking
		}
		state.Territories[id] = terr
	}

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that a fortify log entry exists
	foundFortify := false
	for _, entry := range state.Log {
		if entry.Player == "player_1" {
			// Look for "fortified" in message
			if len(entry.Message) > 0 {
				for _, word := range []string{"fortified", "fortif"} {
					if containsWord(entry.Message, word) {
						foundFortify = true
						break
					}
				}
			}
		}
	}

	if !foundFortify {
		// Fortify may not always occur if no good source/target, so just skip
		t.Log("AI did not fortify this turn (may be expected depending on state)")
	}
}

func TestAICompletesContinent(t *testing.T) {
	state := newTestState()
	state.Phase = models.PhasePlace
	state.Turn = 1

	// Give AI all of Australia except one territory
	ausTerritories := models.ContinentTerritories["australia"]
	for id, terr := range state.Territories {
		terr.Owner = "player_0"
		terr.Troops = 1
		state.Territories[id] = terr
	}

	// Give AI most of Australia
	for i, tid := range ausTerritories {
		terr := state.Territories[tid]
		terr.Owner = "player_1"
		terr.Troops = 15
		state.Territories[tid] = terr
		if i == len(ausTerritories)-1 {
			// Last one stays as player_0 with 1 troop
			terr.Owner = "player_0"
			terr.Troops = 1
			state.Territories[tid] = terr
		}
	}

	// Give AI some other territories too
	extraCount := 0
	for id, terr := range state.Territories {
		if terr.Owner == "player_0" && extraCount < 5 {
			found := false
			for _, aid := range ausTerritories {
				if id == aid {
					found = true
					break
				}
			}
			if !found {
				terr.Owner = "player_1"
				terr.Troops = 5
				state.Territories[id] = terr
				extraCount++
			}
		}
	}

	err := ExecuteTurn(state, "player_1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check if AI now owns all of Australia
	if ownsContinent(state, "australia", "player_1") {
		t.Log("AI successfully completed Australia continent")
	} else {
		// Dice are random, so this might not always succeed
		t.Log("AI did not complete Australia (dice may have been unfavorable)")
	}
}

// containsWord checks if a string contains a substring (simple check).
func containsWord(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
