package game

import "testing"

func TestNew(t *testing.T) {
	g := New()
	if g.Score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", g.Score)
	}
	if g.NumUpgrades != 0 {
		t.Errorf("Expected initial numUpgrades to be 0, got %d", g.NumUpgrades)
	}
	if g.Summator != 2 {
		t.Errorf("Expected initial summator to be 2, got %d", g.Summator)
	}
	if g.State != Running {
		t.Errorf("Expected initial state to be Running, got %v", g.State)
	}
}

func TestTick(t *testing.T) {
	tests := []struct {
		name           string
		initialState   State
		initialScore   int
		summator       int
		expectedScore  int
	}{
		{
			name:          "Running state increases score",
			initialState:  Running,
			initialScore:  0,
			summator:      2,
			expectedScore: 2,
		},
		{
			name:          "Paused state does not increase score",
			initialState:  Paused,
			initialScore:  0,
			summator:      2,
			expectedScore: 0,
		},
		{
			name:          "Multiple ticks accumulate score",
			initialState:  Running,
			initialScore:  10,
			summator:      5,
			expectedScore: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Score:    tt.initialScore,
				Summator: tt.summator,
				State:    tt.initialState,
			}
			g.Tick()
			if g.Score != tt.expectedScore {
				t.Errorf("Expected score %d, got %d", tt.expectedScore, g.Score)
			}
		})
	}
}

func TestTogglePause(t *testing.T) {
	tests := []struct {
		name          string
		initialState  State
		expectedState State
	}{
		{
			name:          "Toggle from Running to Paused",
			initialState:  Running,
			expectedState: Paused,
		},
		{
			name:          "Toggle from Paused to Running",
			initialState:  Paused,
			expectedState: Running,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{State: tt.initialState}
			g.TogglePause()
			if g.State != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, g.State)
			}
		})
	}
}

func TestCanUpgrade(t *testing.T) {
	tests := []struct {
		name       string
		score      int
		summator   int
		canUpgrade bool
	}{
		{
			name:       "Insufficient score",
			score:      10,
			summator:   2,
			canUpgrade: false,
		},
		{
			name:       "Exact score needed",
			score:      20,
			summator:   2,
			canUpgrade: true,
		},
		{
			name:       "More than enough score",
			score:      100,
			summator:   2,
			canUpgrade: true,
		},
		{
			name:       "Zero score",
			score:      0,
			summator:   2,
			canUpgrade: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Score:    tt.score,
				Summator: tt.summator,
			}
			if got := g.CanUpgrade(); got != tt.canUpgrade {
				t.Errorf("Expected CanUpgrade() = %v, got %v", tt.canUpgrade, got)
			}
		})
	}
}

func TestUpgradeCost(t *testing.T) {
	tests := []struct {
		name         string
		summator     int
		expectedCost int
	}{
		{
			name:         "Initial summator",
			summator:     2,
			expectedCost: 20,
		},
		{
			name:         "After one upgrade",
			summator:     4,
			expectedCost: 40,
		},
		{
			name:         "After multiple upgrades",
			summator:     16,
			expectedCost: 160,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{Summator: tt.summator}
			if got := g.UpgradeCost(); got != tt.expectedCost {
				t.Errorf("Expected UpgradeCost() = %d, got %d", tt.expectedCost, got)
			}
		})
	}
}

func TestUpgrade(t *testing.T) {
	tests := []struct {
		name              string
		initialScore      int
		initialSummator   int
		initialUpgrades   int
		expectedSuccess   bool
		expectedScore     int
		expectedSummator  int
		expectedUpgrades  int
	}{
		{
			name:             "Successful upgrade",
			initialScore:     20,
			initialSummator:  2,
			initialUpgrades:  0,
			expectedSuccess:  true,
			expectedScore:    0,
			expectedSummator: 4,
			expectedUpgrades: 1,
		},
		{
			name:             "Failed upgrade - insufficient score",
			initialScore:     10,
			initialSummator:  2,
			initialUpgrades:  0,
			expectedSuccess:  false,
			expectedScore:    10,
			expectedSummator: 2,
			expectedUpgrades: 0,
		},
		{
			name:             "Successful upgrade with excess score",
			initialScore:     100,
			initialSummator:  2,
			initialUpgrades:  0,
			expectedSuccess:  true,
			expectedScore:    80,
			expectedSummator: 4,
			expectedUpgrades: 1,
		},
		{
			name:             "Multiple upgrades progression",
			initialScore:     40,
			initialSummator:  4,
			initialUpgrades:  1,
			expectedSuccess:  true,
			expectedScore:    0,
			expectedSummator: 8,
			expectedUpgrades: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Score:       tt.initialScore,
				Summator:    tt.initialSummator,
				NumUpgrades: tt.initialUpgrades,
			}
			success := g.Upgrade()
			if success != tt.expectedSuccess {
				t.Errorf("Expected Upgrade() = %v, got %v", tt.expectedSuccess, success)
			}
			if g.Score != tt.expectedScore {
				t.Errorf("Expected score %d, got %d", tt.expectedScore, g.Score)
			}
			if g.Summator != tt.expectedSummator {
				t.Errorf("Expected summator %d, got %d", tt.expectedSummator, g.Summator)
			}
			if g.NumUpgrades != tt.expectedUpgrades {
				t.Errorf("Expected numUpgrades %d, got %d", tt.expectedUpgrades, g.NumUpgrades)
			}
		})
	}
}

