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
)

func doTick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	game     *game.Game
	uiState  UIState
	message  string
}

func initialModel() model {
	return model{
		choices: []string{"Pause", "Upgrade", "Quit"},
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
					model.choices = []string{pauseLabel, "Upgrade", "Quit"}
					model.uiState = MainMenu
					model.cursor = 0
				case "No", "Cancel":
					pauseLabel := "Pause"
					if model.game.IsPaused() {
						pauseLabel = "Resume"
					}
					model.choices = []string{pauseLabel, "Upgrade", "Quit"}
					model.uiState = MainMenu
					model.cursor = 0
				}
			}
		}
	}

	return model, nil
}

func (model model) View() string {
	title := "\nTerminal Idle Game\n"
	stats := fmt.Sprintf("| Score: %d | Upgrades: %d | Rate: %d/s |\n", model.game.Score, model.game.NumUpgrades, model.game.Summator)
	breakLine := "----------------------------------------"

	s := title + breakLine + "\n" + stats + breakLine + "\n\n"

	// Show message if any
	if model.message != "" {
		s += fmt.Sprintf("*** %s ***\n\n", model.message)
	}

	// Show upgrade cost in upgrade confirm screen
	if model.uiState == UpgradeConfirm {
		s += fmt.Sprintf("Upgrade Cost: %d points\n", model.game.UpgradeCost())
		s += fmt.Sprintf("Current Score: %d points\n\n", model.game.Score)
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

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
