package game

import (
	"os"
	"path/filepath"
	"testing"
)

func TestToSaveData(t *testing.T) {
	g := New()
	g.Score = 100
	g.NumUpgrades = 2
	g.Summator = 8

	saveData := g.ToSaveData()

	if saveData.Score != 100 {
		t.Errorf("Expected score 100, got %d", saveData.Score)
	}
	if saveData.NumUpgrades != 2 {
		t.Errorf("Expected numUpgrades 2, got %d", saveData.NumUpgrades)
	}
	if saveData.Summator != 8 {
		t.Errorf("Expected summator 8, got %d", saveData.Summator)
	}
	if len(saveData.Generators) != len(g.Generators) {
		t.Errorf("Expected %d generators, got %d", len(g.Generators), len(saveData.Generators))
	}
}

func TestLoadFromSaveData(t *testing.T) {
	saveData := &SaveData{
		Score:       200,
		NumUpgrades: 3,
		Summator:    16,
		Generators: []GeneratorData{
			{Name: "Test1", Level: 2, BaseCost: 10, BaseOutput: 1},
			{Name: "Test2", Level: 1, BaseCost: 50, BaseOutput: 5},
		},
	}

	g := New()
	g.LoadFromSaveData(saveData)

	if g.Score != 200 {
		t.Errorf("Expected score 200, got %d", g.Score)
	}
	if g.NumUpgrades != 3 {
		t.Errorf("Expected numUpgrades 3, got %d", g.NumUpgrades)
	}
	if g.Summator != 16 {
		t.Errorf("Expected summator 16, got %d", g.Summator)
	}
	if len(g.Generators) != 2 {
		t.Errorf("Expected 2 generators, got %d", len(g.Generators))
	}
	if g.Generators[0].Level != 2 {
		t.Errorf("Expected first generator level 2, got %d", g.Generators[0].Level)
	}
	if g.Generators[1].Level != 1 {
		t.Errorf("Expected second generator level 1, got %d", g.Generators[1].Level)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "test_save.json")

	// Create a game with some progress
	g1 := New()
	g1.Score = 500
	g1.NumUpgrades = 3
	g1.Summator = 16
	g1.BuyGenerator(0)
	g1.BuyGenerator(0)
	g1.BuyGenerator(1)

	// Save the game
	err := g1.Save(saveFile)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		t.Fatal("Save file was not created")
	}

	// Create a new game and load the saved state
	g2 := New()
	err = g2.Load(saveFile)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Verify the loaded state matches the saved state
	if g2.Score != g1.Score {
		t.Errorf("Expected score %d, got %d", g1.Score, g2.Score)
	}
	if g2.NumUpgrades != g1.NumUpgrades {
		t.Errorf("Expected numUpgrades %d, got %d", g1.NumUpgrades, g2.NumUpgrades)
	}
	if g2.Summator != g1.Summator {
		t.Errorf("Expected summator %d, got %d", g1.Summator, g2.Summator)
	}
	if len(g2.Generators) != len(g1.Generators) {
		t.Errorf("Expected %d generators, got %d", len(g1.Generators), len(g2.Generators))
	}

	// Verify generator states
	for i := range g1.Generators {
		if g2.Generators[i].Level != g1.Generators[i].Level {
			t.Errorf("Generator %d: expected level %d, got %d", 
				i, g1.Generators[i].Level, g2.Generators[i].Level)
		}
		if g2.Generators[i].Name != g1.Generators[i].Name {
			t.Errorf("Generator %d: expected name %s, got %s", 
				i, g1.Generators[i].Name, g2.Generators[i].Name)
		}
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "subdir", "nested", "test_save.json")

	g := New()
	g.Score = 100

	// Save should create the nested directories
	err := g.Save(saveFile)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		t.Fatal("Save file was not created in nested directory")
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	g := New()
	err := g.Load("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file, got nil")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "invalid.json")
	
	err := os.WriteFile(saveFile, []byte("not valid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	g := New()
	err = g.Load(saveFile)
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}
}

func TestSavePreservesGeneratorOutput(t *testing.T) {
	tmpDir := t.TempDir()
	saveFile := filepath.Join(tmpDir, "test_output.json")

	// Create and modify game
	g1 := New()
	g1.Score = 1000
	g1.BuyGenerator(0)
	g1.BuyGenerator(0)
	g1.BuyGenerator(1)

	output1 := g1.TotalProduction()

	// Save and load
	if err := g1.Save(saveFile); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	g2 := New()
	if err := g2.Load(saveFile); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	output2 := g2.TotalProduction()

	if output1 != output2 {
		t.Errorf("Expected total production %d, got %d", output1, output2)
	}
}

func TestGetDefaultSaveFile(t *testing.T) {
	path := GetDefaultSaveFile()
	if path == "" {
		t.Error("GetDefaultSaveFile returned empty string")
	}
	// Just verify it returns something reasonable
	if filepath.Base(path) != "savegame.json" {
		t.Errorf("Expected filename to be savegame.json, got %s", filepath.Base(path))
	}
}
