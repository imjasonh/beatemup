package renderer

import (
	"testing"
)

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()

	// Check that styles are initialized by verifying they can render
	result := styles.HUD.Render("test")
	if result == "" {
		t.Error("Expected HUD style to render text")
	}

	result = styles.Player.Render("P")
	if result == "" {
		t.Error("Expected Player style to render text")
	}

	result = styles.GameArea.Render("area")
	if result == "" {
		t.Error("Expected GameArea style to render text")
	}
}
