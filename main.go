package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imjasonh/beatemup/entity"
	"github.com/imjasonh/beatemup/game"
	"github.com/imjasonh/beatemup/physics"
	"github.com/imjasonh/beatemup/player"
	"github.com/imjasonh/beatemup/renderer"
)

// Model represents the game state
type Model struct {
	// TUI State
	Width, Height int
	Styles        renderer.StyleSet

	// Game State
	Player        player.Player
	Entities      []entity.Entity
	Score         int
	WorldDistance float64
	CameraOffset  float64 // How much the world has scrolled (X position offset)
}

// GetPlayer implements the interface needed by enemy Update
func (m *Model) GetPlayer() interface{ GetX() float64; GetY() float64 } {
	return &m.Player
}

// Init initializes the model and starts the game tick
func (m Model) Init() tea.Cmd {
	// Seed the random generator for procedural generation
	rand.Seed(time.Now().UnixNano())
	return game.TickCmd()
}

// Update handles input and game tick messages (The Game Loop Update)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "w", "up":
			m.Player.VY = game.PlayerSpeed
		case "s", "down":
			m.Player.VY = -game.PlayerSpeed // Moving down decreases Y cell number
		case "a", "left":
			m.Player.VX = -game.PlayerSpeed
		case "d", "right":
			m.Player.VX = game.PlayerSpeed
		case "j": // Light Attack
			m.Player.State = player.StateLightAttack
			m.Player.AttackTimer = game.LightAttackDuration
		case "k": // Strong Attack
			m.Player.State = player.StateStrongAttack
			m.Player.AttackTimer = game.StrongAttackDuration
		}

	case game.GameTickMsg:
		// 1. --- APPLY PHYSICS & MOVEMENT ---
		m.updatePlayerMovement()
		m.updateEntities()

		// 2. --- CHECK COLLISIONS AND ATTACKS ---
		m.checkCollisions()

		// 3. --- HANDLE WORLD SCROLLING & GENERATION ---
		m.handleScrollingAndGeneration()

		// 4. --- UPDATE TIMERS / COOLDOWNS ---
		m.updateTimers()

		// Check for Game Over condition
		if m.Player.Health <= 0 {
			return m, tea.Quit // For simplicity, quit on death
		}

		return m, game.TickCmd() // Schedule the next tick

	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		m.Styles.GameArea = m.Styles.GameArea.Width(m.Width - 4).Height(m.Height - 5)
	}

	return m, cmd
}

// View renders the TUI (The Game Loop View)
func (m Model) View() string {
	// 1. Render the HUD
	hud := m.renderHUD()

	// 2. Render the Game Area
	gameArea := m.renderGameArea()

	return lipgloss.JoinVertical(lipgloss.Left, hud, gameArea)
}

// updatePlayerMovement updates player position and velocity
func (m *Model) updatePlayerMovement() {
	// Apply velocity to position
	m.Player.X += m.Player.VX
	m.Player.Y += m.Player.VY

	// Apply friction to slow down motion unless keys are pressed continuously
	m.Player.VX *= game.Friction
	m.Player.VY *= game.Friction

	// Clamp Y position to defined lanes for the beat 'em up feel
	m.Player.Y = physics.ClampYToLanes(m.Player.Y)

	// Ensure player does not move off-screen to the left (behind the camera)
	if m.Player.X < 0 {
		m.Player.X = 0
	}
}

// updateEntities updates all active entities
func (m *Model) updateEntities() {
	// Update all active entities (enemies move, projectiles fly, etc.)
	var nextEntities []entity.Entity
	for _, e := range m.Entities {
		e.Update(m) // Entity applies its own AI/movement
		if e.IsAlive() {
			nextEntities = append(nextEntities, e)
		}
	}
	m.Entities = nextEntities
}

// checkCollisions handles collision detection and damage
func (m *Model) checkCollisions() {
	// 1. Player Attack Collision
	if m.Player.State == player.StateLightAttack || m.Player.State == player.StateStrongAttack {
		attackDamage := 10
		attackWidth, attackHeight := 1, 1

		// Define the hitbox location relative to the player
		hitX := m.Player.X + float64(m.Player.Width)
		hitY := m.Player.Y

		if m.Player.State == player.StateStrongAttack {
			attackDamage = 25
			attackWidth, attackHeight = 2, 2
		}

		for _, e := range m.Entities {
			if e.IsEnemy() && e.CollidesWith(hitX, hitY, attackWidth, attackHeight) {
				e.TakeDamage(attackDamage)
				m.Score += 10
			}
		}
	}

	// 2. Enemy Attack / Touch Collision (simplified: enemy touch hurts)
	for _, e := range m.Entities {
		if e.IsEnemy() && e.CollidesWith(m.Player.X, m.Player.Y, m.Player.Width, m.Player.Height) {
			m.Player.Health -= 1
		}
	}

	// Cleanup entities that are now dead
	m.Entities = filterAliveEntities(m.Entities)
}

