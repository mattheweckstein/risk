package game

import (
	"testing"

	"github.com/mattheweckstein/risk/backend/models"
)

// newTestGame creates a minimal 2-player game in a known state for testing.
// Player 0 is human, player 1 is AI.
// All territories are distributed round-robin and have 1 troop each.
// The game is set to the setup phase.
func newTestGame() *models.GameState {
	e := NewGameEngine()
	state := e.NewGame("TestHuman", 1, []string{"TestAI"})
	return state
}

// newTestGameInPhase creates a 2-player game forced into a specific phase.
func newTestGameInPhase(phase models.Phase) *models.GameState {
	state := newTestGame()
	// Force all territories to have enough troops and skip setup
	state.Phase = phase
	state.Turn = 1
	state.CurrentPlayer = "player_0"
	state.TroopsToDeploy = 0

	if phase == models.PhasePlace {
		state.TroopsToDeploy = 5
	}

	// Give each territory some troops
	for id, t := range state.Territories {
		t.Troops = 3
		state.Territories[id] = t
	}

	return state
}

// findOwnedTerritory finds a territory owned by the given player.
func findOwnedTerritory(state *models.GameState, playerID string) string {
	for id, t := range state.Territories {
		if t.Owner == playerID {
			return id
		}
	}
	return ""
}

// findEnemyNeighbor finds an enemy neighbor of the given territory.
func findEnemyNeighbor(state *models.GameState, territoryID, playerID string) string {
	t := state.Territories[territoryID]
	for _, nID := range t.Neighbors {
		if n, ok := state.Territories[nID]; ok && n.Owner != playerID {
			return nID
		}
	}
	return ""
}

// findOwnedNeighbor finds a friendly neighbor of the given territory.
func findOwnedNeighbor(state *models.GameState, territoryID, playerID string) string {
	t := state.Territories[territoryID]
	for _, nID := range t.Neighbors {
		if n, ok := state.Territories[nID]; ok && n.Owner == playerID {
			return nID
		}
	}
	return ""
}

// findAttackPair finds a from/to pair where from is owned by playerID with enough troops and to is enemy adjacent.
func findAttackPair(state *models.GameState, playerID string) (from, to string) {
	for id, t := range state.Territories {
		if t.Owner == playerID && t.Troops >= 2 {
			enemy := findEnemyNeighbor(state, id, playerID)
			if enemy != "" {
				return id, enemy
			}
		}
	}
	return "", ""
}

// === TestNewGame ===

func TestNewGame(t *testing.T) {
	e := NewGameEngine()
	state := e.NewGame("Alice", 1, []string{"Bot"})

	if len(state.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(state.Players))
	}
	if state.Players[0].Name != "Alice" {
		t.Errorf("expected player 0 name 'Alice', got %q", state.Players[0].Name)
	}
	if state.Players[1].Name != "Bot" {
		t.Errorf("expected player 1 name 'Bot', got %q", state.Players[1].Name)
	}
	if !state.Players[1].IsAI {
		t.Error("player 1 should be AI")
	}
	if state.Players[0].IsAI {
		t.Error("player 0 should not be AI")
	}

	// All 42 territories should be distributed
	if len(state.Territories) != 42 {
		t.Fatalf("expected 42 territories, got %d", len(state.Territories))
	}

	p0Count, p1Count := 0, 0
	for _, terr := range state.Territories {
		if terr.Troops != 1 {
			t.Errorf("territory %s should have 1 troop, got %d", terr.ID, terr.Troops)
		}
		if terr.Owner == "player_0" {
			p0Count++
		} else if terr.Owner == "player_1" {
			p1Count++
		}
	}
	if p0Count+p1Count != 42 {
		t.Errorf("territories not fully distributed: p0=%d, p1=%d", p0Count, p1Count)
	}

	if state.Phase != models.PhaseSetup {
		t.Errorf("expected phase setup, got %s", state.Phase)
	}
	if state.CurrentPlayer != "player_0" {
		t.Errorf("expected current player player_0, got %s", state.CurrentPlayer)
	}

	// Deck should have 44 cards (42 territories + 2 wilds)
	if len(state.Deck) != 44 {
		t.Errorf("expected 44 cards in deck, got %d", len(state.Deck))
	}

	// Players should be alive
	for _, p := range state.Players {
		if !p.IsAlive {
			t.Errorf("player %s should be alive", p.ID)
		}
	}
}

