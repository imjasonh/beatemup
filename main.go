package main

import (
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// --- CONSTANTS AND TYPES ---

const (
    FPS = 60
    TickRate = time.Second / FPS
    PlayerSpeed = 0.5   // Base speed per tick
    Friction = 0.85     // Velocity decay factor
    ScrollingThreshold = 0.6 // Player X-position trigger for scrolling (e.g., 60% of screen)
    GroundLevel = 15.0  // Default Y position
)

// Fixed Y-lanes for beat 'em up perspective (in terminal rows)
var LANE_HEIGHTS = []float64{GroundLevel - 6.0, GroundLevel - 2.0, GroundLevel + 2.0, GroundLevel + 6.0}

// Player States
type PlayerState int
const (
    StateIdle PlayerState = iota
    StateMove
    StateLightAttack
    StateStrongAttack
    StateHit
)

// Custom message for time-based updates (the game tick)
type GameTickMsg time.Time

// tickCmd returns a command that sends a GameTickMsg at the specified tick rate
func tickCmd() tea.Cmd {
    return tea.Tick(TickRate, func(t time.Time) tea.Msg {
        return GameTickMsg(t)
    })
}

// --- STYLES ---

type StyleSet struct {
    HUD lipgloss.Style
    HealthBarFull lipgloss.Style
    HealthBarEmpty lipgloss.Style
    Player lipgloss.Style
    EnemyDefault lipgloss.Style
    MiniBoss lipgloss.Style
    GameArea lipgloss.Style
}

func defaultStyles() StyleSet {
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

// --- ENTITY INTERFACE ---

// All dynamic game objects implement this interface
type Entity interface {
    GetX() float64
    GetY() float64
    Update(m *Model)
    View() string
    // Required fields for collision/damage
    IsEnemy() bool
    TakeDamage(damage int)
    Health() int
    IsAlive() bool
    CollidesWith(x, y float64, w, h int) bool
}

// --- PLAYER STRUCT ---

type Player struct {
    Health    int
    MaxHealth int
    X, Y      float64 // X and Y coordinates (Y is lane)
    VX, VY    float64 // Velocity
    State     PlayerState
    AttackTimer time.Duration // Time remaining for attack animation/cooldown
    Width     int
    Height    int
}

// Simple collision check for a cell-based TUI (2x2 bounding box)
func (p *Player) CollidesWith(x, y float64, w, h int) bool {
    pX, pY := int(p.X), int(p.Y)
    
    // Check if rectangles overlap (using integer cells)
    return pX < int(x)+w &&
        pX+p.Width > int(x) &&
        pY < int(y)+h &&
        pY+p.Height > int(y)
}

// --- MODEL (THE GAME STATE) ---

type Model struct {
    // TUI State
    Width, Height int
    Styles        StyleSet
    
    // Game State
    Player        Player
    Entities      []Entity
    Score         int
    WorldDistance float64
    CameraOffset  float64 // How much the world has scrolled (X position offset)
}

// --- CORE BUBBLE TEA IMPLEMENTATIONS ---

// Initializes the model and starts the game tick
func (m Model) Init() tea.Cmd {
    // Seed the random generator for procedural generation
    rand.Seed(time.Now().UnixNano())
    return tickCmd()
}

// Handles input and game tick messages (The Game Loop Update)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "w", "up":
            m.Player.VY = PlayerSpeed
        case "s", "down":
            m.Player.VY = -PlayerSpeed // Moving down decreases Y cell number
        case "a", "left":
            m.Player.VX = -PlayerSpeed
        case "d", "right":
            m.Player.VX = PlayerSpeed
        case "j": // Light Attack
            // Placeholder: Create a temporary Hitbox entity
            // In a real implementation, you'd check cooldown/state
            m.Player.State = StateLightAttack
            m.Player.AttackTimer = time.Duration(200) * time.Millisecond
        case "k": // Strong Attack
            m.Player.State = StateStrongAttack
            m.Player.AttackTimer = time.Duration(500) * time.Millisecond
        }

    case GameTickMsg:
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

        return m, tickCmd() // Schedule the next tick

    case tea.WindowSizeMsg:
        m.Width, m.Height = msg.Width, msg.Height
        m.Styles.GameArea = m.Styles.GameArea.Width(m.Width - 4).Height(m.Height - 5)
    }

    // Clear velocity if key is released (simplified model for terminal)
    // In a full game, you'd track key down/up events to handle persistent velocity.
    // For this skeleton, we will apply friction in updatePlayerMovement.

    return m, cmd
}

// Renders the TUI (The Game Loop View)
func (m Model) View() string {
    // 1. Render the HUD
    hud := m.renderHUD()

    // 2. Render the Game Area
    gameArea := m.renderGameArea()

    return lipgloss.JoinVertical(lipgloss.Left, hud, gameArea)
}

