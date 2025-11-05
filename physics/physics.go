package physics

import (
	"github.com/imjasonh/beatemup/game"
)

// ClampYToLanes snaps a Y coordinate to the nearest lane
func ClampYToLanes(y float64) float64 {
	minDist := 1000.0
	closestY := game.LaneHeights[0]

	// Find the closest lane height to snap to
	for _, laneY := range game.LaneHeights {
		dist := laneY - y
		if dist < 0 {
			dist = -dist
		} // abs

		if dist < minDist {
			minDist = dist
			closestY = laneY
		}
	}
	return closestY
}