func TestIsRunning(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected bool
	}{
		{
			name:     "Running state",
			state:    Running,
			expected: true,
		},
		{
			name:     "Paused state",
			state:    Paused,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{State: tt.state}
			if got := g.IsRunning(); got != tt.expected {
				t.Errorf("Expected IsRunning() = %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestIsPaused(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected bool
	}{
		{
			name:     "Running state",
			state:    Running,
			expected: false,
		},
		{
			name:     "Paused state",
			state:    Paused,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{State: tt.state}
			if got := g.IsPaused(); got != tt.expected {
				t.Errorf("Expected IsPaused() = %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestBuyGenerator(t *testing.T) {
	tests := []struct {
		name            string
		initialScore    int
		generatorIndex  int
		expectedSuccess bool
		expectedScore   int
		expectedLevel   int
	}{
		{
			name:            "Successful purchase",
			initialScore:    100,
			generatorIndex:  0,
			expectedSuccess: true,
			expectedScore:   85, // 100 - 15
			expectedLevel:   1,
		},
		{
			name:            "Insufficient funds",
			initialScore:    10,
			generatorIndex:  0,
			expectedSuccess: false,
			expectedScore:   10,
			expectedLevel:   0,
		},
		{
			name:            "Invalid index negative",
			initialScore:    100,
			generatorIndex:  -1,
			expectedSuccess: false,
			expectedScore:   100,
			expectedLevel:   0,
		},
		{
			name:            "Invalid index too high",
			initialScore:    100,
			generatorIndex:  10,
			expectedSuccess: false,
			expectedScore:   100,
			expectedLevel:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			g.Score = tt.initialScore
			success := g.BuyGenerator(tt.generatorIndex)
			if success != tt.expectedSuccess {
				t.Errorf("Expected BuyGenerator() = %v, got %v", tt.expectedSuccess, success)
			}
			if g.Score != tt.expectedScore {
				t.Errorf("Expected score %d, got %d", tt.expectedScore, g.Score)
			}
			if tt.generatorIndex >= 0 && tt.generatorIndex < len(g.Generators) {
				if g.Generators[tt.generatorIndex].Level != tt.expectedLevel {
					t.Errorf("Expected generator level %d, got %d", tt.expectedLevel, g.Generators[tt.generatorIndex].Level)
				}
			}
		})
	}
}

func TestTotalProduction(t *testing.T) {
	g := New()
	
	// Initially only summator produces
	expected := g.Summator
	if got := g.TotalProduction(); got != expected {
		t.Errorf("Expected total production %d, got %d", expected, got)
	}

	// Buy first generator
	g.Score = 100
	g.BuyGenerator(0)
	expected = g.Summator + g.Generators[0].Output()
	if got := g.TotalProduction(); got != expected {
		t.Errorf("Expected total production %d, got %d", expected, got)
	}

	// Buy second generator
	g.Score = 200
	g.BuyGenerator(1)
	expected = g.Summator + g.Generators[0].Output() + g.Generators[1].Output()
	if got := g.TotalProduction(); got != expected {
		t.Errorf("Expected total production %d, got %d", expected, got)
	}
}

func TestTickWithGenerators(t *testing.T) {
	g := New()
	g.Score = 100
	g.BuyGenerator(0) // Buy clicker (produces 1/s)
	
	initialScore := g.Score
	g.Tick()
	
	expectedScore := initialScore + g.Summator + 1
	if g.Score != expectedScore {
		t.Errorf("Expected score %d after tick with generator, got %d", expectedScore, g.Score)
	}
}

