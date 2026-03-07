package game

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/mattheweckstein/risk/backend/models"
)

// GameEngine manages Risk game state and enforces rules.
type GameEngine struct{}

// NewGameEngine creates a new GameEngine instance.
func NewGameEngine() *GameEngine {
	return &GameEngine{}
}

// NewGame creates a new game with one human player and the specified number of AI players.
func (e *GameEngine) NewGame(playerName string, aiCount int, aiNames []string) *models.GameState {
	colors := []string{"red", "blue", "green", "yellow"}
	totalPlayers := 1 + aiCount
	if totalPlayers < 2 {
		totalPlayers = 2
		aiCount = 1
	}
	if totalPlayers > 4 {
		totalPlayers = 4
		aiCount = totalPlayers - 1
	}

	// Determine starting troops
	var startingTroops int
	switch totalPlayers {
	case 2:
		startingTroops = 40
	case 3:
		startingTroops = 35
	case 4:
		startingTroops = 30
	default:
		startingTroops = 30
	}

	// Create players
	players := make([]models.Player, 0, totalPlayers)
	players = append(players, models.Player{
		ID:      "player_0",
		Name:    playerName,
		IsAI:    false,
		Color:   colors[0],
		Cards:   []models.Card{},
		IsAlive: true,
	})
	for i := 0; i < aiCount; i++ {
		name := fmt.Sprintf("AI %d", i+1)
		if i < len(aiNames) && aiNames[i] != "" {
			name = aiNames[i]
		}
		players = append(players, models.Player{
			ID:      fmt.Sprintf("player_%d", i+1),
			Name:    name,
			IsAI:    true,
			Color:   colors[i+1],
			Cards:   []models.Card{},
			IsAlive: true,
		})
	}

	// Build territories map
	territories := make(map[string]models.Territory, len(models.AllTerritories))
	terrIDs := make([]string, 0, len(models.AllTerritories))
	for _, td := range models.AllTerritories {
		territories[td.ID] = models.Territory{
			ID:        td.ID,
			Name:      td.Name,
			Continent: td.Continent,
			Neighbors: td.Neighbors,
			Troops:    0,
		}
		terrIDs = append(terrIDs, td.ID)
	}

	// Shuffle and distribute territories
	rand.Shuffle(len(terrIDs), func(i, j int) {
		terrIDs[i], terrIDs[j] = terrIDs[j], terrIDs[i]
	})
	for i, tid := range terrIDs {
		t := territories[tid]
		t.Owner = players[i%totalPlayers].ID
		t.Troops = 1
		territories[tid] = t
	}

	// Build card deck
	deck := buildDeck()

	// Calculate troops to deploy for first player (setup: each placed 1 per territory already)
	// Count territories per player to know how many troops are already placed
	terrCounts := make(map[string]int)
	for _, t := range territories {
		terrCounts[t.Owner]++
	}

	state := &models.GameState{
		Phase:          models.PhaseSetup,
		Turn:           0,
		CurrentPlayer:  players[0].ID,
		Players:        players,
		Territories:    territories,
		Deck:           deck,
		Log:            []models.LogEntry{},
		TroopsToDeploy: startingTroops - terrCounts[players[0].ID],
		CardTradeCount: 0,
	}

	state.Log = append(state.Log, models.LogEntry{
		Turn:    0,
		Player:  "",
		Message: fmt.Sprintf("Game started with %d players. Each player has %d starting troops.", totalPlayers, startingTroops),
	})

	return state
}