// --- GAME LOGIC FUNCTIONS ---

func (m *Model) updatePlayerMovement() {
    // Apply velocity to position
    m.Player.X += m.Player.VX
    m.Player.Y += m.Player.VY

    // Apply friction to slow down motion unless keys are pressed continuously
    m.Player.VX *= Friction
    m.Player.VY *= Friction

    // Clamp Y position to defined lanes for the beat 'em up feel
    m.Player.Y = clampYToLanes(m.Player.Y)

    // Ensure player does not move off-screen to the left (behind the camera)
    if m.Player.X < 0 {
        m.Player.X = 0
    }
}

func (m *Model) updateEntities() {
    // Update all active entities (enemies move, projectiles fly, etc.)
    var nextEntities []Entity
    for _, e := range m.Entities {
        e.Update(m) // Entity applies its own AI/movement
        if e.IsAlive() {
            nextEntities = append(nextEntities, e)
        } else if !e.IsEnemy() {
            // Remove non-enemy entities like temporary hitboxes immediately
            // Dead enemies might need a final 'splode animation tick before removal
        }
    }
    m.Entities = nextEntities
}

func (m *Model) checkCollisions() {
    // Collision logic: Player attacks vs Enemies, Enemies vs Player

    // 1. Player Attack Collision
    if m.Player.State == StateLightAttack || m.Player.State == StateStrongAttack {
        attackDamage := 10
        attackWidth, attackHeight := 1, 1
        
        // Define the hitbox location relative to the player
        hitX := m.Player.X + float64(m.Player.Width)
        hitY := m.Player.Y
        
        if m.Player.State == StateStrongAttack {
            attackDamage = 25
            attackWidth, attackHeight = 2, 2
        }

        for _, e := range m.Entities {
            if e.IsEnemy() && e.CollidesWith(hitX, hitY, attackWidth, attackHeight) {
                e.TakeDamage(attackDamage)
                // Score update, hit effect...
                m.Score += 10 
            }
        }
    }
    
    // 2. Enemy Attack / Touch Collision (simplified: enemy touch hurts)
    for _, e := range m.Entities {
        if e.IsEnemy() && e.CollidesWith(m.Player.X, m.Player.Y, m.Player.Width, m.Player.Height) {
             // Simplified enemy damage upon touch
             // In a real game, you'd check if the enemy is in its attack state
             m.Player.Health -= 1
        }
    }
    
    // Cleanup entities that are now dead
    m.Entities = filterAliveEntities(m.Entities)
}

func (m *Model) handleScrollingAndGeneration() {
    // World Scrolling
    if m.Player.X > float64(m.Width) * ScrollingThreshold {
        scrollAmount := m.Player.X - float64(m.Width) * ScrollingThreshold
        m.CameraOffset += scrollAmount
        m.Player.X -= scrollAmount
        m.WorldDistance += scrollAmount * 0.1 // Update distance traveled
        
        // Trigger procedural generation of a new segment off-screen
        m.generateNewSegment()
    }
}

func (m *Model) updateTimers() {
    // Update attack timer
    if m.Player.AttackTimer > 0 {
        // Decrement by the time passed in one tick
        m.Player.AttackTimer -= TickRate 
    } else if m.Player.State != StateIdle {
        // Reset to idle once attack is done
        m.Player.State = StateIdle
        m.Player.AttackTimer = 0
    }

    // Reset velocity if no input is being applied (handled by friction, but V=0 if player stopped)
    if m.Player.VX < 0.01 && m.Player.VX > -0.01 { m.Player.VX = 0 }
    if m.Player.VY < 0.01 && m.Player.VY > -0.01 { m.Player.VY = 0 }
}

func (m *Model) generateNewSegment() {
    // Logic to procedurally generate enemies and power-ups
    // Spawn new enemies far to the right, using CameraOffset + Width
    
    spawnX := m.CameraOffset + float64(m.Width) + 5.0

    // Simplified spawn logic: always spawn one Brawler
    enemy := NewEnemy(spawnX, LANE_HEIGHTS[rand.Intn(len(LANE_HEIGHTS))], EnemyTypeBrawler)
    m.Entities = append(m.Entities, enemy)
    
    // You would implement complex, weighted spawning here based on m.WorldDistance
    // For Mini-Bosses (types 9 & 10), check if m.WorldDistance % 1000 == 0
}