func TestNewGamePlayerCounts(t *testing.T) {
	e := NewGameEngine()

	// 3 players
	state3 := e.NewGame("Alice", 2, nil)
	if len(state3.Players) != 3 {
		t.Errorf("expected 3 players, got %d", len(state3.Players))
	}

	// 4 players
	state4 := e.NewGame("Alice", 3, nil)
	if len(state4.Players) != 4 {
		t.Errorf("expected 4 players, got %d", len(state4.Players))
	}

	// Clamped at max 4
	stateMax := e.NewGame("Alice", 10, nil)
	if len(stateMax.Players) != 4 {
		t.Errorf("expected 4 players (clamped), got %d", len(stateMax.Players))
	}

	// Minimum 2
	stateMin := e.NewGame("Alice", 0, nil)
	if len(stateMin.Players) != 2 {
		t.Errorf("expected 2 players (minimum), got %d", len(stateMin.Players))
	}
}

// === TestPlaceTroops ===

func TestPlaceTroops(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	ownedID := findOwnedTerritory(state, "player_0")
	if ownedID == "" {
		t.Fatal("could not find owned territory")
	}

	origTroops := state.Territories[ownedID].Troops
	err := e.PlaceTroops(state, ownedID, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Territories[ownedID].Troops != origTroops+2 {
		t.Errorf("expected troops %d, got %d", origTroops+2, state.Territories[ownedID].Troops)
	}
	if state.TroopsToDeploy != 3 {
		t.Errorf("expected 3 troops remaining, got %d", state.TroopsToDeploy)
	}
}

func TestPlaceTroopsWrongOwner(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	enemyID := findOwnedTerritory(state, "player_1")
	if enemyID == "" {
		t.Fatal("could not find enemy territory")
	}

	err := e.PlaceTroops(state, enemyID, 1)
	if err == nil {
		t.Error("expected error when placing on enemy territory")
	}
}

func TestPlaceTroopsWrongPhase(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	ownedID := findOwnedTerritory(state, "player_0")
	err := e.PlaceTroops(state, ownedID, 1)
	if err == nil {
		t.Error("expected error when placing during attack phase")
	}
}

func TestPlaceTroopsTooMany(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	ownedID := findOwnedTerritory(state, "player_0")
	err := e.PlaceTroops(state, ownedID, 100)
	if err == nil {
		t.Error("expected error when placing more troops than available")
	}
}

func TestPlaceTroopsAutoAdvancesToAttack(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	state.TroopsToDeploy = 1
	e := NewGameEngine()

	ownedID := findOwnedTerritory(state, "player_0")
	err := e.PlaceTroops(state, ownedID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Phase != models.PhaseAttack {
		t.Errorf("expected phase attack after all troops placed, got %s", state.Phase)
	}
}

// === TestAttack ===

func TestAttack(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	// Give attacker extra troops
	from, to := findAttackPair(state, "player_0")
	if from == "" {
		t.Fatal("could not find attack pair")
	}
	f := state.Territories[from]
	f.Troops = 10
	state.Territories[from] = f

	result, err := e.Attack(state, from, to, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil attack result")
	}
	if len(result.AttackerRolls) != 3 {
		t.Errorf("expected 3 attacker dice, got %d", len(result.AttackerRolls))
	}
	if result.AttackingTerritory != from {
		t.Errorf("expected attacking territory %s, got %s", from, result.AttackingTerritory)
	}
	if result.DefendingTerritory != to {
		t.Errorf("expected defending territory %s, got %s", to, result.DefendingTerritory)
	}
}

func TestAttackNotAdjacent(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	// Find two non-adjacent territories owned by different players
	ownedID := findOwnedTerritory(state, "player_0")
	f := state.Territories[ownedID]
	f.Troops = 10
	state.Territories[ownedID] = f

	// Find an enemy territory that is NOT a neighbor
	neighbors := make(map[string]bool)
	for _, nID := range f.Neighbors {
		neighbors[nID] = true
	}
	var nonAdjacentEnemy string
	for id, terr := range state.Territories {
		if terr.Owner != "player_0" && !neighbors[id] && id != ownedID {
			nonAdjacentEnemy = id
			break
		}
	}
	if nonAdjacentEnemy == "" {
		t.Skip("could not find non-adjacent enemy territory")
	}

	_, err := e.Attack(state, ownedID, nonAdjacentEnemy, 1)
	if err == nil {
		t.Error("expected error for non-adjacent attack")
	}
}

func TestAttackOwnTerritory(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	ownedID := findOwnedTerritory(state, "player_0")
	f := state.Territories[ownedID]
	f.Troops = 10
	state.Territories[ownedID] = f

	friendlyNeighbor := findOwnedNeighbor(state, ownedID, "player_0")
	if friendlyNeighbor == "" {
		t.Skip("no friendly neighbor found")
	}

	_, err := e.Attack(state, ownedID, friendlyNeighbor, 1)
	if err == nil {
		t.Error("expected error for attacking own territory")
	}
}

func TestAttackNotEnoughTroops(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	from, to := findAttackPair(state, "player_0")
	if from == "" {
		t.Fatal("could not find attack pair")
	}

	// Set from to only 1 troop
	f := state.Territories[from]
	f.Troops = 1
	state.Territories[from] = f

	_, err := e.Attack(state, from, to, 1)
	if err == nil {
		t.Error("expected error for not enough troops")
	}
}

func TestAttackConquest(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	from, to := findAttackPair(state, "player_0")
	if from == "" {
		t.Fatal("could not find attack pair")
	}

	// Give attacker many troops, defender 1
	f := state.Territories[from]
	f.Troops = 50
	state.Territories[from] = f
	d := state.Territories[to]
	d.Troops = 1
	state.Territories[to] = d

	// Keep attacking until conquered (dice are random)
	conquered := false
	for i := 0; i < 20; i++ {
		// Resolve any pending conquest
		if state.PendingConquest != nil {
			e.MoveAfterConquest(state, 0)
		}
		result, err := e.Attack(state, from, to, 3)
		if err != nil {
			// Territory may already be conquered
			break
		}
		if result.Conquered {
			conquered = true
			break
		}
		// Restore defender troops to 1 for next attempt
		d = state.Territories[to]
		d.Troops = 1
		state.Territories[to] = d
	}

	if !conquered {
		t.Error("expected to conquer territory after multiple attacks")
		return
	}

	// After conquest, the territory should be owned by player_0
	if state.Territories[to].Owner != "player_0" {
		t.Errorf("expected conquered territory to be owned by player_0")
	}
	if !state.ConqueredThisTurn {
		t.Error("expected ConqueredThisTurn to be true")
	}
}

func TestAttackWrongPhase(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	_, err := e.Attack(state, "alaska", "kamchatka", 1)
	if err == nil {
		t.Error("expected error for attacking during place phase")
	}
}

func TestAttackWithPendingConquest(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	state.PendingConquest = &models.PendingConquest{
		From:      "alaska",
		To:        "kamchatka",
		MinTroops: 1,
		MaxTroops: 3,
	}

	from, to := findAttackPair(state, "player_0")
	if from == "" {
		t.Skip("no attack pair found")
	}
	f := state.Territories[from]
	f.Troops = 10
	state.Territories[from] = f

	_, err := e.Attack(state, from, to, 1)
	if err == nil {
		t.Error("expected error when pending conquest exists")
	}
}

// === TestMoveAfterConquest ===

func TestMoveAfterConquest(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)

	// Set up a pending conquest scenario
	from := findOwnedTerritory(state, "player_0")
	f := state.Territories[from]
	f.Troops = 5
	state.Territories[from] = f

	enemy := findEnemyNeighbor(state, from, "player_0")
	if enemy == "" {
		t.Fatal("no enemy neighbor found")
	}

	// Simulate conquest
	en := state.Territories[enemy]
	en.Owner = "player_0"
	en.Troops = 2
	state.Territories[enemy] = en

	state.PendingConquest = &models.PendingConquest{
		From:      from,
		To:        enemy,
		MinTroops: 2,
		MaxTroops: 3,
	}

	e := NewGameEngine()
	err := e.MoveAfterConquest(state, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.PendingConquest != nil {
		t.Error("expected pending conquest to be nil after resolve")
	}
	if state.Territories[from].Troops != 3 {
		t.Errorf("expected from troops 3, got %d", state.Territories[from].Troops)
	}
	if state.Territories[enemy].Troops != 4 {
		t.Errorf("expected to troops 4, got %d", state.Territories[enemy].Troops)
	}
}

func TestMoveAfterConquestNoPending(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	err := e.MoveAfterConquest(state, 1)
	if err == nil {
		t.Error("expected error when no pending conquest")
	}
}

func TestMoveAfterConquestTooManyTroops(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	state.PendingConquest = &models.PendingConquest{
		From:      "alaska",
		To:        "kamchatka",
		MinTroops: 1,
		MaxTroops: 2,
	}

	err := e.MoveAfterConquest(state, 5)
	if err == nil {
		t.Error("expected error for moving too many troops")
	}
}

// === TestFortify ===

func TestFortify(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	from := findOwnedTerritory(state, "player_0")
	friend := findOwnedNeighbor(state, from, "player_0")
	if friend == "" {
		t.Skip("no friendly neighbor found")
	}

	f := state.Territories[from]
	f.Troops = 5
	state.Territories[from] = f

	err := e.Fortify(state, from, friend, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Territories[from].Troops != 3 {
		t.Errorf("expected from troops 3, got %d", state.Territories[from].Troops)
	}
	if state.Territories[friend].Troops != 5 {
		t.Errorf("expected to troops 5, got %d", state.Territories[friend].Troops)
	}
}

func TestFortifyNotConnected(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	// Make only two territories owned by player_0, not connected
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		terr.Troops = 3
		state.Territories[id] = terr
	}

	// Give player_0 two non-adjacent territories
	t1 := state.Territories["alaska"]
	t1.Owner = "player_0"
	t1.Troops = 5
	state.Territories["alaska"] = t1

	t2 := state.Territories["argentina"]
	t2.Owner = "player_0"
	t2.Troops = 3
	state.Territories["argentina"] = t2

	err := e.Fortify(state, "alaska", "argentina", 2)
	if err == nil {
		t.Error("expected error for non-connected fortify")
	}
}

func TestFortifyNotEnoughTroops(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	from := findOwnedTerritory(state, "player_0")
	friend := findOwnedNeighbor(state, from, "player_0")
	if friend == "" {
		t.Skip("no friendly neighbor found")
	}

	f := state.Territories[from]
	f.Troops = 2
	state.Territories[from] = f

	err := e.Fortify(state, from, friend, 5)
	if err == nil {
		t.Error("expected error for not enough troops to fortify")
	}
}

func TestFortifyWrongOwner(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	enemyID := findOwnedTerritory(state, "player_1")
	ownedID := findOwnedTerritory(state, "player_0")

	err := e.Fortify(state, enemyID, ownedID, 1)
	if err == nil {
		t.Error("expected error for fortifying from enemy territory")
	}
}

func TestFortifyWrongPhase(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	err := e.Fortify(state, "alaska", "kamchatka", 1)
	if err == nil {
		t.Error("expected error for fortifying during attack phase")
	}
}

// === TestEndPhase ===

func TestEndPhasePlaceToAttack(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	state.TroopsToDeploy = 0
	e := NewGameEngine()

	err := e.EndPhase(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Phase != models.PhaseAttack {
		t.Errorf("expected phase attack, got %s", state.Phase)
	}
}

func TestEndPhasePlaceWithTroopsRemaining(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	state.TroopsToDeploy = 3
	e := NewGameEngine()

	err := e.EndPhase(state)
	if err == nil {
		t.Error("expected error when troops remain")
	}
}

func TestEndPhaseAttackToFortify(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	err := e.EndPhase(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Phase != models.PhaseFortify {
		t.Errorf("expected phase fortify, got %s", state.Phase)
	}
}

func TestEndPhaseFortifyToNextPlayerPlace(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	origPlayer := state.CurrentPlayer
	err := e.EndPhase(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Phase != models.PhasePlace {
		t.Errorf("expected phase place, got %s", state.Phase)
	}
	if state.CurrentPlayer == origPlayer {
		t.Error("expected current player to advance")
	}
	if state.TroopsToDeploy < 3 {
		t.Error("expected at least 3 troops to deploy for next player")
	}
}

func TestEndPhaseCardAwardOnConquest(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	state.ConqueredThisTurn = true
	deckSize := len(state.Deck)
	e := NewGameEngine()

	// Count cards before
	var cardsBefore int
	for _, p := range state.Players {
		if p.ID == state.CurrentPlayer {
			cardsBefore = len(p.Cards)
		}
	}

	err := e.EndPhase(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The card was awarded to player_0 (the previous current player, before advancing)
	// After advancing, current player is player_1, so check player_0's cards
	var cardsAfter int
	for _, p := range state.Players {
		if p.ID == "player_0" {
			cardsAfter = len(p.Cards)
		}
	}

	if cardsAfter != cardsBefore+1 {
		t.Errorf("expected card count %d, got %d", cardsBefore+1, cardsAfter)
	}
	if len(state.Deck) != deckSize-1 {
		t.Errorf("expected deck size %d, got %d", deckSize-1, len(state.Deck))
	}
}

func TestEndPhaseEnded(t *testing.T) {
	state := newTestGameInPhase(models.PhaseEnded)
	e := NewGameEngine()

	err := e.EndPhase(state)
	if err == nil {
		t.Error("expected error for ending phase during 'ended'")
	}
}

// === TestTradeCards ===

func TestTradeCardsValidSet(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	// Give player cards: one of each type
	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "cavalry"},
				{Territory: "china", Type: "artillery"},
			}
		}
	}

	origDeploy := state.TroopsToDeploy
	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have received bonus troops
	if state.TroopsToDeploy <= origDeploy {
		t.Error("expected troops to deploy to increase after trade")
	}

	// Cards should be removed
	for _, p := range state.Players {
		if p.ID == "player_0" && len(p.Cards) != 0 {
			t.Errorf("expected 0 cards after trade, got %d", len(p.Cards))
		}
	}
}

func TestTradeCardsInvalidSet(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "infantry"},
				{Territory: "china", Type: "cavalry"},
			}
		}
	}

	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err == nil {
		t.Error("expected error for invalid card set (2 infantry + 1 cavalry)")
	}
}

