package ai

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/mattheweckstein/risk/backend/models"
)

// ExecuteTurn executes a full AI turn (placement, attacks, fortify), mutating
// the game state directly. During setup phase it places a single troop.
func ExecuteTurn(state *models.GameState, playerID string) error {
	if state.Phase == models.PhaseSetup {
		territoryID := PlaceSetupTroop(state, playerID)
		t := state.Territories[territoryID]
		t.Troops++
		state.Territories[territoryID] = t
		name := getPlayerName(state, playerID)
		addLog(state, playerID, fmt.Sprintf("%s placed 1 troop on %s", name, t.Name))
		return nil
	}

	// Full turn: place -> attack -> fortify
	if err := executePlacement(state, playerID); err != nil {
		return fmt.Errorf("placement phase: %w", err)
	}

	executeAttacks(state, playerID)

	executeFortify(state, playerID)

	return nil
}

// PlaceSetupTroop returns the territory ID where the AI wants to place a troop
// during setup phase. It picks an owned territory with the highest score, or
// an unowned territory if available.
func PlaceSetupTroop(state *models.GameState, playerID string) string {
	// First, check for unowned territories (initial claiming phase)
	var unowned []string
	for id, t := range state.Territories {
		if t.Owner == "" {
			unowned = append(unowned, id)
		}
	}

	if len(unowned) > 0 {
		// Score unowned territories and pick the best
		type scored struct {
			id    string
			score float64
		}
		var candidates []scored
		for _, id := range unowned {
			s := ScoreTerritory(state, id, playerID)
			candidates = append(candidates, scored{id, s})
		}
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score > candidates[j].score
		})
		// Add some randomness: pick from top 3
		top := minInt(3, len(candidates))
		return candidates[rand.Intn(top)].id
	}

	// All territories claimed — reinforce an owned territory
	owned := getOwnedTerritories(state, playerID)
	if len(owned) == 0 {
		// Shouldn't happen, but be safe
		for id := range state.Territories {
			return id
		}
	}

	return pickBestTerritory(state, owned, playerID)
}

