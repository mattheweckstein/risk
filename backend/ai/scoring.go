package ai

import (
	"github.com/mattheweckstein/risk/backend/models"
)

// ScoreTerritory evaluates how valuable a territory is for the given player.
// Higher scores indicate more strategically important territories.
func ScoreTerritory(state *models.GameState, territoryID, playerID string) float64 {
	t, ok := state.Territories[territoryID]
	if !ok {
		return 0
	}

	score := 0.0
	continent := t.Continent

	// Count how many territories in the same continent this player owns
	continentTerritories := models.ContinentTerritories[continent]
	ownedInContinent := 0
	for _, cTerr := range continentTerritories {
		if ct, exists := state.Territories[cTerr]; exists && ct.Owner == playerID {
			ownedInContinent++
		}
	}

	// +10 for each territory in the same continent already owned
	score += float64(ownedInContinent) * 10

	// +5 * continent bonus if player owns > 50% of the continent
	totalInContinent := len(continentTerritories)
	if totalInContinent > 0 && float64(ownedInContinent)/float64(totalInContinent) > 0.5 {
		bonus := models.ContinentBonuses[continent]
		score += 5 * float64(bonus)
	}

	// +3 for being a border of a completed continent
	for contName, contTerrs := range models.ContinentTerritories {
		if contName == continent {
			continue
		}
		if ownsContinent(state, contName, playerID) {
			// Check if this territory neighbors any territory in that completed continent
			for _, neighbor := range t.Neighbors {
				for _, ct := range contTerrs {
					if neighbor == ct {
						score += 3
						goto nextContinent
					}
				}
			}
		}
	nextContinent:
	}

	// -2 for each neighbor owned by an enemy
	for _, neighborID := range t.Neighbors {
		if n, exists := state.Territories[neighborID]; exists && n.Owner != playerID && n.Owner != "" {
			score -= 2
		}
	}

	return score
}

// ownsContinent returns true if the player owns every territory in the continent.
func ownsContinent(state *models.GameState, continent, playerID string) bool {
	for _, tID := range models.ContinentTerritories[continent] {
		if t, ok := state.Territories[tID]; !ok || t.Owner != playerID {
			return false
		}
	}
	return true
}

// continentCompletionRatio returns the fraction of a continent owned by the player.
func continentCompletionRatio(state *models.GameState, continent, playerID string) float64 {
	terrs := models.ContinentTerritories[continent]
	if len(terrs) == 0 {
		return 0
	}
	owned := 0
	for _, tID := range terrs {
		if t, ok := state.Territories[tID]; ok && t.Owner == playerID {
			owned++
		}
	}
	return float64(owned) / float64(len(terrs))
}
