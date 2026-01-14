# terminal-idle-game
An idle game made for the terminal, developed in Golang

## Features

- **Idle Game Mechanics**: Watch your score increase automatically over time
- **Pause/Resume**: Control the game flow
- **Base Rate Upgrades**: Spend points to double your base production rate
- **Multiple Generators**: Purchase and upgrade 4 different generator types:
  - **Clicker**: Low cost, produces 1/s per level
  - **Farm**: Medium cost, produces 5/s per level
  - **Mine**: High cost, produces 20/s per level
  - **Factory**: Very high cost, produces 50/s per level
- **Save/Load**: Game progress is automatically loaded on startup and can be manually saved/loaded
- **Intuitive UI**: Navigate with arrow keys or vi-style keybindings (j/k)

## Installation

```bash
go build -o terminal-idle-game .
```

## Usage

```bash
./terminal-idle-game
```

### Controls

- **Arrow Keys** or **j/k**: Navigate menu
- **Enter** or **Space**: Select option
- **q**: Quit game

## Architecture

The project is organized into packages:

- `game/`: Core game logic
  - `game.go`: Main game state and mechanics
  - `generator.go`: Generator types and logic
  - `save.go`: Save/load functionality
  - `*_test.go`: Comprehensive test suite (22 test functions)

- `main.go`: Terminal UI using Bubble Tea

## Testing

Run all tests:

```bash
go test -v ./...
```

All game logic is thoroughly tested with unit tests covering:
- Game state management
- Upgrade mechanics
- Generator functionality
- Save/load persistence
- Edge cases and error handling

## Save File Location

Game saves are stored at: `~/.terminal-idle-game/savegame.json`

## Development

This game demonstrates:
- Clean separation of concerns (game logic vs UI)
- Comprehensive test coverage
- Idiomatic Go code structure
- Terminal UI with Bubble Tea framework
