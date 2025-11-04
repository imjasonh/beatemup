package main

import (
    "testing"
)

func TestPlayerInitialization(t *testing.T) {
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
}

func TestLaneHeights(t *testing.T) {
    if len(LANE_HEIGHTS) != 4 {
        t.Errorf("Expected 4 lanes, got %d", len(LANE_HEIGHTS))
    }
}

func TestTickRate(t *testing.T) {
    if FPS != 60 {
        t.Errorf("Expected FPS to be 60, got %d", FPS)
    }
}