// --- VIEW RENDERING FUNCTIONS ---

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
    groundRow := int(GroundLevel)
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
    if m.Player.State == StateLightAttack { playerChar = "P!" }
    if m.Player.State == StateStrongAttack { playerChar = "P!!" }
    
    m.drawChar(canvas, int(m.Player.X), int(m.Player.Y), m.Styles.Player.Render(playerChar))
    
    // Combine lines into a single string
    var sb strings.Builder
    for _, line := range canvas {
        sb.WriteString(strings.Join(line, ""))
        sb.WriteString("\n")
    }
    
    return m.Styles.GameArea.Render(sb.String())
}

// Utility to safely draw a character string on the canvas
func (m *Model) drawChar(canvas [][]string, x, y int, char string) {
    if y >= 0 && y < len(canvas) && x >= 0 && x < len(canvas[y]) {
        // Only draw the first character of the string for simplicity
        canvas[y][x] = char
    }
}

// --- UTILITY FUNCTIONS ---

func clampYToLanes(y float64) float64 {
    minDist := 1000.0
    closestY := LANE_HEIGHTS[0]
    
    // Find the closest lane height to snap to
    for _, laneY := range LANE_HEIGHTS {
        dist := laneY - y
        if dist < 0 { dist = -dist } // abs
        
        if dist < minDist {
            minDist = dist
            closestY = laneY
        }
    }
    return closestY
}

func filterAliveEntities(entities []Entity) []Entity {
    var alive []Entity
    for _, e := range entities {
        if e.IsAlive() {
            alive = append(alive, e)
        }
    }
    return alive
}

// --- PLACEHOLDER ENEMY IMPLEMENTATION (Brawler for testing) ---

type EnemyType int
const (
    EnemyTypeBrawler EnemyType = iota // Type 1: Simple X-only approach
    // ... define all 10 types
)

type Enemy struct {
    Type EnemyType
    HealthVal int
    X, Y float64
    Width, Height int
    // Other AI state fields
}

func NewEnemy(x, y float64, t EnemyType) *Enemy {
    health := 50
    if t == EnemyTypeBrawler { health = 30 }
    
    return &Enemy{
        Type: t,
        HealthVal: health,
        X: x, Y: y,
        Width: 1, Height: 1,
    }
}

func (e *Enemy) GetX() float64 { return e.X }
func (e *Enemy) GetY() float64 { return e.Y }
func (e *Enemy) IsEnemy() bool { return true }
func (e *Enemy) Health() int { return e.HealthVal }
func (e *Enemy) IsAlive() bool { return e.HealthVal > 0 }
func (e *Enemy) TakeDamage(damage int) { e.HealthVal -= damage }

func (e *Enemy) CollidesWith(x, y float64, w, h int) bool {
    eX, eY := int(e.X), int(e.Y)
    
    // Check if rectangles overlap (using integer cells)
    return eX < int(x)+w &&
        eX+e.Width > int(x) &&
        eY < int(y)+h &&
        eY+e.Height > int(y)
}

func (e *Enemy) Update(m *Model) {
    // Basic AI for Brawler: always move towards the player
    if e.Type == EnemyTypeBrawler {
        
        // Move horizontally towards player
        if e.X > m.Player.X {
            e.X -= PlayerSpeed * 0.5
        } else if e.X < m.Player.X {
            e.X += PlayerSpeed * 0.5
        }
        
        // Move vertically towards player's lane
        if e.Y > m.Player.Y {
            e.Y -= PlayerSpeed * 0.2
        } else if e.Y < m.Player.Y {
            e.Y += PlayerSpeed * 0.2
        }
        e.Y = clampYToLanes(e.Y) // Snap to the nearest lane
    }
    // All 10 enemy types would have their specific logic here.
}

func (e *Enemy) View() string {
    // Render based on type and health
    if e.HealthVal <= 0 {
        return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7F50")).Render("X") // Dead
    }
    
    char := "E"
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4500"))
    
    switch e.Type {
        case EnemyTypeBrawler: char = "B"
        // ...
    }
    
    return style.Render(char)
}

// --- INITIALIZATION ---

func initialModel() Model {
    p := Player{
        Health: 100, MaxHealth: 100,
        X: 10.0, Y: GroundLevel, // Start at X=10, middle lane
        Width: 1, Height: 1,
    }
    
    m := Model{
        Player: p,
        Styles: defaultStyles(),
        WorldDistance: 0,
        Score: 0,
        CameraOffset: 0,
    }

    // Spawn initial enemies (for testing)
    m.Entities = append(m.Entities, NewEnemy(25.0, LANE_HEIGHTS[0], EnemyTypeBrawler))
    m.Entities = append(m.Entities, NewEnemy(35.0, LANE_HEIGHTS[2], EnemyTypeBrawler))
    
    return m
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error starting game: %v\n", err)
        os.Exit(1)
    }
}