func TestTradeCardsWrongPhase(t *testing.T) {
	state := newTestGameInPhase(models.PhaseFortify)
	e := NewGameEngine()

	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "cavalry"},
				{Territory: "china", Type: "artillery"},
			}
		}
	}

	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err == nil {
		t.Error("expected error for trading cards during fortify phase")
	}
}

func TestTradeCardsWithWild(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "", Type: "wild"},
				{Territory: "brazil", Type: "infantry"},
				{Territory: "china", Type: "infantry"},
			}
		}
	}

	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err != nil {
		t.Fatalf("wild card set should be valid: %v", err)
	}
}

func TestTradeCardsForcedDuringAttackWith6Cards(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "cavalry"},
				{Territory: "china", Type: "artillery"},
				{Territory: "india", Type: "infantry"},
				{Territory: "japan", Type: "cavalry"},
				{Territory: "peru", Type: "artillery"},
			}
		}
	}

	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err != nil {
		t.Fatalf("should allow trade with 6+ cards during attack: %v", err)
	}
}

func TestTradeCardsThreeOfSameType(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)
	e := NewGameEngine()

	for i := range state.Players {
		if state.Players[i].ID == "player_0" {
			state.Players[i].Cards = []models.Card{
				{Territory: "alaska", Type: "infantry"},
				{Territory: "brazil", Type: "infantry"},
				{Territory: "china", Type: "infantry"},
			}
		}
	}

	err := e.TradeCards(state, [3]int{0, 1, 2})
	if err != nil {
		t.Fatalf("three of same type should be valid: %v", err)
	}
}

