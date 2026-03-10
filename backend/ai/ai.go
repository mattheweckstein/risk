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
		// Pick from top 2 to keep some variety but be more strategic
		top := minInt(2, len(candidates))
		return candidates[rand.Intn(top)].id
	}

	// All territories claimed — reinforce an owned territory
	owned := getOwnedTerritories(state, playerID)
	if len(owned) == 0 {
		for id := range state.Territories {
			return id
		}
	}

	return pickBestSetupReinforcement(state, owned, playerID)
}

// pickBestSetupReinforcement picks the best territory to reinforce during setup.
// Focuses on border territories, especially those near continent completion.
func pickBestSetupReinforcement(state *models.GameState, territories []models.Territory, playerID string) string {
	type scored struct {
		id    string
		score float64
	}
	var candidates []scored
	for _, t := range territories {
		s := scorePlacement(state, t.ID, playerID)
		candidates = append(candidates, scored{t.ID, s})
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	if len(candidates) == 0 {
		return territories[0].ID
	}

	// Mostly pick the best, small chance of second best
	if len(candidates) >= 2 && rand.Float64() < 0.2 {
		return candidates[1].id
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

	// Get border territories for placement
	borders := getBorderTerritories(state, playerID)
	if len(borders) == 0 {
		borders = owned
	}

	// Score borders using the strategic placement scorer
	type scored struct {
		id    string
		score float64
	}
	var candidates []scored
	for _, t := range borders {
		s := scorePlacement(state, t.ID, playerID)
		candidates = append(candidates, scored{t.ID, s})
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	// Concentrate troops more aggressively on top candidates
	// Put 60% on best, 25% on second, 15% on third
	remaining := troops
	topCount := minInt(len(candidates), 3)
	weights := []float64{0.60, 0.25, 0.15}

	for i := 0; i < topCount && remaining > 0; i++ {
		place := int(float64(troops) * weights[i])
		if i == topCount-1 {
			// Last one gets everything remaining
			place = remaining
		}
		if place > remaining {
			place = remaining
		}
		if place < 1 && remaining > 0 {
			place = 1
		}

		tID := candidates[i].id
		t := state.Territories[tID]
		t.Troops += place
		state.Territories[tID] = t
		remaining -= place

		addLog(state, playerID, fmt.Sprintf("%s placed %d troops on %s", name, place, t.Name))
	}

	// If any remaining (shouldn't happen, but safety)
	if remaining > 0 && len(candidates) > 0 {
		tID := candidates[0].id
		t := state.Territories[tID]
		t.Troops += remaining
		state.Territories[tID] = t
		addLog(state, playerID, fmt.Sprintf("%s placed %d troops on %s", name, remaining, t.Name))
	}

	return nil
}

// tradeCards handles card trading for the AI.
func tradeCards(state *models.GameState, player *models.Player) {
	for len(player.Cards) >= 4 {
		// Strategic timing: if we have exactly 4 cards and the next bonus is small,
		// hold until 5 (forced trade) to let the global trade count go up
		if len(player.Cards) == 4 && state.CardTradeCount < 3 {
			// Early game — holding is worth it since bonus values increase
			break
		}

		indices, found := findValidCardSet(player.Cards)
		if !found {
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

		// Place bonus troops strategically -- on the best border territory
		borders := getBorderTerritories(state, player.ID)
		if len(borders) == 0 {
			borders = getOwnedTerritories(state, player.ID)
		}
		if len(borders) > 0 {
			bestID := pickBestPlacement(state, borders, player.ID)
			t := state.Territories[bestID]
			t.Troops += bonus
			state.Territories[bestID] = t
			addLog(state, player.ID, fmt.Sprintf("%s placed %d bonus troops on %s", name, bonus, t.Name))
		}

		if len(player.Cards) < 4 {
			break
		}
	}
}

// pickBestPlacement picks the best territory for placing troops using the
// strategic placement scorer. Deterministic -- always picks the best.
func pickBestPlacement(state *models.GameState, territories []models.Territory, playerID string) string {
	bestID := territories[0].ID
	bestScore := scorePlacement(state, territories[0].ID, playerID)

	for _, t := range territories[1:] {
		s := scorePlacement(state, t.ID, playerID)
		if s > bestScore {
			bestScore = s
			bestID = t.ID
		}
	}
	return bestID
}

// cardTradeBonus returns the troop bonus for the Nth card trade in the game.
func cardTradeBonus(tradeNumber int) int {
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
		return 15 + (tradeNumber-6)*5
	}
}

// attackCandidate represents a potential attack.
type attackCandidate struct {
	fromID string
	toID   string
	score  float64
}

// executeAttacks handles the attack phase with strategic goal-driven behavior.
func executeAttacks(state *models.GameState, playerID string) {
	name := getPlayerName(state, playerID)
	maxAttacks := 30 // Increased limit for more aggressive play

	for attacks := 0; attacks < maxAttacks; attacks++ {
		candidates := buildAttackCandidates(state, playerID)
		if len(candidates) == 0 {
			break
		}

		// Sort by score descending
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score > candidates[j].score
		})

		// Pick the best attack — almost never skip a scored attack
		attacked := false
		for _, c := range candidates {

			from := state.Territories[c.fromID]
			to := state.Territories[c.toID]
			defenderID := to.Owner

			conquered := executeOneAttack(state, playerID, c.fromID, c.toID)

			if conquered {
				addLog(state, playerID, fmt.Sprintf("%s conquered %s from %s", name, to.Name, from.Name))
				state.ConqueredThisTurn = true

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

								player := getPlayer(state, playerID)
								defender := &state.Players[i]
								if player != nil {
									player.Cards = append(player.Cards, defender.Cards...)
									defender.Cards = nil
									tradeCards(state, player)
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
}

// buildAttackCandidates finds all viable attack opportunities with strategic scoring.
func buildAttackCandidates(state *models.GameState, playerID string) []attackCandidate {
	var candidates []attackCandidate

	// Are we behind the leader? If so, be more aggressive
	myStrength := getPlayerStrength(state, playerID)
	leaderID := getLeaderID(state)
	leaderStrength := getPlayerStrength(state, leaderID)
	behindLeader := leaderID != playerID && leaderStrength > myStrength
	strengthGap := leaderStrength - myStrength

	// Base ratio adjustment: the further behind, the more aggressive
	aggressionMod := 0.0
	if behindLeader {
		aggressionMod = strengthGap * 0.01 // Up to ~0.3 reduction in required ratio
		if aggressionMod > 0.3 {
			aggressionMod = 0.3
		}
	}

	for _, t := range state.Territories {
		if t.Owner != playerID || t.Troops < 2 {
			continue
		}

		for _, nID := range t.Neighbors {
			n, ok := state.Territories[nID]
			if !ok || n.Owner == playerID || n.Owner == "" {
				continue
			}

			attackerTroops := t.Troops - 1
			defenderTroops := maxInt(1, n.Troops)
			ratio := float64(attackerTroops) / float64(defenderTroops)

			// Determine minimum ratio based on context
			minRatio := 1.4 - aggressionMod // Default: 1.4x, lower when behind

			// Check if this attack would complete a continent
			wouldComplete := false
			continent := n.Continent
			needed := 0
			for _, tID := range models.ContinentTerritories[continent] {
				if ct, ok2 := state.Territories[tID]; ok2 && ct.Owner != playerID {
					needed++
				}
			}
			if needed == 1 {
				wouldComplete = true
				minRatio = 1.0 // Very aggressive for continent completion
			}

			// Check if this would break an opponent's continent
			wouldBreak := false
			if ownsContinent(state, continent, n.Owner) {
				wouldBreak = true
				minRatio = 1.0 // Always worth breaking a continent
			}

			// Check if this would eliminate a player
			wouldEliminate := false
			defenderCount := countPlayerTerritories(state, n.Owner)
			defenderCards := countPlayerCards(state, n.Owner)
			if defenderCount == 1 {
				wouldEliminate = true
				minRatio = 0.9 // Even slightly bad odds acceptable to eliminate
				if defenderCards >= 3 {
					minRatio = 0.7 // Worth a real gamble for cards
				}
			}

			// Extra aggression against the leader
			targetingLeader := n.Owner == leaderID && leaderID != playerID
			if targetingLeader {
				minRatio -= 0.2
			}

			// Floor the ratio at 0.6 — don't suicide
			if minRatio < 0.6 {
				minRatio = 0.6
			}

			// Skip if below minimum ratio
			if ratio < minRatio {
				continue
			}

			// Score the attack using the strategic scorer
			score := scoreAttack(state, t.ID, nID, playerID)

			// Additional bonuses
			if wouldComplete {
				score += 50
			}
			if wouldBreak {
				score += 40
			}
			if wouldEliminate {
				score += 40
			}

			// Leader targeting bonus — all AIs focus the strongest player
			score += leaderThreatBonus(state, n.Owner, playerID)

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

	// Pre-compute strategic value to determine how aggressively to press the attack
	to := state.Territories[toID]
	continent := to.Continent
	isStrategic := false

	// Check continent completion
	needed := 0
	for _, tID := range models.ContinentTerritories[continent] {
		if ct, ok := state.Territories[tID]; ok && ct.Owner != playerID {
			needed++
		}
	}
	if needed == 1 {
		isStrategic = true
	}

	// Check opponent continent breaking
	if to.Owner != "" && ownsContinent(state, continent, to.Owner) {
		isStrategic = true
	}

	// Check player elimination
	if to.Owner != "" && countPlayerTerritories(state, to.Owner) == 1 {
		isStrategic = true
	}

	// Attacking the leader is always strategic
	leaderID := getLeaderID(state)
	targetingLeader := to.Owner == leaderID && leaderID != playerID
	if targetingLeader {
		isStrategic = true
	}

	for {
		from := state.Territories[fromID]
		to = state.Territories[toID]

		attackerAvailable := from.Troops - 1
		if attackerAvailable <= 0 {
			return false
		}

		// Determine stop threshold based on strategic value
		stopRatio := 1.2 // Default: continue while we have 1.2x advantage
		if isStrategic {
			stopRatio = 0.8 // Press much harder on strategic attacks
		}

		if float64(attackerAvailable) < float64(maxInt(1, to.Troops))*stopRatio {
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
			// Conquered! Move troops in strategically
			moveTroops := decideConquestMove(state, from, to, playerID)
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

// decideConquestMove decides how many troops to move into a conquered territory.
// Considers whether the conquered territory has enemy neighbors that threaten it,
// or whether the source territory still needs troops for further attacks.
func decideConquestMove(state *models.GameState, from, to models.Territory, playerID string) int {
	available := from.Troops - 1
	if available < 1 {
		return 1
	}

	// Count enemy threat around the conquered territory
	enemyThreatTo := 0
	for _, nID := range to.Neighbors {
		if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
			enemyThreatTo += n.Troops
		}
	}

	// Count enemy threat around the source territory (excluding the now-conquered one)
	enemyThreatFrom := 0
	for _, nID := range from.Neighbors {
		if nID == to.ID {
			continue
		}
		if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
			enemyThreatFrom += n.Troops
		}
	}

	// If source has no more enemies, move everything
	if enemyThreatFrom == 0 {
		return available
	}

	// If conquered territory has more enemy threat, move more troops there
	if enemyThreatTo > enemyThreatFrom {
		// Move most troops forward
		move := available * 3 / 4
		if move < 1 {
			move = 1
		}
		return move
	}

	// Default: move half, keep half for source defense
	move := available / 2
	if move < 1 {
		move = 1
	}
	return move
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

// executeFortify handles the fortify phase — strategically moves troops
// toward the most important front lines.
// In free fortify mode, makes multiple moves to optimize troop distribution.
func executeFortify(state *models.GameState, playerID string) {
	name := getPlayerName(state, playerID)

	// In free fortify mode, do up to 5 moves. In classic, do 1.
	maxMoves := 1
	if state.FreeFortify {
		maxMoves = 5
	}

	for move := 0; move < maxMoves; move++ {
		borders := getBorderTerritories(state, playerID)
		if len(borders) == 0 {
			return
		}

		// Score all border territories by strategic importance and need
		type fortifyTargetT struct {
			territory models.Territory
			score     float64
		}

		var targets []fortifyTargetT
		for _, b := range borders {
			enemyThreat := 0.0
			for _, nID := range b.Neighbors {
				if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
					enemyThreat += float64(n.Troops)
				}
			}

			stratScore := scorePlacement(state, b.ID, playerID)
			deficit := enemyThreat - float64(b.Troops)
			score := stratScore
			if deficit > 0 {
				score += deficit * 5
			}
			if ownsContinent(state, b.Continent, playerID) {
				score += 30 + float64(models.ContinentBonuses[b.Continent])*5
			}

			targets = append(targets, fortifyTargetT{b, score})
		}

		sort.Slice(targets, func(i, j int) bool {
			return targets[i].score > targets[j].score
		})

		if len(targets) == 0 {
			return
		}

		bestTarget := targets[0].territory

		// Find the best source: interior territory or low-priority border with excess troops
		owned := getOwnedTerritories(state, playerID)

		type fortifySourceT struct {
			territory models.Territory
			score     float64
		}

		var sources []fortifySourceT
		for _, t := range owned {
			if t.ID == bestTarget.ID || t.Troops <= 1 {
				continue
			}

			if !isConnected(state, t.ID, bestTarget.ID, playerID) {
				continue
			}

			isInterior := true
			localThreat := 0.0
			for _, nID := range t.Neighbors {
				if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
					isInterior = false
					localThreat += float64(n.Troops)
				}
			}

			excess := float64(t.Troops - 1)
			sourceScore := excess

			if isInterior {
				sourceScore += 100
			} else {
				safeExcess := float64(t.Troops) - localThreat*1.2
				if safeExcess <= 0 {
					continue
				}
				sourceScore += safeExcess
			}

			sources = append(sources, fortifySourceT{t, sourceScore})
		}

		sort.Slice(sources, func(i, j int) bool {
			return sources[i].score > sources[j].score
		})

		if len(sources) == 0 {
			return
		}

		bestSource := sources[0].territory

		// Determine how many troops to move
		moveTroops := bestSource.Troops - 1

		isSourceInterior := true
		sourceThreat := 0.0
		for _, nID := range bestSource.Neighbors {
			if n, ok := state.Territories[nID]; ok && n.Owner != playerID && n.Owner != "" {
				isSourceInterior = false
				sourceThreat += float64(n.Troops)
			}
		}
		if !isSourceInterior {
			keep := int(sourceThreat)
			if keep < 1 {
				keep = 1
			}
			moveTroops = bestSource.Troops - keep
			if moveTroops < 1 {
				if move == 0 {
					// Try next best source on first move
					continue
				}
				return
			}
		}

		src := state.Territories[bestSource.ID]
		dst := state.Territories[bestTarget.ID]
		src.Troops -= moveTroops
		dst.Troops += moveTroops
		state.Territories[bestSource.ID] = src
		state.Territories[bestTarget.ID] = dst

		addLog(state, playerID, fmt.Sprintf(
			"%s fortified %s with %d troops from %s",
			name, dst.Name, moveTroops, src.Name,
		))
	}
}