// pickBestTerritory returns the ID of the highest-scored territory from the
// given list. Adds a little randomness.
func pickBestTerritory(state *models.GameState, territories []models.Territory, playerID string) string {
	type scored struct {
		id    string
		score float64
	}
	var candidates []scored
	for _, t := range territories {
		s := ScoreTerritory(state, t.ID, playerID)
		candidates = append(candidates, scored{t.ID, s})
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// Pick from the top 3 with weighted randomness
	top := minInt(3, len(candidates))
	weights := []float64{0.6, 0.25, 0.15}
	r := rand.Float64()
	cumulative := 0.0
	for i := 0; i < top; i++ {
		cumulative += weights[i]
		if r < cumulative {
			return candidates[i].id
		}
	}
	return candidates[0].id
}

// executePlacement handles the troop deployment phase.
func executePlacement(state *models.GameState, playerID string) error {
	player := getPlayer(state, playerID)
	if player == nil {
		return fmt.Errorf("player %s not found", playerID)
	}
	name := player.Name

	// Trade cards if holding 4+ (or forced at 5+)
	tradeCards(state, player)

	// Calculate troops to deploy
	owned := getOwnedTerritories(state, playerID)
	troops := maxInt(3, len(owned)/3)

	// Add continent bonuses
	for continent, bonus := range models.ContinentBonuses {
		if ownsContinent(state, continent, playerID) {
			troops += bonus
		}
	}

	addLog(state, playerID, fmt.Sprintf("%s receives %d troops to deploy", name, troops))

	// Distribute troops among highest-scored territories
	borders := getBorderTerritories(state, playerID)
	if len(borders) == 0 {
		borders = owned
	}

	// Score and sort border territories
	type scored struct {
		id    string
		score float64
	}
	var candidates []scored
	for _, t := range borders {
		s := ScoreTerritory(state, t.ID, playerID)
		// Boost score for territories in continents we're close to completing
		continent := t.Continent
		ratio := continentCompletionRatio(state, continent, playerID)
		if ratio >= 0.6 {
			s += 20 * ratio
		}
		candidates = append(candidates, scored{t.ID, s})
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// Place troops on top candidates
	remaining := troops
	for remaining > 0 && len(candidates) > 0 {
		// Distribute: put more on the top territory, some on others
		topCount := minInt(len(candidates), 3)
		for i := 0; i < topCount && remaining > 0; i++ {
			place := remaining
			if i == 0 {
				// Place roughly half on the top territory
				place = maxInt(1, remaining/2)
			} else {
				place = maxInt(1, remaining/3)
			}
			if place > remaining {
				place = remaining
			}

			tID := candidates[i].id
			t := state.Territories[tID]
			t.Troops += place
			state.Territories[tID] = t
			remaining -= place

			addLog(state, playerID, fmt.Sprintf("%s placed %d troops on %s", name, place, t.Name))
		}
	}

	return nil
}

// tradeCards handles card trading for the AI.
func tradeCards(state *models.GameState, player *models.Player) {
	for len(player.Cards) >= 4 {
		indices, found := findValidCardSet(player.Cards)
		if !found {
			if len(player.Cards) >= 5 {
				// Forced to trade but no valid set found — shouldn't happen with 5 cards
				break
			}
			break
		}

		// Calculate trade-in bonus
		state.CardTradeCount++
		bonus := cardTradeBonus(state.CardTradeCount)

		name := player.Name
		addLog(state, player.ID, fmt.Sprintf("%s traded in cards for %d troops", name, bonus))

		// Remove traded cards (in reverse index order to preserve indices)
		sortedIndices := []int{indices[0], indices[1], indices[2]}
		sort.Sort(sort.Reverse(sort.IntSlice(sortedIndices)))

		// Put traded cards back in the deck
		for _, idx := range sortedIndices {
			state.Deck = append(state.Deck, player.Cards[idx])
		}
		for _, idx := range sortedIndices {
			player.Cards = append(player.Cards[:idx], player.Cards[idx+1:]...)
		}

		// Place bonus troops on highest-scored owned territory
		owned := getOwnedTerritories(state, player.ID)
		if len(owned) > 0 {
			bestID := pickBestTerritory(state, owned, player.ID)
			t := state.Territories[bestID]
			t.Troops += bonus
			state.Territories[bestID] = t
			addLog(state, player.ID, fmt.Sprintf("%s placed %d bonus troops on %s", name, bonus, t.Name))
		}

		// Only trade once per loop iteration if we still have 4+
		if len(player.Cards) < 4 {
			break
		}
	}
}

// cardTradeBonus returns the troop bonus for the Nth card trade in the game.
func cardTradeBonus(tradeNumber int) int {
	// Standard Risk escalating trade values
	switch tradeNumber {
	case 1:
		return 4
	case 2:
		return 6
	case 3:
		return 8
	case 4:
		return 10
	case 5:
		return 12
	case 6:
		return 15
	default:
		// After 6th trade, increase by 5 each time
		return 15 + (tradeNumber-6)*5
	}
}

// attackCandidate represents a potential attack.
type attackCandidate struct {
	fromID string
	toID   string
	score  float64
}

// executeAttacks handles the attack phase.
func executeAttacks(state *models.GameState, playerID string) {
	name := getPlayerName(state, playerID)
	maxAttacks := 20 // Safety limit to prevent infinite loops

	for attacks := 0; attacks < maxAttacks; attacks++ {
		candidates := buildAttackCandidates(state, playerID)
		if len(candidates) == 0 {
			break
		}

		// Sort by score descending
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score > candidates[j].score
		})

		// Pick the best attack (occasionally skip it for randomness)
		attacked := false
		for _, c := range candidates {
			// 15% chance to skip any given attack
			if rand.Float64() < 0.15 {
				continue
			}

			from := state.Territories[c.fromID]
			to := state.Territories[c.toID]
			defenderID := to.Owner // capture before attack mutates the state

			// Attack until we conquer or it's no longer favorable
			conquered := executeOneAttack(state, playerID, c.fromID, c.toID)

			if conquered {
				addLog(state, playerID, fmt.Sprintf("%s conquered %s from %s", name, to.Name, from.Name))
				state.ConqueredThisTurn = true

				// Check if we completed a continent
				if ownsContinent(state, to.Continent, playerID) {
					addLog(state, playerID, fmt.Sprintf("%s now controls all of %s!", name, to.Continent))
				}

				// Check if defending player is eliminated
				if defenderID != "" && defenderID != playerID {
					defenderAlive := false
					for _, t := range state.Territories {
						if t.Owner == defenderID {
							defenderAlive = true
							break
						}
					}
					if !defenderAlive {
						for i := range state.Players {
							if state.Players[i].ID == defenderID {
								state.Players[i].IsAlive = false
								defName := state.Players[i].Name
								addLog(state, playerID, fmt.Sprintf("%s eliminated %s!", name, defName))

								// Take eliminated player's cards
								player := getPlayer(state, playerID)
								defender := &state.Players[i]
								if player != nil {
									player.Cards = append(player.Cards, defender.Cards...)
									defender.Cards = nil
								}
								break
							}
						}
					}
				}
			}
			attacked = true
			break
		}

		if !attacked {
			break
		}
	}

	// Draw a card if conquered at least one territory
	if state.ConqueredThisTurn {
		drawCard(state, playerID)
	}
}

