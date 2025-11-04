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

## Architecture

The game follows the design document specifications:

### Game Loop
- Uses Bubble Tea's reactive model with a custom `GameTickMsg` at 60 FPS
- All physics, movement, AI, and scrolling are tick-based
- Separate handling for player input (KeyMsg) and game updates (GameTickMsg)

### Data Structures
- **Player**: Health, position, velocity, state machine, and attack timers
- **Entity Interface**: Polymorphic system for enemies, power-ups, and projectiles
- **Model**: Core game state including player, entities, score, and camera

### Rendering
- Canvas-based rendering system
- Layer-based drawing (ground → entities → player)
- Styled output using lipgloss for colors and formatting

## Technical Details

- **Language**: Go 1.20+
- **Framework**: github.com/charmbracelet/bubbletea
- **Styling**: github.com/charmbracelet/lipgloss
- **FPS**: 60 frames per second
- **Lanes**: 4 vertical positions for depth simulation

## Development

Run tests:
```bash
go test -v
```

Build:
```bash
go build -o beatemup .
```

## Future Enhancements

The current implementation includes:
- ✅ Core game loop and tick system
- ✅ Player movement and attacks
- ✅ Basic enemy AI (Brawler type)
- ✅ Collision detection
- ✅ Scrolling and generation

Potential additions (from design doc):
- 9 additional enemy types with unique behaviors
- Power-up system
- Mini-boss encounters
- Combo system
- Sound effects (using beep library)
- Save/load high scores

## License

See repository license.
