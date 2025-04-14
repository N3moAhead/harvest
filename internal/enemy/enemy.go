package enemy

import (
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyInterface interface {
	Update(player *player.Player, dt float64)
	Draw(screen *ebiten.Image, camX, camY float64)
	GetPosition() component.Vector2D
	IsAlive() bool
}
type Enemy struct {
	Pos    component.Vector2D
	Health component.Health
	Speed  float64
	Damage int
	AttackCooldown  float64
	attackTimer     float64
}

func (e *Enemy) MoveTowards(target component.Vector2D, dt float64) {
	dir := target.Sub(e.Pos)
	if dir.Len() > 1 {
		dir = dir.Normalize()
		e.Pos = e.Pos.Add(dir.Mul(e.Speed * dt))
	}
}

func (e *Enemy) IsAlive() bool {
	return e.Health.HP > 0
}

func (e *Enemy) GetPosition() component.Vector2D {
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