package player

import (
	"testing"

	"github.com/imjasonh/beatemup/game"
)

func TestNew(t *testing.T) {
	p := New()

	if p.Health != 100 {
		t.Errorf("Expected health to be 100, got %d", p.Health)
	}

	if p.MaxHealth != 100 {
		t.Errorf("Expected max health to be 100, got %d", p.MaxHealth)
	}

	if p.X != 10.0 {
		t.Errorf("Expected X to be 10.0, got %f", p.X)
	}

	if p.Y != game.GroundLevel {
		t.Errorf("Expected Y to be %f, got %f", game.GroundLevel, p.Y)
	}
}

func TestPlayerCollision(t *testing.T) {
	p := Player{X: 5, Y: 5, Width: 2, Height: 2}

	// Test collision
	if !p.CollidesWith(6, 6, 2, 2) {
		t.Error("Expected collision at (6,6)")
	}

	// Test no collision
	if p.CollidesWith(10, 10, 2, 2) {
		t.Error("Expected no collision at (10,10)")
	}
}

func TestPlayerGetters(t *testing.T) {
	p := Player{X: 15.5, Y: 20.3}

	if p.GetX() != 15.5 {
		t.Errorf("Expected GetX to return 15.5, got %f", p.GetX())
	}

	if p.GetY() != 20.3 {
		t.Errorf("Expected GetY to return 20.3, got %f", p.GetY())
	}
}