// buildAttackCandidates finds all favorable attack opportunities.
func buildAttackCandidates(state *models.GameState, playerID string) []attackCandidate {
	var candidates []attackCandidate

	for _, t := range state.Territories {
		if t.Owner != playerID || t.Troops < 2 {
			continue
		}

		for _, nID := range t.Neighbors {
			n, ok := state.Territories[nID]
			if !ok || n.Owner == playerID || n.Owner == "" {
				continue
			}

			// Only attack when we have >= 2x defender troops (favorable odds)
			attackerTroops := t.Troops - 1 // leave 1 behind
			if attackerTroops < n.Troops*2 {
				continue
			}

			// Score the attack
			score := float64(attackerTroops) / float64(maxInt(1, n.Troops))
			score += ScoreTerritory(state, nID, playerID) * 0.5

			// Bonus if this would complete a continent
			continent := n.Continent
			ratio := continentCompletionRatio(state, continent, playerID)
			if ratio >= 0.6 {
				score += 30 * ratio
			}

			candidates = append(candidates, attackCandidate{
				fromID: t.ID,
				toID:   nID,
				score:  score,
			})
		}
	}

	return candidates
}

// executeOneAttack carries out repeated dice rolls for a single attack until
// either the territory is conquered or the attack is no longer favorable.
// Returns true if the territory was conquered.
func executeOneAttack(state *models.GameState, playerID, fromID, toID string) bool {
	name := getPlayerName(state, playerID)

	for {
		from := state.Territories[fromID]
		to := state.Territories[toID]

		attackerAvailable := from.Troops - 1
		if attackerAvailable <= 0 {
			return false
		}

		// Stop if odds are no longer favorable (less than 1.5x)
		if float64(attackerAvailable) < float64(to.Troops)*1.5 {
			return false
		}

		// Roll dice
		attackDice := minInt(3, attackerAvailable)
		defendDice := minInt(2, to.Troops)

		attackRolls := rollDice(attackDice)
		defendRolls := rollDice(defendDice)

		// Compare dice
		attackerLosses := 0
		defenderLosses := 0
		comparisons := minInt(len(attackRolls), len(defendRolls))
		for i := 0; i < comparisons; i++ {
			if attackRolls[i] > defendRolls[i] {
				defenderLosses++
			} else {
				attackerLosses++
			}
		}

		// Record attack result
		state.LastAttackResult = &models.AttackResult{
			AttackerRolls:      attackRolls,
			DefenderRolls:      defendRolls,
			AttackerLosses:     attackerLosses,
			DefenderLosses:     defenderLosses,
			Conquered:          false,
			AttackingTerritory: fromID,
			DefendingTerritory: toID,
		}

		// Apply losses
		from.Troops -= attackerLosses
		to.Troops -= defenderLosses

		addLog(state, playerID, fmt.Sprintf(
			"%s attacked %s from %s (rolled %v vs %v)",
			name, to.Name, from.Name, attackRolls, defendRolls,
		))

		// Check for conquest
		if to.Troops <= 0 {
			// Conquered! Move troops in
			moveTroops := from.Troops - 1
			if moveTroops < 1 {
				moveTroops = 1
			}
			from.Troops -= moveTroops
			if from.Troops < 1 {
				from.Troops = 1
			}
			to.Owner = playerID
			to.Troops = moveTroops

			state.Territories[fromID] = from
			state.Territories[toID] = to
			state.LastAttackResult.Conquered = true

			return true
		}

		state.Territories[fromID] = from
		state.Territories[toID] = to
	}
}

