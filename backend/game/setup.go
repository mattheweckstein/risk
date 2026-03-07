package game

import (
	"fmt"

	"github.com/mattheweckstein/risk/backend/models"
)

// SetupPlaceTroop handles placing a single troop during the setup phase.
// During setup, players take turns placing 1 troop at a time on territories they own.
// When all players have placed all starting troops, the game transitions to the first player's place phase.
func (e *GameEngine) SetupPlaceTroop(state *models.GameState, territory string) error {
	if state.Phase != models.PhaseSetup {
		return fmt.Errorf("not in setup phase")
	}

	return e.PlaceTroops(state, territory, 1)
}

// IsSetupComplete returns true if all players have placed all their starting troops.
func (e *GameEngine) IsSetupComplete(state *models.GameState) bool {
	for _, p := range state.Players {
		if !p.IsAlive {
			continue
		}
		if e.setupTroopsRemaining(state, p.ID) > 0 {
			return false
		}
	}
	return true
}

// GetSetupTroopsRemaining returns how many troops a player still needs to place during setup.
func (e *GameEngine) GetSetupTroopsRemaining(state *models.GameState, playerID string) int {
	return e.setupTroopsRemaining(state, playerID)
}

// GetCurrentSetupPlayer returns the ID of the player who should place troops next in setup.
func (e *GameEngine) GetCurrentSetupPlayer(state *models.GameState) string {
	return state.CurrentPlayer
}

// IsCurrentPlayerAI returns true if the current player is an AI player.
func (e *GameEngine) IsCurrentPlayerAI(state *models.GameState) bool {
	for _, p := range state.Players {
		if p.ID == state.CurrentPlayer {
			return p.IsAI
		}
	}
	return false
}
