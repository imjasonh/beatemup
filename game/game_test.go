package game

import (
	"testing"
	"time"
)

func TestConstants(t *testing.T) {
	if FPS != 60 {
		t.Errorf("Expected FPS to be 60, got %d", FPS)
	}

	if TickRate != time.Second/60 {
		t.Errorf("Expected TickRate to be %v, got %v", time.Second/60, TickRate)
	}

	if len(LaneHeights) != 4 {
		t.Errorf("Expected 4 lanes, got %d", len(LaneHeights))
	}
}

func TestTickCmd(t *testing.T) {
	cmd := TickCmd()
	if cmd == nil {
		t.Error("TickCmd should not return nil")
	}
}
