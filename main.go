package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time
type GameCurrentState int

const (
	Running GameCurrentState = iota
	Paused
	Upgrading
)

func doTick() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type GameState struct {
	score       int
	numUpgrades int
	summator    int
	state       GameCurrentState
}

var GameStateName = map[GameCurrentState]string{
	Running:   "Running",
	Paused:    "Paused",
	Upgrading: "Upgrading",
}

type model struct {
	cursor    int
	choices   []string
	selected  map[int]struct{}
	gameState GameState
}

func initialModel() model {
	return model{
		choices: []string{"Pause", "Upgrade", "Quit"},
		gameState: GameState{
			score:       0,
			numUpgrades: 0,
			summator:    2,
			state:       Running,
		},
		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (model *model) upgradeSummator() {
	model.gameState.numUpgrades++
	model.gameState.summator *= 2
}

func (model model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle("Terminal Idle Game"), doTick())
}

func (model model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if model.gameState.state == Running {
			model.gameState.score += model.gameState.summator
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
			switch model.gameState.state {
			case Running:
				switch model.choices[model.cursor] {
				case "Pause":
					model.choices[0] = "Resume"
					model.gameState.state = Paused
				case "Upgrade":
					model.choices[0] = "Yes"
					model.choices[1] = "No"
					model.gameState.state = Upgrading
				case "Quit":
					return model, tea.Quit
				}
			case Paused:
				switch model.choices[model.cursor] {
				case "Resume":
					model.choices[0] = "Pause"
					model.gameState.state = Running
				case "Upgrade":
					model.choices[0] = "Yes"
					model.choices[1] = "No"
					model.gameState.state = Upgrading
				case "Quit":
					return model, tea.Quit
				}
			case Upgrading:
				switch model.choices[model.cursor] {
				case "Yes":
					model.choices[0] = "Pause"
					model.choices[1] = "Upgrade"
					if model.gameState.score >= 10*model.gameState.summator {
						model.upgradeSummator()
						model.gameState.state = Running
					} else {
						go func() {
							time.Sleep(time.Second * 2)
							model.choices[0] = "Pause"
							model.choices[1] = "Upgrade"
						}()
						model.choices[0] = "Not enough score points..."
						model.choices[1] = fmt.Sprintf("Need %d or more points", 10*model.gameState.summator)
						model.gameState.state = Running
					}
				case "No":
					model.choices[0] = "Pause"
					model.choices[1] = "Upgrade"
					model.gameState.state = Running
				case "Quit":
					return model, tea.Quit
				}
			}
		}
	}

	return model, nil
}

func (model model) View() string {
	title := "\nTerminal Idle Game\n"
	stats := fmt.Sprintf("| Score: %d | Number of Upgrades: %d | Summator: %d |\n", model.gameState.score, model.gameState.numUpgrades, model.gameState.summator)
	break_line := "--------------------------------"

	s := title + break_line + break_line + "\n" + stats + break_line + break_line + "\n\n"

	for i, choice := range model.choices {
		cursor := " "
		if model.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s |%s|\n", cursor, choice)
	}

	s += "\n" + break_line + "\n"
	s += "cursor: " + strconv.Itoa(model.cursor) + " -- selected: "
	if _, ok := model.selected[model.cursor]; ok {
		s += "true"
	} else {
		s += "false"
	}

	s += "\n"

	s += "state: " + GameStateName[model.gameState.state] + "\n"

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
