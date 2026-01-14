package game

// Generator represents a type of generator that produces resources
type Generator struct {
	Name       string
	Level      int
	BaseCost   int
	BaseOutput int
	CostMult   float64 // Cost multiplier per level
}

// NewGenerator creates a new generator
func NewGenerator(name string, baseCost, baseOutput int) *Generator {
	return &Generator{
		Name:       name,
		Level:      0,
		BaseCost:   baseCost,
		BaseOutput: baseOutput,
		CostMult:   1.5,
	}
}

// Cost returns the cost to upgrade this generator to the next level
func (g *Generator) Cost() int {
	if g.Level == 0 {
		return g.BaseCost
	}
	cost := float64(g.BaseCost)
	for i := 0; i < g.Level; i++ {
		cost *= g.CostMult
	}
	return int(cost)
}

// Output returns the current output per tick of this generator
func (g *Generator) Output() int {
	if g.Level == 0 {
		return 0
	}
	return g.BaseOutput * g.Level
}

// Upgrade increases the generator level by 1
func (g *Generator) Upgrade() {
	g.Level++
}
