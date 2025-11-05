package entity

import (
	"testing"
)

func TestNewEnemy(t *testing.T) {
	e := NewEnemy(10.0, 15.0, TypeBrawler)

	if e.X != 10.0 {
		t.Errorf("Expected X to be 10.0, got %f", e.X)
	}

	if e.Y != 15.0 {
		t.Errorf("Expected Y to be 15.0, got %f", e.Y)
	}

	if e.Type != TypeBrawler {
		t.Errorf("Expected type to be TypeBrawler, got %v", e.Type)
	}

	if e.HealthVal != 30 {
		t.Errorf("Expected Brawler health to be 30, got %d", e.HealthVal)
	}
}

func TestEnemyInterface(t *testing.T) {
	e := NewEnemy(5.0, 10.0, TypeBrawler)

	if e.GetX() != 5.0 {
		t.Errorf("Expected GetX to return 5.0, got %f", e.GetX())
	}

	if e.GetY() != 10.0 {
		t.Errorf("Expected GetY to return 10.0, got %f", e.GetY())
	}

	if !e.IsEnemy() {
		t.Error("Expected IsEnemy to return true")
	}

	if !e.IsAlive() {
		t.Error("Expected enemy to be alive")
	}

	if e.Health() != 30 {
		t.Errorf("Expected Health to return 30, got %d", e.Health())
	}
}

func TestEnemyDamage(t *testing.T) {
	e := NewEnemy(0, 0, TypeBrawler)
	initialHealth := e.Health()

	e.TakeDamage(10)

	if e.Health() != initialHealth-10 {
		t.Errorf("Expected health to be %d, got %d", initialHealth-10, e.Health())
	}

	e.TakeDamage(100)

	if e.IsAlive() {
		t.Error("Expected enemy to be dead after taking fatal damage")
	}
}

func TestEnemyCollision(t *testing.T) {
	e := NewEnemy(5.0, 5.0, TypeBrawler)

	if !e.CollidesWith(5, 5, 1, 1) {
		t.Error("Expected collision at (5,5)")
	}

	if e.CollidesWith(10, 10, 1, 1) {
		t.Error("Expected no collision at (10,10)")
	}
}
