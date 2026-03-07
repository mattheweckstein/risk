package models

// Phase represents the current phase of the game
type Phase string

const (
	PhaseSetup   Phase = "setup"
	PhasePlace   Phase = "place"
	PhaseAttack  Phase = "attack"
	PhaseFortify Phase = "fortify"
	PhaseEnded   Phase = "ended"
)

// Territory represents a single territory on the map
type Territory struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Continent string   `json:"continent"`
	Neighbors []string `json:"neighbors"`
	Owner     string   `json:"owner"`
	Troops    int      `json:"troops"`
}

// Player represents a player in the game
type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	IsAI  bool   `json:"isAI"`
	Color string `json:"color"`
	Cards []Card `json:"cards"`
	IsAlive bool `json:"isAlive"`
}

// GameState represents the complete state of a game
type GameState struct {
	ID             string               `json:"id"`
	Phase          Phase                `json:"phase"`
	Turn           int                  `json:"turn"`
	CurrentPlayer  string               `json:"currentPlayer"`
	Players        []Player             `json:"players"`
	Territories    map[string]Territory `json:"territories"`
	Deck           []Card               `json:"deck"`
	Log            []LogEntry           `json:"log"`
	Winner         string               `json:"winner,omitempty"`
	TroopsToDeploy int                  `json:"troopsToDeploy"`
	CardTradeCount int                  `json:"cardTradeCount"`
	ConqueredThisTurn bool             `json:"conqueredThisTurn"`
	LastAttackResult *AttackResult      `json:"lastAttackResult,omitempty"`
}

// Card represents a Risk card
type Card struct {
	Territory string `json:"territory"`
	Type      string `json:"type"` // infantry, cavalry, artillery, wild
}

// LogEntry represents a single entry in the game log
type LogEntry struct {
	Turn    int    `json:"turn"`
	Player  string `json:"player"`
	Message string `json:"message"`
}

// AttackResult holds the outcome of an attack
type AttackResult struct {
	AttackerRolls   []int  `json:"attackerRolls"`
	DefenderRolls   []int  `json:"defenderRolls"`
	AttackerLosses  int    `json:"attackerLosses"`
	DefenderLosses  int    `json:"defenderLosses"`
	Conquered       bool   `json:"conquered"`
	AttackingTerritory string `json:"attackingTerritory"`
	DefendingTerritory string `json:"defendingTerritory"`
}

// NewGameRequest is the request body for creating a new game
type NewGameRequest struct {
	PlayerName string   `json:"playerName"`
	AICount    int      `json:"aiCount"`
	AINames    []string `json:"aiNames,omitempty"`
}

// PlaceRequest is the request body for placing troops
type PlaceRequest struct {
	Territory string `json:"territory"`
	Troops    int    `json:"troops"`
}

// AttackRequest is the request body for attacking
type AttackRequest struct {
	From         string `json:"from"`
	To           string `json:"to"`
	AttackerDice int    `json:"attackerDice"`
}

// FortifyRequest is the request body for fortifying
type FortifyRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Troops int    `json:"troops"`
}

// CardTradeRequest is the request body for trading cards
type CardTradeRequest struct {
	CardIndices [3]int `json:"cardIndices"`
}

// ContinentBonus defines the troop bonus for controlling a continent
var ContinentBonuses = map[string]int{
	"north_america": 5,
	"south_america": 2,
	"europe":        5,
	"africa":        3,
	"asia":          7,
	"australia":     2,
}

// ContinentTerritories maps each continent to its territory IDs
var ContinentTerritories = map[string][]string{
	"north_america": {
		"alaska", "northwest_territory", "greenland", "alberta",
		"ontario", "quebec", "western_us", "eastern_us", "central_america",
	},
	"south_america": {
		"venezuela", "peru", "brazil", "argentina",
	},
	"europe": {
		"iceland", "scandinavia", "great_britain", "northern_europe",
		"western_europe", "southern_europe", "ukraine",
	},
	"africa": {
		"north_africa", "egypt", "east_africa", "congo",
		"south_africa", "madagascar",
	},
	"asia": {
		"ural", "siberia", "yakutsk", "kamchatka", "irkutsk",
		"mongolia", "japan", "afghanistan", "china", "india",
		"siam", "middle_east",
	},
	"australia": {
		"indonesia", "new_guinea", "western_australia", "eastern_australia",
	},
}
