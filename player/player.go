package player

import (
	"time"

	"github.com/imjasonh/beatemup/game"
)

// State represents the player's current action state
type State int

const (
	StateIdle State = iota
	StateMove
	StateLightAttack
	StateStrongAttack
	StateHit
)

// Player represents the player character
type Player struct {
	Health      int
	MaxHealth   int
	X, Y        float64       // X and Y coordinates (Y is lane)
	VX, VY      float64       // Velocity
	State       State
	AttackTimer time.Duration // Time remaining for attack animation/cooldown
	Width       int
	Height      int
}

// GetX returns the player's X position
func (p *Player) GetX() float64 { return p.X }

// GetY returns the player's Y position
func (p *Player) GetY() float64 { return p.Y }

// CollidesWith checks if the player collides with a bounding box
func (p *Player) CollidesWith(x, y float64, w, h int) bool {
	pX, pY := int(p.X), int(p.Y)

	// Check if rectangles overlap (using integer cells)
	return pX < int(x)+w &&
		pX+p.Width > int(x) &&
		pY < int(y)+h &&
		pY+p.Height > int(y)
}

// New creates a new player at the starting position
func New() Player {
	return Player{
		Health:    100,
		MaxHealth: 100,
		X:         10.0,
		Y:         game.GroundLevel, // Start at X=10, middle lane
		Width:     1,
		Height:    1,
	}
}
