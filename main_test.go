package main

import (
	"testing"

	"github.com/imjasonh/beatemup/game"
)

func TestInitialModel(t *testing.T) {
	m := initialModel()

	if m.Player.Health != 100 {
		t.Errorf("Expected player health to be 100, got %d", m.Player.Health)
	}

	if m.Player.MaxHealth != 100 {
		t.Errorf("Expected player max health to be 100, got %d", m.Player.MaxHealth)
	}

	if m.Player.X != 10.0 {
		t.Errorf("Expected player X to be 10.0, got %f", m.Player.X)
	}

	if m.Score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", m.Score)
	}

	if len(m.Entities) != 2 {
		t.Errorf("Expected 2 initial entities, got %d", len(m.Entities))
	}
}

func TestFilterAliveEntities(t *testing.T) {
	// This is tested indirectly through the game logic
	// but we can add a simple test here
	m := initialModel()
	before := len(m.Entities)

	filtered := filterAliveEntities(m.Entities)

	if len(filtered) != before {
		t.Errorf("Expected %d alive entities, got %d", before, len(filtered))
	}
}

func TestGameConstants(t *testing.T) {
	if game.FPS != 60 {
		t.Errorf("Expected FPS to be 60, got %d", game.FPS)
	}

	if len(game.LaneHeights) != 4 {
		t.Errorf("Expected 4 lanes, got %d", len(game.LaneHeights))
	}
}