// buildDeck creates a shuffled Risk card deck.
func buildDeck() []models.Card {
	types := []string{"infantry", "cavalry", "artillery"}
	cards := make([]models.Card, 0, len(models.AllTerritories)+2)
	for i, td := range models.AllTerritories {
		cards = append(cards, models.Card{
			Territory: td.ID,
			Type:      types[i%3],
		})
	}
	// Add 2 wild cards
	cards = append(cards, models.Card{Territory: "", Type: "wild"})
	cards = append(cards, models.Card{Territory: "", Type: "wild"})

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// PlaceTroops places troops on an owned territory during setup or place phase.
func (e *GameEngine) PlaceTroops(state *models.GameState, territory string, troops int) error {
	if state.Phase != models.PhaseSetup && state.Phase != models.PhasePlace {
		return fmt.Errorf("cannot place troops during %s phase", state.Phase)
	}

	if troops <= 0 {
		return fmt.Errorf("must place at least 1 troop")
	}

	if troops > state.TroopsToDeploy {
		return fmt.Errorf("not enough troops to deploy: have %d, want %d", state.TroopsToDeploy, troops)
	}

	t, ok := state.Territories[territory]
	if !ok {
		return fmt.Errorf("territory %s does not exist", territory)
	}

	if t.Owner != state.CurrentPlayer {
		return fmt.Errorf("territory %s is not owned by current player", territory)
	}

	if state.Phase == models.PhaseSetup && troops != 1 {
		return fmt.Errorf("during setup, must place exactly 1 troop at a time")
	}

	// Place troops
	t.Troops += troops
	state.Territories[territory] = t
	state.TroopsToDeploy -= troops

	player := e.getPlayer(state, state.CurrentPlayer)
	state.Log = append(state.Log, models.LogEntry{
		Turn:    state.Turn,
		Player:  state.CurrentPlayer,
		Message: fmt.Sprintf("%s placed %d troop(s) on %s", player.Name, troops, territory),
	})

	if state.Phase == models.PhaseSetup && state.TroopsToDeploy == 0 {
		// Advance to next player in setup
		e.advanceSetupPlayer(state)
	}

	if state.Phase == models.PhasePlace && state.TroopsToDeploy == 0 {
		state.Phase = models.PhaseAttack
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: "All troops placed. Entering attack phase.",
		})
	}

	return nil
}

// advanceSetupPlayer advances to the next player in setup phase.
func (e *GameEngine) advanceSetupPlayer(state *models.GameState) {
	nextIdx := e.nextAlivePlayerIndex(state)
	nextPlayer := state.Players[nextIdx]

	// Check if all players have finished setup
	allDone := true
	for _, p := range state.Players {
		if !p.IsAlive {
			continue
		}
		remaining := e.setupTroopsRemaining(state, p.ID)
		if remaining > 0 {
			allDone = false
			break
		}
	}

	if allDone {
		// Transition to first player's place phase
		state.Phase = models.PhasePlace
		state.CurrentPlayer = state.Players[0].ID
		state.Turn = 1
		state.TroopsToDeploy = CalculateTroops(state, state.CurrentPlayer)
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("Setup complete. %s's turn begins with %d troops to place.", state.Players[0].Name, state.TroopsToDeploy),
		})
		return
	}

	state.CurrentPlayer = nextPlayer.ID
	remaining := e.setupTroopsRemaining(state, nextPlayer.ID)
	if remaining > 0 {
		state.TroopsToDeploy = 1 // Setup: place 1 at a time
	} else {
		// This player is done, skip to next
		state.TroopsToDeploy = 0
		e.advanceSetupPlayer(state)
	}
}

