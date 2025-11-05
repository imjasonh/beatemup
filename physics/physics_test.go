package physics

import (
	"testing"

	"github.com/imjasonh/beatemup/game"
)

func TestClampYToLanes(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, game.LaneHeights[0]},
		{game.GroundLevel, game.LaneHeights[1]}, // Should snap to lane 1
		{100.0, game.LaneHeights[3]},
		{game.LaneHeights[2], game.LaneHeights[2]}, // Already on a lane
	}

	for _, tt := range tests {
		result := ClampYToLanes(tt.input)
		if result != tt.expected {
			t.Errorf("ClampYToLanes(%f) = %f, expected %f", tt.input, result, tt.expected)
		}
	}
}

func TestClampYToLanesSnapsToNearest(t *testing.T) {
	// Test a value between lanes 0 and 1
	y := (game.LaneHeights[0] + game.LaneHeights[1]) / 2.0
	result := ClampYToLanes(y)

	// Should snap to one of the two nearest lanes
	if result != game.LaneHeights[0] && result != game.LaneHeights[1] {
		t.Errorf("ClampYToLanes(%f) = %f, expected one of the nearest lanes", y, result)
	}
}