// === TestCalculateTroops ===

func TestCalculateTroopsMinimum(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Give player_0 only 3 territories (would be 3/3=1, but minimum is 3)
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		state.Territories[id] = terr
	}
	count := 0
	for id, terr := range state.Territories {
		if count < 3 {
			terr.Owner = "player_0"
			state.Territories[id] = terr
			count++
		}
	}

	troops := CalculateTroops(state, "player_0")
	if troops < 3 {
		t.Errorf("expected minimum 3 troops, got %d", troops)
	}
}

func TestCalculateTroopsTerritoryBonus(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Give player_0 exactly 12 territories (12/3=4)
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		state.Territories[id] = terr
	}
	count := 0
	for id, terr := range state.Territories {
		if count < 12 {
			terr.Owner = "player_0"
			state.Territories[id] = terr
			count++
		}
	}

	troops := CalculateTroops(state, "player_0")
	if troops < 4 {
		t.Errorf("expected at least 4 troops for 12 territories, got %d", troops)
	}
}

func TestCalculateTroopsContinentBonus(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Give player_0 all of Australia (4 territories) + some others
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		state.Territories[id] = terr
	}

	// Give Australia to player_0
	for _, tid := range models.ContinentTerritories["australia"] {
		terr := state.Territories[tid]
		terr.Owner = "player_0"
		state.Territories[tid] = terr
	}
	// Add some more territories to get above minimum
	extraCount := 0
	for id, terr := range state.Territories {
		if terr.Owner != "player_0" && extraCount < 5 {
			terr.Owner = "player_0"
			state.Territories[id] = terr
			extraCount++
		}
	}

	troops := CalculateTroops(state, "player_0")
	// 9 territories / 3 = 3, + 2 for Australia = 5
	if troops < 5 {
		t.Errorf("expected at least 5 troops with Australia bonus, got %d", troops)
	}
}

