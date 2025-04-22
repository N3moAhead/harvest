package enemy

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// type PeashooterEnemy struct { Different example
// 	Enemy
// 	RangedEnemyData
// }

type CarrotEnemy struct {
	Enemy
	MeleeEnemyData
}

func NewCarrotEnemy(pos component.Vector2D) *CarrotEnemy {
	baseEntity := entity.NewEntity(pos.X, pos.Y)
	return &CarrotEnemy{
		Enemy: Enemy{
			Entity:         *baseEntity,
			Speed:          config.CARROT_SPEED,
			Health:         component.NewHealth(config.CARROT_HEALTH),
			Damage:         config.CARROT_DAMAGE,
			AttackCooldown: config.CARROT_ATTACK_COOLDOWN,
			attackTimer:    config.CARROT_ATTACK_START,
		},
		MeleeEnemyData: MeleeEnemyData{
			AttackRange: config.CARROT_ATTACK_RANGE,
		},
	}
}

func (e *CarrotEnemy) Update(player *player.Player, dt float64) {
	e.UpdateKnockback()
	e.MoveTowards(player.Pos, dt)

	e.attackTimer -= dt
	if e.Pos.Sub(player.Pos).Len() < e.AttackRange && e.attackTimer <= 0 {
		player.Health.Damage(e.Damage)
		e.attackTimer = e.AttackCooldown
	}
}

func (e *CarrotEnemy) Draw(screen *ebiten.Image, camX, camY float64) {
	e.DefaultDraw(screen, camX, camY, config.CARROT_WIDTH, config.CARROT_HEIGHT,
		color.RGBA{R: config.CARROT_COLOR_R, G: config.CARROT_COLOR_G, B: config.CARROT_COLOR_B, A: 255})
}

func (e *CarrotEnemy) IsAlive() bool {
	return e.Enemy.IsAlive()
}

func (e *CarrotEnemy) GetPosition() component.Vector2D {
	return e.Enemy.GetPosition()
}