// setupTroopsRemaining returns how many troops a player has left to place in setup.
func (e *GameEngine) setupTroopsRemaining(state *models.GameState, playerID string) int {
	totalPlayers := 0
	for _, p := range state.Players {
		if p.IsAlive {
			totalPlayers++
		}
	}

	var startingTroops int
	switch totalPlayers {
	case 2:
		startingTroops = 40
	case 3:
		startingTroops = 35
	case 4:
		startingTroops = 30
	default:
		startingTroops = 30
	}

	placed := 0
	for _, t := range state.Territories {
		if t.Owner == playerID {
			placed += t.Troops
		}
	}

	remaining := startingTroops - placed
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

// Attack resolves combat between two territories.
func (e *GameEngine) Attack(state *models.GameState, from, to string, attackerDice int) (*models.AttackResult, error) {
	if state.Phase != models.PhaseAttack {
		return nil, fmt.Errorf("cannot attack during %s phase", state.Phase)
	}

	fromT, ok := state.Territories[from]
	if !ok {
		return nil, fmt.Errorf("territory %s does not exist", from)
	}
	toT, ok := state.Territories[to]
	if !ok {
		return nil, fmt.Errorf("territory %s does not exist", to)
	}

	if fromT.Owner != state.CurrentPlayer {
		return nil, fmt.Errorf("you do not own %s", from)
	}
	if toT.Owner == state.CurrentPlayer {
		return nil, fmt.Errorf("cannot attack your own territory %s", to)
	}

	// Check adjacency
	if !isAdjacent(fromT, to) {
		return nil, fmt.Errorf("%s is not adjacent to %s", from, to)
	}

	if attackerDice < 1 || attackerDice > 3 {
		return nil, fmt.Errorf("attacker must roll 1-3 dice")
	}
	if fromT.Troops <= attackerDice {
		return nil, fmt.Errorf("need more than %d troops in %s to attack with %d dice", attackerDice, from, attackerDice)
	}

	// Roll dice
	attackerRolls := rollDice(attackerDice)
	defenderDice := toT.Troops
	if defenderDice > 2 {
		defenderDice = 2
	}
	defenderRolls := rollDice(defenderDice)

	// Compare pairs
	attackerLosses := 0
	defenderLosses := 0
	pairs := len(attackerRolls)
	if len(defenderRolls) < pairs {
		pairs = len(defenderRolls)
	}
	for i := 0; i < pairs; i++ {
		if attackerRolls[i] > defenderRolls[i] {
			defenderLosses++
		} else {
			attackerLosses++
		}
	}

	// Apply losses
	fromT.Troops -= attackerLosses
	toT.Troops -= defenderLosses

	conquered := false
	if toT.Troops <= 0 {
		conquered = true
		defenderID := toT.Owner
		toT.Owner = state.CurrentPlayer
		toT.Troops = attackerDice
		fromT.Troops -= attackerDice
		state.ConqueredThisTurn = true

		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("Conquered %s from %s!", toT.Name, from),
		})

		// Check if defender is eliminated
		e.checkPlayerEliminated(state, defenderID)
	}

	state.Territories[from] = fromT
	state.Territories[to] = toT

	result := &models.AttackResult{
		AttackerRolls:      attackerRolls,
		DefenderRolls:      defenderRolls,
		AttackerLosses:     attackerLosses,
		DefenderLosses:     defenderLosses,
		Conquered:          conquered,
		AttackingTerritory: from,
		DefendingTerritory: to,
	}
	state.LastAttackResult = result

	attackerPlayer := e.getPlayer(state, state.CurrentPlayer)
	state.Log = append(state.Log, models.LogEntry{
		Turn:   state.Turn,
		Player: state.CurrentPlayer,
		Message: fmt.Sprintf("%s attacked %s from %s. Attacker rolled %v, Defender rolled %v. Attacker lost %d, Defender lost %d.",
			attackerPlayer.Name, to, from, attackerRolls, defenderRolls, attackerLosses, defenderLosses),
	})

	// Check win condition
	if e.checkWinCondition(state) {
		state.Phase = models.PhaseEnded
		state.Winner = state.CurrentPlayer
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("%s wins the game!", attackerPlayer.Name),
		})
	}

	return result, nil
}

