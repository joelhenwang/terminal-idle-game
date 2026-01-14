package game

// State represents the current state of the game
type State int

const (
	Running State = iota
	Paused
)

// Game represents the core game state and logic
type Game struct {
	Score       int
	NumUpgrades int
	Summator    int
	State       State
	Generators  []*Generator
}

// New creates a new game with default values
func New() *Game {
	return &Game{
		Score:       0,
		NumUpgrades: 0,
		Summator:    2,
		State:       Running,
		Generators: []*Generator{
			NewGenerator("Clicker", 15, 1),
			NewGenerator("Farm", 100, 5),
			NewGenerator("Mine", 500, 20),
			NewGenerator("Factory", 2000, 50),
		},
	}
}

// Tick updates the game state for one time unit
func (g *Game) Tick() {
	if g.State == Running {
		g.Score += g.Summator
		// Add production from generators
		for _, gen := range g.Generators {
			g.Score += gen.Output()
		}
	}
}

// TogglePause toggles between Running and Paused states
func (g *Game) TogglePause() {
	if g.State == Running {
		g.State = Paused
	} else if g.State == Paused {
		g.State = Running
	}
}

// CanUpgrade returns true if the player has enough score to upgrade
func (g *Game) CanUpgrade() bool {
	return g.Score >= g.UpgradeCost()
}

// UpgradeCost returns the cost of the next upgrade
func (g *Game) UpgradeCost() int {
	return 10 * g.Summator
}

// Upgrade performs an upgrade if possible and returns true if successful
func (g *Game) Upgrade() bool {
	if !g.CanUpgrade() {
		return false
	}
	g.Score -= g.UpgradeCost()
	g.NumUpgrades++
	g.Summator *= 2
	return true
}

// IsRunning returns true if the game is in Running state
func (g *Game) IsRunning() bool {
	return g.State == Running
}

// IsPaused returns true if the game is in Paused state
func (g *Game) IsPaused() bool {
	return g.State == Paused
}

// BuyGenerator attempts to buy a generator upgrade and returns true if successful
func (g *Game) BuyGenerator(index int) bool {
	if index < 0 || index >= len(g.Generators) {
		return false
	}
	gen := g.Generators[index]
	cost := gen.Cost()
	if g.Score >= cost {
		g.Score -= cost
		gen.Upgrade()
		return true
	}
	return false
}

// TotalProduction returns the total production rate per tick
func (g *Game) TotalProduction() int {
	total := g.Summator
	for _, gen := range g.Generators {
		total += gen.Output()
	}
	return total
}

