package entity

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/imjasonh/beatemup/game"
	"github.com/imjasonh/beatemup/physics"
)

// Type represents the enemy type
type Type int

const (
	TypeBrawler Type = iota // Type 1: Simple X-only approach
	// ... define all 10 types
)

// Enemy represents an enemy character
type Enemy struct {
	Type      Type
	HealthVal int
	X, Y      float64
	Width     int
	Height    int
	// Other AI state fields
}

// NewEnemy creates a new enemy at the specified position
func NewEnemy(x, y float64, t Type) *Enemy {
	health := 50
	if t == TypeBrawler {
		health = 30
	}

	return &Enemy{
		Type:      t,
		HealthVal: health,
		X:         x,
		Y:         y,
		Width:     1,
		Height:    1,
	}
}

func (e *Enemy) GetX() float64       { return e.X }
func (e *Enemy) GetY() float64       { return e.Y }
func (e *Enemy) IsEnemy() bool       { return true }
func (e *Enemy) Health() int         { return e.HealthVal }
func (e *Enemy) IsAlive() bool       { return e.HealthVal > 0 }
func (e *Enemy) TakeDamage(damage int) { e.HealthVal -= damage }

func (e *Enemy) CollidesWith(x, y float64, w, h int) bool {
	eX, eY := int(e.X), int(e.Y)

	// Check if rectangles overlap (using integer cells)
	return eX < int(x)+w &&
		eX+e.Width > int(x) &&
		eY < int(y)+h &&
		eY+e.Height > int(y)
}

// Update updates the enemy AI and movement
func (e *Enemy) Update(m interface{}) {
	// Type assertion to access model fields
	type modelInterface interface {
		GetPlayer() interface{ GetX() float64; GetY() float64 }
	}
	
	model, ok := m.(modelInterface)
	if !ok {
		return
	}
	
	player := model.GetPlayer()

	// Basic AI for Brawler: always move towards the player
	if e.Type == TypeBrawler {
		// Move horizontally towards player
		if e.X > player.GetX() {
			e.X -= game.PlayerSpeed * 0.5
		} else if e.X < player.GetX() {
			e.X += game.PlayerSpeed * 0.5
		}

		// Move vertically towards player's lane
		if e.Y > player.GetY() {
			e.Y -= game.PlayerSpeed * 0.2
		} else if e.Y < player.GetY() {
			e.Y += game.PlayerSpeed * 0.2
		}
		e.Y = physics.ClampYToLanes(e.Y) // Snap to the nearest lane
	}
	// All 10 enemy types would have their specific logic here.
}

// View renders the enemy
func (e *Enemy) View() string {
	// Render based on type and health
	if e.HealthVal <= 0 {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7F50")).Render("X") // Dead
	}

	char := "E"
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))

	switch e.Type {
	case TypeBrawler:
		char = "B"
		// ...
	}

	return style.Render(char)
}