// Fortify moves troops between connected owned territories.
func (e *GameEngine) Fortify(state *models.GameState, from, to string, troops int) error {
	if state.Phase != models.PhaseFortify {
		return fmt.Errorf("cannot fortify during %s phase", state.Phase)
	}

	if troops <= 0 {
		return fmt.Errorf("must move at least 1 troop")
	}

	fromT, ok := state.Territories[from]
	if !ok {
		return fmt.Errorf("territory %s does not exist", from)
	}
	toT, ok := state.Territories[to]
	if !ok {
		return fmt.Errorf("territory %s does not exist", to)
	}

	if fromT.Owner != state.CurrentPlayer {
		return fmt.Errorf("you do not own %s", from)
	}
	if toT.Owner != state.CurrentPlayer {
		return fmt.Errorf("you do not own %s", to)
	}

	if fromT.Troops-troops < 1 {
		return fmt.Errorf("must leave at least 1 troop in %s", from)
	}

	if !IsConnected(state, from, to, state.CurrentPlayer) {
		return fmt.Errorf("%s and %s are not connected through your territories", from, to)
	}

	fromT.Troops -= troops
	toT.Troops += troops
	state.Territories[from] = fromT
	state.Territories[to] = toT

	player := e.getPlayer(state, state.CurrentPlayer)
	state.Log = append(state.Log, models.LogEntry{
		Turn:    state.Turn,
		Player:  state.CurrentPlayer,
		Message: fmt.Sprintf("%s moved %d troop(s) from %s to %s", player.Name, troops, from, to),
	})

	return nil
}

// EndPhase advances the game to the next phase.
func (e *GameEngine) EndPhase(state *models.GameState) error {
	switch state.Phase {
	case models.PhasePlace:
		if state.TroopsToDeploy > 0 {
			return fmt.Errorf("must place all %d remaining troops before ending place phase", state.TroopsToDeploy)
		}
		state.Phase = models.PhaseAttack
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: "Entering attack phase.",
		})

	case models.PhaseAttack:
		state.Phase = models.PhaseFortify
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: "Entering fortify phase.",
		})

	case models.PhaseFortify:
		// Award card if conquered this turn
		if state.ConqueredThisTurn && len(state.Deck) > 0 {
			card := state.Deck[0]
			state.Deck = state.Deck[1:]
			for i := range state.Players {
				if state.Players[i].ID == state.CurrentPlayer {
					state.Players[i].Cards = append(state.Players[i].Cards, card)
					break
				}
			}
			state.Log = append(state.Log, models.LogEntry{
				Turn:    state.Turn,
				Player:  state.CurrentPlayer,
				Message: "Earned a card for conquering a territory this turn.",
			})
		}
		state.ConqueredThisTurn = false
		state.LastAttackResult = nil

		// Advance to next player
		nextIdx := e.nextAlivePlayerIndex(state)
		state.CurrentPlayer = state.Players[nextIdx].ID
		state.Turn++
		state.Phase = models.PhasePlace
		state.TroopsToDeploy = CalculateTroops(state, state.CurrentPlayer)

		// Check forced card trade-in (5+ cards)
		for i := range state.Players {
			if state.Players[i].ID == state.CurrentPlayer && len(state.Players[i].Cards) >= 5 {
				state.Log = append(state.Log, models.LogEntry{
					Turn:    state.Turn,
					Player:  state.CurrentPlayer,
					Message: fmt.Sprintf("%s has %d cards and must trade in.", state.Players[i].Name, len(state.Players[i].Cards)),
				})
				break
			}
		}

		nextPlayer := e.getPlayer(state, state.CurrentPlayer)
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("%s's turn begins with %d troops to place.", nextPlayer.Name, state.TroopsToDeploy),
		})

	default:
		return fmt.Errorf("cannot end phase during %s", state.Phase)
	}

	return nil
}

