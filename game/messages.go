package game

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// GameTickMsg is a custom message for time-based updates (the game tick)
type GameTickMsg time.Time

// TickCmd returns a command that sends a GameTickMsg at the specified tick rate
func TickCmd() tea.Cmd {
	return tea.Tick(TickRate, func(t time.Time) tea.Msg {
		return GameTickMsg(t)
	})
}