// === TestPlayerElimination ===

func TestPlayerElimination(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	// Give player_1 only one territory, player_0 gets the rest
	var lastEnemy string
	for id, terr := range state.Territories {
		terr.Owner = "player_0"
		terr.Troops = 5
		state.Territories[id] = terr
	}

	// Find an adjacent pair for the attack
	fromID := "alaska"
	// Give kamchatka to player_1
	lastEnemy = "kamchatka"
	en := state.Territories[lastEnemy]
	en.Owner = "player_1"
	en.Troops = 1
	state.Territories[lastEnemy] = en

	f := state.Territories[fromID]
	f.Troops = 50
	state.Territories[fromID] = f

	// Give player_1 some cards to test transfer
	for i := range state.Players {
		if state.Players[i].ID == "player_1" {
			state.Players[i].Cards = []models.Card{
				{Territory: "brazil", Type: "infantry"},
				{Territory: "china", Type: "cavalry"},
			}
		}
	}

	// Attack until conquered
	for i := 0; i < 20; i++ {
		if state.PendingConquest != nil {
			e.MoveAfterConquest(state, 0)
		}
		result, err := e.Attack(state, fromID, lastEnemy, 3)
		if err != nil {
			break
		}
		if result.Conquered {
			break
		}
		// Reset defender
		d := state.Territories[lastEnemy]
		d.Troops = 1
		d.Owner = "player_1"
		state.Territories[lastEnemy] = d
	}

	// Check player_1 is eliminated
	for _, p := range state.Players {
		if p.ID == "player_1" && p.IsAlive {
			t.Error("expected player_1 to be eliminated")
		}
	}

	// Check cards transferred to player_0
	var p0Cards int
	for _, p := range state.Players {
		if p.ID == "player_0" {
			p0Cards = len(p.Cards)
		}
	}
	if p0Cards < 2 {
		t.Errorf("expected player_0 to have at least 2 cards from elimination, got %d", p0Cards)
	}
}