// TradeCards validates and processes a card trade-in.
func (e *GameEngine) TradeCards(state *models.GameState, indices [3]int) error {
	if state.Phase != models.PhasePlace {
		return fmt.Errorf("can only trade cards during place phase")
	}

	var currentPlayer *models.Player
	for i := range state.Players {
		if state.Players[i].ID == state.CurrentPlayer {
			currentPlayer = &state.Players[i]
			break
		}
	}
	if currentPlayer == nil {
		return fmt.Errorf("current player not found")
	}

	// Validate indices
	for _, idx := range indices {
		if idx < 0 || idx >= len(currentPlayer.Cards) {
			return fmt.Errorf("card index %d out of range", idx)
		}
	}
	if indices[0] == indices[1] || indices[0] == indices[2] || indices[1] == indices[2] {
		return fmt.Errorf("card indices must be unique")
	}

	cards := [3]models.Card{
		currentPlayer.Cards[indices[0]],
		currentPlayer.Cards[indices[1]],
		currentPlayer.Cards[indices[2]],
	}

	if !isValidCardSet(cards) {
		return fmt.Errorf("invalid card combination")
	}

	// Calculate bonus
	bonus := cardTradeBonus(state.CardTradeCount)
	state.CardTradeCount++
	state.TroopsToDeploy += bonus

	state.Log = append(state.Log, models.LogEntry{
		Turn:    state.Turn,
		Player:  state.CurrentPlayer,
		Message: fmt.Sprintf("%s traded cards for %d bonus troops.", currentPlayer.Name, bonus),
	})

	// Check if traded card matches owned territory -> +2 troops
	for _, card := range cards {
		if card.Territory == "" {
			continue
		}
		if t, ok := state.Territories[card.Territory]; ok && t.Owner == state.CurrentPlayer {
			t.Troops += 2
			state.Territories[card.Territory] = t
			state.Log = append(state.Log, models.LogEntry{
				Turn:    state.Turn,
				Player:  state.CurrentPlayer,
				Message: fmt.Sprintf("Bonus: +2 troops on %s (matched traded card).", card.Territory),
			})
		}
	}

	// Remove cards (sort indices descending to remove from end first)
	sortedIndices := []int{indices[0], indices[1], indices[2]}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedIndices)))
	for _, idx := range sortedIndices {
		currentPlayer.Cards = append(currentPlayer.Cards[:idx], currentPlayer.Cards[idx+1:]...)
	}

	return nil
}

// CalculateTroops calculates the number of troops a player receives at the start of their turn.
func CalculateTroops(state *models.GameState, playerID string) int {
	// Count territories
	terrCount := 0
	for _, t := range state.Territories {
		if t.Owner == playerID {
			terrCount++
		}
	}

	troops := terrCount / 3
	if troops < 3 {
		troops = 3
	}

	// Continent bonuses
	for continent, terrIDs := range models.ContinentTerritories {
		ownsAll := true
		for _, tid := range terrIDs {
			if t, ok := state.Territories[tid]; !ok || t.Owner != playerID {
				ownsAll = false
				break
			}
		}
		if ownsAll {
			troops += models.ContinentBonuses[continent]
		}
	}

	return troops
}

// IsConnected checks if two territories are connected through territories owned by the same player using BFS.
func IsConnected(state *models.GameState, from, to, playerID string) bool {
	if from == to {
		return true
	}

	visited := make(map[string]bool)
	queue := []string{from}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		t, ok := state.Territories[current]
		if !ok {
			continue
		}

		for _, neighbor := range t.Neighbors {
			if visited[neighbor] {
				continue
			}
			nt, ok := state.Territories[neighbor]
			if !ok || nt.Owner != playerID {
				continue
			}
			if neighbor == to {
				return true
			}
			visited[neighbor] = true
			queue = append(queue, neighbor)
		}
	}

	return false
}

// CheckElimination marks players with 0 territories as not alive.
func CheckElimination(state *models.GameState) {
	terrCounts := make(map[string]int)
	for _, t := range state.Territories {
		terrCounts[t.Owner]++
	}
	for i := range state.Players {
		if terrCounts[state.Players[i].ID] == 0 {
			state.Players[i].IsAlive = false
		}
	}
}