// handleScrollingAndGeneration manages world scrolling and enemy spawning
func (m *Model) handleScrollingAndGeneration() {
	// World Scrolling
	if m.Player.X > float64(m.Width)*game.ScrollingThreshold {
		scrollAmount := m.Player.X - float64(m.Width)*game.ScrollingThreshold
		m.CameraOffset += scrollAmount
		m.Player.X -= scrollAmount
		m.WorldDistance += scrollAmount * 0.1 // Update distance traveled

		// Trigger procedural generation of a new segment off-screen
		m.generateNewSegment()
	}
}

// updateTimers updates attack timers and cooldowns
func (m *Model) updateTimers() {
	// Update attack timer
	if m.Player.AttackTimer > 0 {
		// Decrement by the time passed in one tick
		m.Player.AttackTimer -= game.TickRate
	} else if m.Player.State != player.StateIdle {
		// Reset to idle once attack is done
		m.Player.State = player.StateIdle
		m.Player.AttackTimer = 0
	}

	// Reset velocity if no input is being applied (handled by friction, but V=0 if player stopped)
	if m.Player.VX < 0.01 && m.Player.VX > -0.01 {
		m.Player.VX = 0
	}
	if m.Player.VY < 0.01 && m.Player.VY > -0.01 {
		m.Player.VY = 0
	}
}

// generateNewSegment spawns new enemies procedurally
func (m *Model) generateNewSegment() {
	spawnX := m.CameraOffset + float64(m.Width) + 5.0

	// Simplified spawn logic: always spawn one Brawler
	enemy := entity.NewEnemy(spawnX, game.LaneHeights[rand.Intn(len(game.LaneHeights))], entity.TypeBrawler)
	m.Entities = append(m.Entities, enemy)
}

// renderHUD renders the heads-up display
func (m *Model) renderHUD() string {
	health := m.Player.Health
	maxHealth := m.Player.MaxHealth

	// Build Health Bar
	fullBlocks := health / 10 // Assuming 100 max health, 10 blocks
	emptyBlocks := (maxHealth / 10) - fullBlocks
	healthBar := strings.Repeat("█", fullBlocks) + strings.Repeat("░", emptyBlocks)

	healthView := fmt.Sprintf("HEALTH: %s", m.Styles.HealthBarFull.Render(healthBar))
	scoreView := fmt.Sprintf("SCORE: %d", m.Score)
	distanceView := fmt.Sprintf("DISTANCE: %.0fM", m.WorldDistance)

	hudContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		scoreView, " | ",
		healthView, " | ",
		distanceView,
	)

	return m.Styles.HUD.Render(hudContent)
}

// renderGameArea renders the main game area
func (m *Model) renderGameArea() string {
	canvas := make([][]string, m.Height)
	for i := range canvas {
		// Initialize line with spaces
		canvas[i] = make([]string, m.Width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	// Draw the ground line (visual aid) - draw it first so entities can render on top
	groundRow := int(game.GroundLevel)
	if groundRow >= 0 && groundRow < len(canvas) {
		for j := range canvas[groundRow] {
			canvas[groundRow][j] = "―"
		}
	}

	// Draw Entities (enemies, power-ups, etc.)
	for _, e := range m.Entities {
		// Offset entity position by CameraOffset before drawing
		renderX := int(e.GetX() - m.CameraOffset)
		renderY := int(e.GetY())

		if renderX >= 0 && renderX < m.Width {
			m.drawChar(canvas, renderX, renderY, e.View())
		}
	}

	// Draw Player on top of everything
	playerChar := "P"
	if m.Player.State == player.StateLightAttack {
		playerChar = "P!"
	}
	if m.Player.State == player.StateStrongAttack {
		playerChar = "P!!"
	}

	m.drawChar(canvas, int(m.Player.X), int(m.Player.Y), m.Styles.Player.Render(playerChar))

	// Combine lines into a single string
	var sb strings.Builder
	for _, line := range canvas {
		sb.WriteString(strings.Join(line, ""))
		sb.WriteString("\n")
	}

	return m.Styles.GameArea.Render(sb.String())
}

// drawChar safely draws a character string on the canvas
func (m *Model) drawChar(canvas [][]string, x, y int, char string) {
	if y >= 0 && y < len(canvas) && x >= 0 && x < len(canvas[y]) {
		canvas[y][x] = char
	}
}

// filterAliveEntities filters out dead entities
func filterAliveEntities(entities []entity.Entity) []entity.Entity {
	var alive []entity.Entity
	for _, e := range entities {
		if e.IsAlive() {
			alive = append(alive, e)
		}
	}
	return alive
}

// initialModel creates the initial game state
func initialModel() Model {
	p := player.New()

	m := Model{
		Player:        p,
		Styles:        renderer.DefaultStyles(),
		WorldDistance: 0,
		Score:         0,
		CameraOffset:  0,
	}

	// Spawn initial enemies (for testing)
	m.Entities = append(m.Entities, entity.NewEnemy(25.0, game.LaneHeights[0], entity.TypeBrawler))
	m.Entities = append(m.Entities, entity.NewEnemy(35.0, game.LaneHeights[2], entity.TypeBrawler))

	return m
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting game: %v\n", err)
		os.Exit(1)
	}
}
