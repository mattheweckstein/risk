package ai

import (
	"math"

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

	// --- Continent completion scoring ---
	continentTerritories := models.ContinentTerritories[continent]
	totalInContinent := len(continentTerritories)
	ownedInContinent := 0
	for _, cTerr := range continentTerritories {
		if ct, exists := state.Territories[cTerr]; exists && ct.Owner == playerID {
			ownedInContinent++
		}
	}

	ratio := 0.0
	if totalInContinent > 0 {
		ratio = float64(ownedInContinent) / float64(totalInContinent)
	}

	bonus := models.ContinentBonuses[continent]
	// Efficiency = bonus per territory (how valuable is this continent to hold)
	efficiency := float64(bonus) / float64(maxInt(1, totalInContinent))

	// Base continent affinity: the more we own, the more valuable remaining ones are
	// This creates exponential urgency as we get close to completing
	score += float64(ownedInContinent) * 8

	// Strong bonus for near-completion continents (scales with bonus value)
	if ratio >= 0.75 {
		// Almost complete -- very high priority
		score += 40 * efficiency * ratio
	} else if ratio >= 0.5 {
		score += 15 * efficiency * ratio
	}

	// Small continents (Australia, South America) get an extra bump because
	// they're easier to hold with fewer border territories
	if totalInContinent <= 4 {
		score += 8 * ratio
	}

	// --- Border defense of owned continents ---
	// If we own a continent, its border territories are critical to defend
	for contName := range models.ContinentTerritories {
		if contName == continent {
			continue
		}
		if ownsContinent(state, contName, playerID) {
			// Check if this territory borders that completed continent
			if territoryBordersContinent(state, territoryID, contName, playerID) {
				score += 15 + float64(models.ContinentBonuses[contName])*2
			}
		}
	}

	// --- Enemy continent breaking ---
	// Territories inside a continent fully owned by an opponent are high-value targets
	if t.Owner != playerID {
		for _, p := range state.Players {
			if p.ID == playerID || !p.IsAlive {
				continue
			}
			if ownsContinent(state, continent, p.ID) {
				// Breaking an opponent's continent is very valuable
				score += 25 + float64(bonus)*5
			}
		}
	}

	// --- Neighbor threat/value ---
	enemyNeighborTroops := 0
	enemyNeighbors := 0
	friendlyNeighbors := 0
	for _, neighborID := range t.Neighbors {
		if n, exists := state.Territories[neighborID]; exists {
			if n.Owner == playerID {
				friendlyNeighbors++
			} else if n.Owner != "" {
				enemyNeighbors++
				enemyNeighborTroops += n.Troops
			}
		}
	}

	// Territories with fewer enemy neighbors are easier to hold (chokepoints)
	if enemyNeighbors > 0 {
		score -= float64(enemyNeighbors) * 1.5
	}
	// Interior territories are less valuable for placement (not on the front)
	if enemyNeighbors == 0 && t.Owner == playerID {
		score -= 10
	}
	// Bonus for territories with many friendly neighbors (connected position)
	score += float64(friendlyNeighbors) * 1.5

	return score
}

// scoreAttack evaluates how good a specific attack is considering strategic context.
// Returns a score where higher = more desirable attack.
func scoreAttack(state *models.GameState, fromID, toID, playerID string) float64 {
	from := state.Territories[fromID]
	to := state.Territories[toID]

	attackerTroops := from.Troops - 1
	defenderTroops := maxInt(1, to.Troops)

	// Base score: troop ratio (higher is better odds)
	ratio := float64(attackerTroops) / float64(defenderTroops)
	score := ratio * 10

	// --- Continent completion ---
	continent := to.Continent
	completionRatio := continentCompletionRatio(state, continent, playerID)
	continentSize := len(models.ContinentTerritories[continent])
	bonus := models.ContinentBonuses[continent]

	// How many territories do we still need in this continent?
	needed := 0
	for _, tID := range models.ContinentTerritories[continent] {
		if ct, ok := state.Territories[tID]; ok && ct.Owner != playerID {
			needed++
		}
	}

	if needed == 1 {
		// This attack would complete the continent!
		score += 80 + float64(bonus)*10
	} else if needed == 2 && completionRatio >= 0.5 {
		score += 40 + float64(bonus)*5
	} else if completionRatio >= 0.5 {
		score += 15 * completionRatio * float64(bonus) / float64(maxInt(1, continentSize))
	}

	// --- Breaking opponent continents ---
	defenderID := to.Owner
	if defenderID != "" && ownsContinent(state, continent, defenderID) {
		// Breaking an opponent's continent bonus is very valuable
		score += 50 + float64(bonus)*8
	}

	// --- Player elimination ---
	if defenderID != "" {
		defenderTerritoryCount := countPlayerTerritories(state, defenderID)
		defenderCards := countPlayerCards(state, defenderID)

		if defenderTerritoryCount == 1 {
			// Eliminating this player!
			score += 60
			// Extra value if they have cards we'd steal
			score += float64(defenderCards) * 15
		} else if defenderTerritoryCount == 2 {
			// Getting close to elimination
			score += 25
			score += float64(defenderCards) * 8
		} else if defenderTerritoryCount <= 4 && defenderCards >= 4 {
			// Weak player with cards worth targeting
			score += 20 + float64(defenderCards)*5
		}
	}

	// --- Favor attacking weak defenders ---
	// A territory with 1 troop is much easier than one with 3
	if defenderTroops == 1 {
		score += 10
	}

	// Small penalty for attacks that would leave the source very weak
	remainingAfterEstimate := float64(attackerTroops) - float64(defenderTroops)*0.7
	if remainingAfterEstimate < 2 {
		score -= 5
	}

	return score
}

