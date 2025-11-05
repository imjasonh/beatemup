# 🥋 ARCADE-GO - Beat 'Em Up TUI Game

A terminal-based side-scrolling beat 'em up game built with Go and Bubble Tea.

## Features

- **60 FPS Game Loop**: Smooth gameplay powered by time-based game ticks
- **2D Movement**: Navigate across 4 distinct lanes with WASD or arrow keys
- **Combat System**: Light (J) and strong (K) attacks with cooldowns
- **Enemy AI**: Intelligent enemies that pursue and attack the player
- **Procedural Generation**: Endless gameplay with dynamic enemy spawning
- **World Scrolling**: Smooth camera movement as you progress
- **Styled UI**: Beautiful terminal rendering with health bars and HUD

## Project Structure

The codebase is organized into modular packages:

```
beatemup/
├── game/           # Core game constants and messages
│   ├── constants.go    # Game configuration (FPS, speeds, lanes)
│   ├── messages.go     # Bubble Tea messages (GameTickMsg)
│   └── game_test.go    # Game package tests
├── player/         # Player character logic
│   ├── player.go       # Player struct and methods
│   └── player_test.go  # Player tests
├── entity/         # Game entities (enemies, items)
│   ├── entity.go       # Entity interface
│   ├── enemy.go        # Enemy implementation and AI
│   └── entity_test.go  # Entity tests
├── physics/        # Physics and collision utilities
│   ├── physics.go      # Movement and collision helpers
│   └── physics_test.go # Physics tests
├── renderer/       # Rendering and styling
│   ├── styles.go       # lipgloss style definitions
│   └── renderer_test.go # Renderer tests
└── main.go         # Entry point and game loop
```

## Installation

```bash
go build -o beatemup .
```

## How to Play

Run the game:
```bash
./beatemup
```

### Controls

- **W/Up Arrow**: Move up (to higher lane)
- **S/Down Arrow**: Move down (to lower lane)
- **A/Left Arrow**: Move left
- **D/Right Arrow**: Move right
- **J**: Light attack (quick, less damage)
- **K**: Strong attack (slower, more damage)
- **Q or Ctrl+C**: Quit game

### Objective

Survive as long as possible while defeating endless waves of enemies. Your score increases with each enemy defeated, and the difficulty scales with distance traveled.

## Development

Run all tests:
```bash
go test ./... -v
```

Run tests for a specific package:
```bash
go test ./game -v
go test ./player -v
go test ./entity -v
```

Build:
```bash
go build -o beatemup .
```

## Architecture

The game follows the design document specifications with a modular structure:

### Game Loop
- Uses Bubble Tea's reactive model with a custom `GameTickMsg` at 60 FPS
- All physics, movement, AI, and scrolling are tick-based
- Separate handling for player input (KeyMsg) and game updates (GameTickMsg)

### Packages
- **game**: Core constants, tick rate, and game messages
- **player**: Player state, movement, and collision
- **entity**: Entity interface and enemy AI implementations
- **physics**: Collision detection and lane snapping utilities
- **renderer**: lipgloss styles for terminal rendering
- **main**: Bubble Tea model, game loop, and rendering logic

## Technical Details

- **Language**: Go 1.20+
- **Framework**: github.com/charmbracelet/bubbletea
- **Styling**: github.com/charmbracelet/lipgloss
- **FPS**: 60 frames per second
- **Lanes**: 4 vertical positions for depth simulation

## Testing

The project includes comprehensive tests for all packages:
- Unit tests for game constants and messages
- Player movement and collision tests
- Enemy AI and damage system tests
- Physics and lane snapping tests
- Renderer style tests

Run all tests with coverage:
```bash
go test ./... -cover
```

## Future Enhancements

The current implementation includes:
- ✅ Core game loop and tick system
- ✅ Player movement and attacks
- ✅ Basic enemy AI (Brawler type)
- ✅ Collision detection
- ✅ Scrolling and generation
- ✅ Modular package structure
- ✅ Comprehensive test coverage

Potential additions (from design doc):
- 9 additional enemy types with unique behaviors
- Power-up system
- Mini-boss encounters
- Combo system
- Sound effects (using beep library)
- Save/load high scores

## License

See repository license.
