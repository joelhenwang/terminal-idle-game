package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SaveData represents the serializable game state
type SaveData struct {
	Score       int
	NumUpgrades int
	Summator    int
	Generators  []GeneratorData
}

// GeneratorData represents the serializable generator state
type GeneratorData struct {
	Name       string
	Level      int
	BaseCost   int
	BaseOutput int
}

// ToSaveData converts the game state to a saveable format
func (g *Game) ToSaveData() *SaveData {
	generators := make([]GeneratorData, len(g.Generators))
	for i, gen := range g.Generators {
		generators[i] = GeneratorData{
			Name:       gen.Name,
			Level:      gen.Level,
			BaseCost:   gen.BaseCost,
			BaseOutput: gen.BaseOutput,
		}
	}
	return &SaveData{
		Score:       g.Score,
		NumUpgrades: g.NumUpgrades,
		Summator:    g.Summator,
		Generators:  generators,
	}
}

// LoadFromSaveData loads the game state from save data
func (g *Game) LoadFromSaveData(data *SaveData) {
	g.Score = data.Score
	g.NumUpgrades = data.NumUpgrades
	g.Summator = data.Summator
	
	// Recreate generators with saved state
	g.Generators = make([]*Generator, len(data.Generators))
	for i, genData := range data.Generators {
		g.Generators[i] = &Generator{
			Name:       genData.Name,
			Level:      genData.Level,
			BaseCost:   genData.BaseCost,
			BaseOutput: genData.BaseOutput,
			CostMult:   1.5,
		}
	}
}

// Save saves the game state to a file
func (g *Game) Save(filename string) error {
	saveData := g.ToSaveData()
	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Load loads the game state from a file
func (g *Game) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return err
	}

	g.LoadFromSaveData(&saveData)
	return nil
}

// GetDefaultSaveFile returns the default save file path
func GetDefaultSaveFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "savegame.json"
	}
	return filepath.Join(homeDir, ".terminal-idle-game", "savegame.json")
}
