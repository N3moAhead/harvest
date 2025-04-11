package enemy

import (
	"github.com/N3moAhead/harvest/internal/component"
)

type BaseEnemy struct {
	Pos    component.Vector2D
	Speed  float64
	Health component.Health
	Damage int
	AttackCooldown  float64
	attackTimer     float64
}

func (e *BaseEnemy) MoveTowards(target component.Vector2D, dt float64) {
	dir := target.Sub(e.Pos)
	if dir.Len() > 1 {
		dir = dir.Normalize()
		e.Pos = e.Pos.Add(dir.Mul(e.Speed * dt))
	}
}

func (e *BaseEnemy) IsAlive() bool {
	return e.Health.HP > 0
}

func (e *BaseEnemy) GetPosition() component.Vector2D {
	return e.Pos
}

// seperation into these types must be discussed
type MeleeEnemyData struct {
	AttackRange    float64
}

type RangedEnemyData struct {
	ProjectileSpeed float64
	// more
}