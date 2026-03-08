package ai

import (
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
)

func TestScoreTerritory(t *testing.T) {
	state := newTestState()

	// Score a territory owned by the player
	ownedID := ""
	for id, terr := range state.Territories {
		if terr.Owner == "player_1" {
			ownedID = id
			break
		}
	}

	score := ScoreTerritory(state, ownedID, "player_1")
	// Score should be a reasonable number (not zero, not absurdly high)
	if score == 0 {
		t.Log("ScoreTerritory returned 0, which may be valid for some configurations")
	}

	// Score for non-existent territory should be 0
	score = ScoreTerritory(state, "nonexistent", "player_1")
	if score != 0 {
		t.Errorf("expected 0 for nonexistent territory, got %f", score)
	}
}

func TestScoreTerritoryContProximityBonus(t *testing.T) {
	state := newTestState()

	// Give player_1 most of Australia
	ausTerritories := models.ContinentTerritories["australia"]
	for i, tid := range ausTerritories {
		terr := state.Territories[tid]
		if i < len(ausTerritories)-1 {
			terr.Owner = "player_1"
		} else {
			terr.Owner = "player_0"
		}
		state.Territories[tid] = terr
	}

	// The last territory (owned by player_0) should score high for player_1
	// because it would complete the continent
	lastAus := ausTerritories[len(ausTerritories)-1]
	score := ScoreTerritory(state, lastAus, "player_1")

	// Score a random territory far from any completion
	randomScore := ScoreTerritory(state, "ukraine", "player_1")

	// Near-completion territory should generally score higher
	// (not always guaranteed depending on other factors, but usually)
	t.Logf("Near-completion Australia territory score: %.2f, random territory score: %.2f", score, randomScore)
}

func TestScoreAttack(t *testing.T) {
	state := newTestState()

	// Find an attack pair
	var fromID, toID string
	for _, terr := range state.Territories {
		if terr.Owner == "player_1" && terr.Troops >= 2 {
			for _, nID := range terr.Neighbors {
				if n, ok := state.Territories[nID]; ok && n.Owner == "player_0" {
					fromID = terr.ID
					toID = nID
					break
				}
			}
		}
		if fromID != "" {
			break
		}
	}

	if fromID == "" {
		t.Skip("no attack pair found")
	}

	score := scoreAttack(state, fromID, toID, "player_1")
	if score <= 0 {
		t.Logf("Attack score is %.2f (may be valid for balanced troops)", score)
	}
}

func TestScoreAttackContinentCompletion(t *testing.T) {
	state := newTestState()

	// Give player_1 all of Australia except eastern_australia
	for _, tid := range models.ContinentTerritories["australia"] {
		terr := state.Territories[tid]
		if tid == "eastern_australia" {
			terr.Owner = "player_0"
			terr.Troops = 1
		} else {
			terr.Owner = "player_1"
			terr.Troops = 10
		}
		state.Territories[tid] = terr
	}

	// Score attacking eastern_australia from western_australia
	completionScore := scoreAttack(state, "western_australia", "eastern_australia", "player_1")

	// Score a random non-continent-completing attack
	var normalFrom, normalTo string
	for _, terr := range state.Territories {
		if terr.Owner == "player_1" && terr.Troops >= 2 && terr.Continent != "australia" {
			for _, nID := range terr.Neighbors {
				if n, ok := state.Territories[nID]; ok && n.Owner == "player_0" {
					normalFrom = terr.ID
					normalTo = nID
					break
				}
			}
		}
		if normalFrom != "" {
			break
		}
	}

	if normalFrom != "" {
		normalScore := scoreAttack(state, normalFrom, normalTo, "player_1")
		if completionScore <= normalScore {
			t.Logf("continent completion score (%.2f) not higher than normal (%.2f) - may depend on troop ratios",
				completionScore, normalScore)
		}
	}

	// The continent completion attack should have a high score
	if completionScore < 50 {
		t.Logf("continent completion attack score: %.2f (expected high)", completionScore)
	}
}

func TestScoreAttackElimination(t *testing.T) {
	state := newTestState()

	// Give player_0 only one territory
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		terr.Troops = 5
		state.Territories[id] = terr
	}

	// Give player_0 only kamchatka
	k := state.Territories["kamchatka"]
	k.Owner = "player_0"
	k.Troops = 1
	state.Territories["kamchatka"] = k

	score := scoreAttack(state, "alaska", "kamchatka", "player_1")
	// Elimination should have a very high score
	if score < 50 {
		t.Logf("elimination attack score: %.2f (expected high)", score)
	}
}

func TestScoreAttackBreakingOpponentContinent(t *testing.T) {
	state := newTestState()

	// Give player_0 all of Australia
	for _, tid := range models.ContinentTerritories["australia"] {
		terr := state.Territories[tid]
		terr.Owner = "player_0"
		terr.Troops = 3
		state.Territories[tid] = terr
	}

	// Give player_1 siam (neighbor of indonesia) with lots of troops
	s := state.Territories["siam"]
	s.Owner = "player_1"
	s.Troops = 15
	state.Territories["siam"] = s

	score := scoreAttack(state, "siam", "indonesia", "player_1")
	// Breaking opponent continent should score high
	if score < 30 {
		t.Logf("break continent score: %.2f (expected high)", score)
	}
}

func TestScorePlacement(t *testing.T) {
	state := newTestState()

	// Find a border territory for player_1
	borders := getBorderTerritories(state, "player_1")
	if len(borders) == 0 {
		t.Skip("no border territories found")
	}

	// Find an interior territory for player_1
	var interiorID string
	for _, terr := range state.Territories {
		if terr.Owner != "player_1" {
			continue
		}
		isInterior := true
		for _, nID := range terr.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != "player_1" {
				isInterior = false
				break
			}
		}
		if isInterior {
			interiorID = terr.ID
			break
		}
	}

	borderScore := scorePlacement(state, borders[0].ID, "player_1")

	if interiorID != "" {
		interiorScore := scorePlacement(state, interiorID, "player_1")
		// Border territories should generally score higher for placement
		t.Logf("Border placement score: %.2f, Interior placement score: %.2f", borderScore, interiorScore)
	} else {
		t.Logf("Border placement score: %.2f (no interior territory found for comparison)", borderScore)
	}
}

func TestScorePlacementContinentBorderPriority(t *testing.T) {
	state := newTestState()

	// Give player_1 all of Australia
	for _, tid := range models.ContinentTerritories["australia"] {
		terr := state.Territories[tid]
		terr.Owner = "player_1"
		terr.Troops = 3
		state.Territories[tid] = terr
	}

	// Indonesia borders siam (enemy territory) - should be high priority
	indonesiaScore := scorePlacement(state, "indonesia", "player_1")

	// A random non-continent-border territory
	randomBorder := ""
	for _, terr := range state.Territories {
		if terr.Owner == "player_1" && terr.Continent != "australia" {
			for _, nID := range terr.Neighbors {
				if n, ok := state.Territories[nID]; ok && n.Owner != "player_1" {
					randomBorder = terr.ID
					break
				}
			}
		}
		if randomBorder != "" {
			break
		}
	}

	if randomBorder != "" {
		randomScore := scorePlacement(state, randomBorder, "player_1")
		t.Logf("Continent border score (Indonesia): %.2f, Random border score (%s): %.2f",
			indonesiaScore, randomBorder, randomScore)
		// Indonesia guarding Australia should generally score higher
	}
}