// === TestWinCondition ===

func TestWinCondition(t *testing.T) {
	state := newTestGameInPhase(models.PhaseAttack)
	e := NewGameEngine()

	// Give all territories to player_0 except one
	for id, terr := range state.Territories {
		terr.Owner = "player_0"
		terr.Troops = 5
		state.Territories[id] = terr
	}

	// Give kamchatka to player_1 with 1 troop
	en := state.Territories["kamchatka"]
	en.Owner = "player_1"
	en.Troops = 1
	state.Territories["kamchatka"] = en

	f := state.Territories["alaska"]
	f.Troops = 50
	state.Territories["alaska"] = f

	// Attack until conquered
	for i := 0; i < 20; i++ {
		if state.PendingConquest != nil {
			e.MoveAfterConquest(state, 0)
		}
		result, err := e.Attack(state, "alaska", "kamchatka", 3)
		if err != nil {
			break
		}
		if result.Conquered {
			break
		}
		// Reset defender
		d := state.Territories["kamchatka"]
		d.Troops = 1
		d.Owner = "player_1"
		state.Territories["kamchatka"] = d
	}

	if state.Phase != models.PhaseEnded {
		t.Error("expected game to end after winning")
	}
	if state.Winner != "player_0" {
		t.Errorf("expected winner player_0, got %s", state.Winner)
	}
}

