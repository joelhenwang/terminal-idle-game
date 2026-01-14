package main

import (
	"fmt"
	"os"
	"terminal-idle-game/game"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time
type UIState int

const (
	MainMenu UIState = iota
	UpgradeConfirm
	GeneratorsMenu
	BuyGeneratorConfirm
)

func doTick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type model struct {
	cursor            int
	choices           []string
	selected          map[int]struct{}
	game              *game.Game
	uiState           UIState
	message           string
	selectedGenerator int
}

func initialModel() model {
	return model{
		choices: []string{"Pause", "Upgrade", "Generators", "Quit"},
		game:    game.New(),
		uiState: MainMenu,
		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (model model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Terminal Idle Game"), doTick())
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		model.game.Tick()
		// Clear message after a few seconds
		if model.message != "" {
			// Message will be cleared on next interaction
		}
		return model, doTick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "up", "k":
			if model.cursor > 0 {
				model.cursor--
			}
		case "down", "j":
			if model.cursor < len(model.choices)-1 {
				model.cursor++
			}
		case "enter", " ":
			model.message = "" // Clear message on any action
			switch model.uiState {
			case MainMenu:
				switch model.choices[model.cursor] {
				case "Pause":
					model.game.TogglePause()
					model.choices[0] = "Resume"
				case "Resume":
					model.game.TogglePause()
					model.choices[0] = "Pause"
				case "Upgrade":
					model.choices = []string{"Yes", "No", "Cancel"}
					model.uiState = UpgradeConfirm
					model.cursor = 0
				case "Generators":
					model.choices = []string{}
					for _, gen := range model.game.Generators {
						model.choices = append(model.choices, gen.Name)
					}
					model.choices = append(model.choices, "Back")
					model.uiState = GeneratorsMenu
					model.cursor = 0
				case "Quit":
					return model, tea.Quit
				}
			case UpgradeConfirm:
				switch model.choices[model.cursor] {
				case "Yes":
					if model.game.Upgrade() {
						model.message = "Upgrade successful!"
					} else {
						model.message = fmt.Sprintf("Not enough score! Need %d points", model.game.UpgradeCost())
					}
					pauseLabel := "Pause"
					if model.game.IsPaused() {
						pauseLabel = "Resume"
					}
					model.choices = []string{pauseLabel, "Upgrade", "Generators", "Quit"}
					model.uiState = MainMenu
					model.cursor = 0
				case "No", "Cancel":
					pauseLabel := "Pause"
					if model.game.IsPaused() {
						pauseLabel = "Resume"
					}
					model.choices = []string{pauseLabel, "Upgrade", "Generators", "Quit"}
					model.uiState = MainMenu
					model.cursor = 0
				}
			case GeneratorsMenu:
				if model.choices[model.cursor] == "Back" {
					pauseLabel := "Pause"
					if model.game.IsPaused() {
						pauseLabel = "Resume"
					}
					model.choices = []string{pauseLabel, "Upgrade", "Generators", "Quit"}
					model.uiState = MainMenu
					model.cursor = 0
				} else {
					model.selectedGenerator = model.cursor
					model.choices = []string{"Buy", "Cancel"}
					model.uiState = BuyGeneratorConfirm
					model.cursor = 0
				}
			case BuyGeneratorConfirm:
				switch model.choices[model.cursor] {
				case "Buy":
					if model.game.BuyGenerator(model.selectedGenerator) {
						gen := model.game.Generators[model.selectedGenerator]
						model.message = fmt.Sprintf("Bought %s! Now level %d", gen.Name, gen.Level)
					} else {
						gen := model.game.Generators[model.selectedGenerator]
						model.message = fmt.Sprintf("Not enough score! Need %d points", gen.Cost())
					}
					model.choices = []string{}
					for _, gen := range model.game.Generators {
						model.choices = append(model.choices, gen.Name)
					}
					model.choices = append(model.choices, "Back")
					model.uiState = GeneratorsMenu
					model.cursor = 0
				case "Cancel":
					model.choices = []string{}
					for _, gen := range model.game.Generators {
						model.choices = append(model.choices, gen.Name)
					}
					model.choices = append(model.choices, "Back")
					model.uiState = GeneratorsMenu
					model.cursor = 0
				}
			}
		}
	}

	return model, nil
}

func (model model) View() string {
	title := "\nTerminal Idle Game\n"
	stats := fmt.Sprintf("| Score: %d | Upgrades: %d | Base Rate: %d/s | Total Rate: %d/s |\n", 
		model.game.Score, model.game.NumUpgrades, model.game.Summator, model.game.TotalProduction())
	breakLine := "=========================================="

	s := title + breakLine + "\n" + stats + breakLine + "\n\n"

	// Show message if any
	if model.message != "" {
		s += fmt.Sprintf("*** %s ***\n\n", model.message)
	}

	// Show upgrade cost in upgrade confirm screen
	if model.uiState == UpgradeConfirm {
		s += fmt.Sprintf("Upgrade Cost: %d points (doubles base rate)\n", model.game.UpgradeCost())
		s += fmt.Sprintf("Current Score: %d points\n\n", model.game.Score)
	}

	// Show generator details in generators menu
	if model.uiState == GeneratorsMenu {
		s += "GENERATORS:\n"
		for _, gen := range model.game.Generators {
			s += fmt.Sprintf("  %s - Level %d - Produces %d/s - Cost: %d\n", 
				gen.Name, gen.Level, gen.Output(), gen.Cost())
		}
		s += "\n"
	}

	// Show selected generator details in buy confirm screen
	if model.uiState == BuyGeneratorConfirm {
		gen := model.game.Generators[model.selectedGenerator]
		s += fmt.Sprintf("Generator: %s\n", gen.Name)
		s += fmt.Sprintf("Current Level: %d\n", gen.Level)
		s += fmt.Sprintf("Current Output: %d/s\n", gen.Output())
		s += fmt.Sprintf("Next Output: %d/s\n", gen.BaseOutput*(gen.Level+1))
		s += fmt.Sprintf("Cost: %d points\n", gen.Cost())
		s += fmt.Sprintf("Your Score: %d points\n\n", model.game.Score)
	}

	for i, choice := range model.choices {
		cursor := " "
		if model.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s |%s|\n", cursor, choice)
	}

	s += "\n" + breakLine + "\n"

	stateStr := "Running"
	if model.game.IsPaused() {
		stateStr = "Paused"
	}
	s += "State: " + stateStr + "\n"

	s += "\nPress q to quit. Use arrow keys or j/k to navigate.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
