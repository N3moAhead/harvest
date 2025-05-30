package enemy

import (
	"image/color"
	"math/rand"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type EnemyInterface interface {
	Update(player *player.Player, dt float64)
	Draw(screen *ebiten.Image, camX, camY float64)
	GetPosition() component.Vector2D
	IsAlive() bool
	TakeDamage(damage float64)
	AddKnockback(from *component.Vector2D, distance float64)
	TryDrop(elapsedMinutes float32) []item.Item
	GetType() EnemyType
}

type EnemyType int

const (
	TypeCarrot EnemyType = iota
	TypePotato
	TypeOnion
	maxEnemyType // Move this type to also allow the spawning of TypePeashooter
	TypePeashooter
)

func (t EnemyType) String() string {
	switch t {
	case TypeCarrot:
		return "carrot"
	case TypePotato:
		return "potato"
	case TypeOnion:
		return "onion"
	case TypePeashooter:
		return "peashooter"
	default:
		return "unknown"
	}
}

func RandomEnemyType() EnemyType {
	return EnemyType(rand.Intn(int(maxEnemyType)))
}

type Enemy struct {
	entity.Entity
	enemyType           EnemyType
	Health              component.Health
	Knockback           component.Knockback
	Speed               float64
	Damage              float64
	AttackCooldown      float64
	attackTimer         float64
	animationStore      *animation.AnimationStore
	DropProb            float32
	DropAmount          int
	DropAmountPerMinute float32
}

func (e *Enemy) MoveTowards(target component.Vector2D, dt float64) {
	diff := target.Sub(e.Pos)
	dist := diff.Len()

	dir := diff.Normalize()
	step := dir.Mul(e.Speed * dt)

	if step.Len() > dist { // if overstep, set to target
		e.Pos = target
	} else {
		e.Pos = e.Pos.Add(step)
	}
}

func (e *Enemy) AddKnockback(from *component.Vector2D, dist float64) {
	e.Knockback.Init(from, &e.Pos, dist)
}

func (e *Enemy) UpdateKnockback() {
	e.Knockback.Update(&e.Pos)
}

func (e *Enemy) DefaultDraw(screen *ebiten.Image, camX, camY float64, width int, height int, color color.RGBA) {
	x := float32(e.Pos.X - camX)
	y := float32(e.Pos.Y - camY)
	vector.DrawFilledRect(
		screen,
		x, y,
		float32(width), float32(height),
		color,
		false,
	)
}

func (e *Enemy) IsAlive() bool {
	return e.Health.HP > 0
}

func (e *Enemy) TakeDamage(damage float64) {
	// TODO flash the dealt damage
	e.Health.Damage(damage)
}

func (e *Enemy) GetPosition() component.Vector2D {
	return e.Pos
}

func (e *Enemy) GetType() EnemyType {
	return e.enemyType
}

func DefaultDrop(elapsedMinutes float32, x, y float64) []item.Item {
	return []item.Item{*item.NewPotato(x, y)} // default drop is a potato
}

type MeleeEnemyData struct {
	AttackRange float64
}

type RangedEnemyData struct {
	ProjectileSpeed float64
	// more
}