// scorePlacement evaluates how valuable it is to place troops on a specific
// border territory. Considers continent defense, attack staging, and threat level.
func scorePlacement(state *models.GameState, territoryID, playerID string) float64 {
	t := state.Territories[territoryID]
	score := ScoreTerritory(state, territoryID, playerID)

	continent := t.Continent

	// --- Continent border defense ---
	// If we own this continent, borders need heavy reinforcement
	if ownsContinent(state, continent, playerID) {
		maxEnemyThreat := 0
		for _, nID := range t.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
				if n.Troops > maxEnemyThreat {
					maxEnemyThreat = n.Troops
				}
			}
		}
		if maxEnemyThreat > 0 {
			// We need to be strong enough to hold
			deficit := float64(maxEnemyThreat) - float64(t.Troops)
			if deficit > 0 {
				score += deficit * 3
			}
			score += 20 + float64(models.ContinentBonuses[continent])*3
		}
	}

	// --- Attack staging ---
	// Check if this territory borders a continent we're close to completing
	for contName, contTerrs := range models.ContinentTerritories {
		ratio := continentCompletionRatio(state, contName, playerID)
		if ratio < 0.5 {
			continue
		}
		// Check if this territory can attack into that continent
		for _, nID := range t.Neighbors {
			for _, ct := range contTerrs {
				if nID == ct {
					n := state.Territories[nID]
					if n.Owner != playerID {
						score += 20 * ratio
						goto doneStaging
					}
				}
			}
		}
	}
doneStaging:

	// --- Threat response ---
	// Territories facing large enemy forces need reinforcement
	totalEnemyThreat := 0.0
	bordersLeader := false
	leaderID := getLeaderID(state)
	for _, nID := range t.Neighbors {
		if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
			totalEnemyThreat += float64(n.Troops)
			if n.Owner == leaderID && leaderID != playerID {
				bordersLeader = true
			}
		}
	}
	if totalEnemyThreat > 0 {
		// Reinforce if enemy troops nearby significantly exceed ours
		threatRatio := totalEnemyThreat / math.Max(1, float64(t.Troops))
		if threatRatio > 1.5 {
			score += threatRatio * 5
		}
	}

	// --- Leader targeting ---
	// Prioritize placing troops on borders with the leading player
	if bordersLeader {
		score += 15
	}

	return score
}

// territoryBordersContinent returns true if the territory (owned by playerID,
// outside the continent) neighbors any territory in the given continent.
func territoryBordersContinent(state *models.GameState, territoryID, continent, playerID string) bool {
	t := state.Territories[territoryID]
	contTerrs := models.ContinentTerritories[continent]
	for _, nID := range t.Neighbors {
		for _, ct := range contTerrs {
			if nID == ct {
				n := state.Territories[nID]
				if n.Owner != playerID {
					return true
				}
			}
		}
	}
	return false
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

// countPlayerTerritories returns how many territories a player owns.
func countPlayerTerritories(state *models.GameState, playerID string) int {
	count := 0
	for _, t := range state.Territories {
		if t.Owner == playerID {
			count++
		}
	}
	return count
}

// countPlayerCards returns how many cards a player holds.
func countPlayerCards(state *models.GameState, playerID string) int {
	for i := range state.Players {
		if state.Players[i].ID == playerID {
			return len(state.Players[i].Cards)
		}
	}
	return 0
}

// getLeaderID returns the player who owns the most territories.
// In case of tie, prefers the human player (they're the real threat).
func getLeaderID(state *models.GameState) string {
	bestID := ""
	bestCount := 0
	isHuman := false
	for _, p := range state.Players {
		if !p.IsAlive {
			continue
		}
		count := countPlayerTerritories(state, p.ID)
		if count > bestCount || (count == bestCount && !p.IsAI && !isHuman) {
			bestCount = count
			bestID = p.ID
			isHuman = !p.IsAI
		}
	}
	return bestID
}

// getPlayerStrength returns a rough measure of a player's power:
// territory count + continent bonuses + cards held.
func getPlayerStrength(state *models.GameState, playerID string) float64 {
	terrs := countPlayerTerritories(state, playerID)
	strength := float64(terrs)

	// Add continent bonus value
	for continent, bonus := range models.ContinentBonuses {
		if ownsContinent(state, continent, playerID) {
			strength += float64(bonus) * 3
		}
	}

	// Cards are potential future troops
	cards := countPlayerCards(state, playerID)
	strength += float64(cards) * 2

	// Total troop count matters
	totalTroops := 0
	for _, t := range state.Territories {
		if t.Owner == playerID {
			totalTroops += t.Troops
		}
	}
	strength += float64(totalTroops) * 0.3

	return strength
}

// isLeader returns true if the given player is the strongest player in the game.
func isLeader(state *models.GameState, playerID string) bool {
	return getLeaderID(state) == playerID
}

// leaderThreatBonus returns extra attack score when targeting the game leader.
func leaderThreatBonus(state *models.GameState, defenderID, attackerID string) float64 {
	if defenderID == "" {
		return 0
	}

	leaderID := getLeaderID(state)
	if defenderID != leaderID {
		return 0
	}

	// The bigger the lead, the more important to attack them
	leaderStrength := getPlayerStrength(state, leaderID)
	myStrength := getPlayerStrength(state, attackerID)

	if leaderStrength <= myStrength {
		return 0
	}

	gap := leaderStrength - myStrength
	// Scale bonus with the gap: bigger lead = more urgency
	bonus := 15 + gap*0.5
	if bonus > 50 {
		bonus = 50
	}
	return bonus
}
