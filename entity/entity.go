package entity

// Entity represents any dynamic object in the game
type Entity interface {
	GetX() float64
	GetY() float64
	Update(m interface{}) // Uses interface{} to avoid circular dependency
	View() string
	// Required fields for collision/damage
	IsEnemy() bool
	TakeDamage(damage int)
	Health() int
	IsAlive() bool
	CollidesWith(x, y float64, w, h int) bool
}