// checkPlayerEliminated checks if a specific player has been eliminated and transfers their cards.
func (e *GameEngine) checkPlayerEliminated(state *models.GameState, playerID string) {
	count := 0
	for _, t := range state.Territories {
		if t.Owner == playerID {
			count++
		}
	}
	if count > 0 {
		return
	}

	// Mark as eliminated
	var eliminatedPlayer *models.Player
	for i := range state.Players {
		if state.Players[i].ID == playerID {
			state.Players[i].IsAlive = false
			eliminatedPlayer = &state.Players[i]
			break
		}
	}

	if eliminatedPlayer == nil {
		return
	}

	// Transfer cards to the conquering player
	if len(eliminatedPlayer.Cards) > 0 {
		for i := range state.Players {
			if state.Players[i].ID == state.CurrentPlayer {
				state.Players[i].Cards = append(state.Players[i].Cards, eliminatedPlayer.Cards...)
				break
			}
		}
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("%s eliminated! Transferred %d cards.", eliminatedPlayer.Name, len(eliminatedPlayer.Cards)),
		})
		eliminatedPlayer.Cards = []models.Card{}
	} else {
		state.Log = append(state.Log, models.LogEntry{
			Turn:    state.Turn,
			Player:  state.CurrentPlayer,
			Message: fmt.Sprintf("%s has been eliminated!", eliminatedPlayer.Name),
		})
	}
}

// checkWinCondition checks if the current player owns all territories.
func (e *GameEngine) checkWinCondition(state *models.GameState) bool {
	for _, t := range state.Territories {
		if t.Owner != state.CurrentPlayer {
			return false
		}
	}
	return true
}

// rollDice rolls n dice and returns the results sorted descending.
func rollDice(n int) []int {
	rolls := make([]int, n)
	for i := range rolls {
		rolls[i] = rand.Intn(6) + 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(rolls)))
	return rolls
}

// isAdjacent checks if a territory is adjacent to another by ID.
func isAdjacent(t models.Territory, neighborID string) bool {
	for _, n := range t.Neighbors {
		if n == neighborID {
			return true
		}
	}
	return false
}

// isValidCardSet checks if three cards form a valid trade-in set.
func isValidCardSet(cards [3]models.Card) bool {
	types := make(map[string]int)
	wildCount := 0
	for _, c := range cards {
		if c.Type == "wild" {
			wildCount++
		} else {
			types[c.Type]++
		}
	}

	// Any set with a wild card is valid
	if wildCount > 0 {
		return true
	}

	// Three of the same type
	for _, count := range types {
		if count == 3 {
			return true
		}
	}

	// One of each type
	if len(types) == 3 {
		return true
	}

	return false
}

// cardTradeBonus returns the bonus for the nth card trade (0-indexed).
func cardTradeBonus(tradeCount int) int {
	bonuses := []int{4, 6, 8, 10, 12, 15}
	if tradeCount < len(bonuses) {
		return bonuses[tradeCount]
	}
	// After the 6th trade, increase by 5 each time
	return 20 + (tradeCount-6)*5
}

// getPlayer returns a copy of the player with the given ID.
func (e *GameEngine) getPlayer(state *models.GameState, playerID string) models.Player {
	for _, p := range state.Players {
		if p.ID == playerID {
			return p
		}
	}
	return models.Player{}
}

// nextAlivePlayerIndex returns the index of the next alive player after the current player.
func (e *GameEngine) nextAlivePlayerIndex(state *models.GameState) int {
	currentIdx := -1
	for i, p := range state.Players {
		if p.ID == state.CurrentPlayer {
			currentIdx = i
			break
		}
	}

	for offset := 1; offset <= len(state.Players); offset++ {
		idx := (currentIdx + offset) % len(state.Players)
		if state.Players[idx].IsAlive {
			return idx
		}
	}
	return currentIdx
}