// drawCard draws a card from the deck for the player (if available).
func drawCard(state *models.GameState, playerID string) {
	if len(state.Deck) == 0 {
		return
	}
	player := getPlayer(state, playerID)
	if player == nil {
		return
	}
	card := state.Deck[len(state.Deck)-1]
	state.Deck = state.Deck[:len(state.Deck)-1]
	player.Cards = append(player.Cards, card)
	addLog(state, playerID, fmt.Sprintf("%s drew a card", player.Name))
}

// executeFortify handles the fortify phase — moving troops from a safe
// interior territory to a threatened border territory.
func executeFortify(state *models.GameState, playerID string) {
	name := getPlayerName(state, playerID)

	borders := getBorderTerritories(state, playerID)
	if len(borders) == 0 {
		return
	}

	// Find the most threatened border territory
	var mostThreatened models.Territory
	highestThreat := 0.0
	for _, b := range borders {
		threat := 0.0
		for _, nID := range b.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != playerID {
				threat += float64(n.Troops)
			}
		}
		if threat > highestThreat {
			highestThreat = threat
			mostThreatened = b
		}
	}

	if mostThreatened.ID == "" {
		return
	}

	// Find an interior or safe territory with excess troops
	owned := getOwnedTerritories(state, playerID)
	var bestSource models.Territory
	bestExcess := 0

	for _, t := range owned {
		if t.ID == mostThreatened.ID {
			continue
		}
		if t.Troops <= 1 {
			continue
		}

		// Prefer interior territories (all neighbors owned by us)
		isInterior := true
		for _, nID := range t.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != playerID {
				isInterior = false
				break
			}
		}

		excess := t.Troops - 1
		effectiveExcess := excess
		if isInterior {
			effectiveExcess += 10 // strongly prefer interior territories
		}

		if effectiveExcess > bestExcess {
			// Check connectivity
			if isConnected(state, t.ID, mostThreatened.ID, playerID) {
				bestExcess = effectiveExcess
				bestSource = t
			}
		}
	}

	if bestSource.ID == "" || bestSource.Troops <= 1 {
		return
	}

	moveTroops := bestSource.Troops - 1

	src := state.Territories[bestSource.ID]
	dst := state.Territories[mostThreatened.ID]
	src.Troops -= moveTroops
	dst.Troops += moveTroops
	state.Territories[bestSource.ID] = src
	state.Territories[mostThreatened.ID] = dst

	addLog(state, playerID, fmt.Sprintf(
		"%s fortified %s with %d troops from %s",
		name, dst.Name, moveTroops, src.Name,
	))
}
