package ai

import (
	"math/rand"
	"sort"

	"github.com/mattheweckstein/risk/backend/models"
)

// getOwnedTerritories returns all territories owned by the given player.
func getOwnedTerritories(state *models.GameState, playerID string) []models.Territory {
	var owned []models.Territory
	for _, t := range state.Territories {
		if t.Owner == playerID {
			owned = append(owned, t)
		}
	}
	return owned
}

// getBorderTerritories returns territories owned by the player that have at
// least one neighbor owned by a different player.
func getBorderTerritories(state *models.GameState, playerID string) []models.Territory {
	var borders []models.Territory
	for _, t := range state.Territories {
		if t.Owner != playerID {
			continue
		}
		for _, nID := range t.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != playerID {
				borders = append(borders, t)
				break
			}
		}
	}
	return borders
}

// isConnected checks whether two territories are connected through territories
// owned by the given player, using BFS.
func isConnected(state *models.GameState, from, to, playerID string) bool {
	if from == to {
		return true
	}

	visited := map[string]bool{from: true}
	queue := []string{from}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		ct, ok := state.Territories[current]
		if !ok {
			continue
		}
		for _, nID := range ct.Neighbors {
			if nID == to {
				// The destination just needs to be owned by the player too
				if n, exists := state.Territories[nID]; exists && n.Owner == playerID {
					return true
				}
				continue
			}
			if visited[nID] {
				continue
			}
			if n, exists := state.Territories[nID]; exists && n.Owner == playerID {
				visited[nID] = true
				queue = append(queue, nID)
			}
		}
	}
	return false
}

// findValidCardSet searches the player's cards for a valid set of 3 to trade.
// Returns the indices (into the cards slice) and whether a set was found.
func findValidCardSet(cards []models.Card) (indices [3]int, found bool) {
	n := len(cards)
	if n < 3 {
		return indices, false
	}

	// Try all combinations of 3 cards
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				if isValidSet(cards[i], cards[j], cards[k]) {
					return [3]int{i, j, k}, true
				}
			}
		}
	}
	return indices, false
}

// isValidSet returns true if the three cards form a valid trade-in set.
func isValidSet(a, b, c models.Card) bool {
	types := []string{a.Type, b.Type, c.Type}

	// Count wilds
	wilds := 0
	for _, t := range types {
		if t == "wild" {
			wilds++
		}
	}

	// Any set with a wild card is valid
	if wilds > 0 {
		return true
	}

	// Three of the same type
	if types[0] == types[1] && types[1] == types[2] {
		return true
	}

	// One of each type
	sort.Strings(types)
	if types[0] == "artillery" && types[1] == "cavalry" && types[2] == "infantry" {
		return true
	}

	return false
}

// rollDice rolls n dice and returns the results sorted descending.
func rollDice(n int) []int {
	results := make([]int, n)
	for i := range results {
		results[i] = rand.Intn(6) + 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(results)))
	return results
}

// getPlayerName returns the display name for a player.
func getPlayerName(state *models.GameState, playerID string) string {
	for _, p := range state.Players {
		if p.ID == playerID {
			return p.Name
		}
	}
	return playerID
}

// getPlayer returns a pointer to the player in the state, or nil.
func getPlayer(state *models.GameState, playerID string) *models.Player {
	for i := range state.Players {
		if state.Players[i].ID == playerID {
			return &state.Players[i]
		}
	}
	return nil
}

// addLog appends a log entry to the game state.
func addLog(state *models.GameState, playerID, message string) {
	state.Log = append(state.Log, models.LogEntry{
		Turn:    state.Turn,
		Player:  playerID,
		Message: message,
	})
}

// min returns the smaller of two ints.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxInt returns the larger of two ints.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
