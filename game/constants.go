package game

import "time"

// Game constants
const (
	FPS                = 60
	TickRate           = time.Second / FPS
	PlayerSpeed        = 0.5 // Base speed per tick
	Friction           = 0.85 // Velocity decay factor
	ScrollingThreshold = 0.6 // Player X-position trigger for scrolling (e.g., 60% of screen)
	GroundLevel        = 15.0 // Default Y position
	LightAttackDuration  = 200 * time.Millisecond
	StrongAttackDuration = 500 * time.Millisecond
)

// Fixed Y-lanes for beat 'em up perspective (in terminal rows)
var LaneHeights = []float64{GroundLevel - 6.0, GroundLevel - 2.0, GroundLevel + 2.0, GroundLevel + 6.0}
