package renderer

import (
	"github.com/charmbracelet/lipgloss"
)

// StyleSet holds all the lipgloss styles for the game
type StyleSet struct {
	HUD            lipgloss.Style
	HealthBarFull  lipgloss.Style
	HealthBarEmpty lipgloss.Style
	Player         lipgloss.Style
	EnemyDefault   lipgloss.Style
	MiniBoss       lipgloss.Style
	GameArea       lipgloss.Style
}

// DefaultStyles returns the default style set
func DefaultStyles() StyleSet {
	return StyleSet{
		HUD:            lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("#56007A")).Foreground(lipgloss.Color("#FFF8C8")),
		HealthBarFull:  lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")),
		HealthBarEmpty: lipgloss.NewStyle().Foreground(lipgloss.Color("#7F1D1D")),
		Player:         lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")),
		EnemyDefault:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500")),
		MiniBoss:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true),
		GameArea:       lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).BorderForeground(lipgloss.Color("#7F1D1D")),
	}
}