// === TestIsConnected ===

func TestIsConnectedDirectNeighbors(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Make alaska and northwest_territory both owned by player_0
	a := state.Territories["alaska"]
	a.Owner = "player_0"
	state.Territories["alaska"] = a

	b := state.Territories["northwest_territory"]
	b.Owner = "player_0"
	state.Territories["northwest_territory"] = b

	if !IsConnected(state, "alaska", "northwest_territory", "player_0") {
		t.Error("expected adjacent territories to be connected")
	}
}

func TestIsConnectedSameTerritory(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	if !IsConnected(state, "alaska", "alaska", "player_0") {
		t.Error("expected territory to be connected to itself")
	}
}

func TestIsConnectedThroughChain(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Give player_0 a chain: alaska -> northwest_territory -> alberta
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		state.Territories[id] = terr
	}

	for _, tid := range []string{"alaska", "northwest_territory", "alberta"} {
		terr := state.Territories[tid]
		terr.Owner = "player_0"
		state.Territories[tid] = terr
	}

	if !IsConnected(state, "alaska", "alberta", "player_0") {
		t.Error("expected chain-connected territories to be connected")
	}
}

func TestIsConnectedNotConnected(t *testing.T) {
	state := newTestGameInPhase(models.PhasePlace)

	// Give player_0 only alaska and argentina (not connected)
	for id, terr := range state.Territories {
		terr.Owner = "player_1"
		state.Territories[id] = terr
	}

	a := state.Territories["alaska"]
	a.Owner = "player_0"
	state.Territories["alaska"] = a

	b := state.Territories["argentina"]
	b.Owner = "player_0"
	state.Territories["argentina"] = b

	if IsConnected(state, "alaska", "argentina", "player_0") {
		t.Error("expected disconnected territories to not be connected")
	}
}
