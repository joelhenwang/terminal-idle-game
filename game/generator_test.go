package game

import "testing"

func TestNewGenerator(t *testing.T) {
	g := NewGenerator("Clicker", 10, 1)
	if g.Name != "Clicker" {
		t.Errorf("Expected name 'Clicker', got %s", g.Name)
	}
	if g.Level != 0 {
		t.Errorf("Expected initial level 0, got %d", g.Level)
	}
	if g.BaseCost != 10 {
		t.Errorf("Expected base cost 10, got %d", g.BaseCost)
	}
	if g.BaseOutput != 1 {
		t.Errorf("Expected base output 1, got %d", g.BaseOutput)
	}
	if g.CostMult != 1.5 {
		t.Errorf("Expected cost multiplier 1.5, got %f", g.CostMult)
	}
}

func TestGeneratorCost(t *testing.T) {
	tests := []struct {
		name         string
		baseCost     int
		level        int
		expectedCost int
	}{
		{
			name:         "Level 0 costs base cost",
			baseCost:     10,
			level:        0,
			expectedCost: 10,
		},
		{
			name:         "Level 1 costs base * 1.5",
			baseCost:     10,
			level:        1,
			expectedCost: 15,
		},
		{
			name:         "Level 2 costs base * 1.5^2",
			baseCost:     10,
			level:        2,
			expectedCost: 22,
		},
		{
			name:         "Level 3 costs base * 1.5^3",
			baseCost:     10,
			level:        3,
			expectedCost: 33,
		},
		{
			name:         "Higher base cost",
			baseCost:     100,
			level:        1,
			expectedCost: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				BaseCost: tt.baseCost,
				Level:    tt.level,
				CostMult: 1.5,
			}
			cost := g.Cost()
			if cost != tt.expectedCost {
				t.Errorf("Expected cost %d, got %d", tt.expectedCost, cost)
			}
		})
	}
}

func TestGeneratorOutput(t *testing.T) {
	tests := []struct {
		name           string
		baseOutput     int
		level          int
		expectedOutput int
	}{
		{
			name:           "Level 0 produces nothing",
			baseOutput:     1,
			level:          0,
			expectedOutput: 0,
		},
		{
			name:           "Level 1 produces base output",
			baseOutput:     1,
			level:          1,
			expectedOutput: 1,
		},
		{
			name:           "Level 2 produces 2x base output",
			baseOutput:     1,
			level:          2,
			expectedOutput: 2,
		},
		{
			name:           "Level 5 with base 3",
			baseOutput:     3,
			level:          5,
			expectedOutput: 15,
		},
		{
			name:           "Higher level scaling",
			baseOutput:     10,
			level:          10,
			expectedOutput: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				BaseOutput: tt.baseOutput,
				Level:      tt.level,
			}
			output := g.Output()
			if output != tt.expectedOutput {
				t.Errorf("Expected output %d, got %d", tt.expectedOutput, output)
			}
		})
	}
}

func TestGeneratorUpgrade(t *testing.T) {
	g := NewGenerator("Test", 10, 1)
	initialLevel := g.Level

	g.Upgrade()
	if g.Level != initialLevel+1 {
		t.Errorf("Expected level %d, got %d", initialLevel+1, g.Level)
	}

	g.Upgrade()
	g.Upgrade()
	if g.Level != initialLevel+3 {
		t.Errorf("Expected level %d after 3 upgrades, got %d", initialLevel+3, g.Level)
	}
}

func TestGeneratorIntegration(t *testing.T) {
	g := NewGenerator("Clicker", 10, 1)

	// Initially produces nothing
	if g.Output() != 0 {
		t.Errorf("Expected initial output 0, got %d", g.Output())
	}

	// Cost should be base cost
	if g.Cost() != 10 {
		t.Errorf("Expected initial cost 10, got %d", g.Cost())
	}

	// After first upgrade
	g.Upgrade()
	if g.Output() != 1 {
		t.Errorf("Expected output 1 after first upgrade, got %d", g.Output())
	}
	if g.Cost() != 15 {
		t.Errorf("Expected cost 15 after first upgrade, got %d", g.Cost())
	}

	// After second upgrade
	g.Upgrade()
	if g.Output() != 2 {
		t.Errorf("Expected output 2 after second upgrade, got %d", g.Output())
	}
	if g.Cost() != 22 {
		t.Errorf("Expected cost 22 after second upgrade, got %d", g.Cost())
	}
}
