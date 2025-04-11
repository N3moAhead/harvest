package enemy

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CarrotEnemy struct {
	BaseEnemy
	MeleeEnemyData
}

func NewCarrotEnemy(pos component.Vector2D) *CarrotEnemy {
	return &CarrotEnemy{
		BaseEnemy: BaseEnemy{
			Pos:    pos,
			Speed:  config.CARROT_SPEED,
			Health: component.NewHealth(config.CARROT_HEALTH),
			Damage: config.CARROT_DAMAGE,
			AttackCooldown: config.CARROT_ATTACK_COOLDOWN,
			attackTimer:   config.CARROT_ATTACK_START,
		},
		MeleeEnemyData: MeleeEnemyData{
			AttackRange:    config.CARROT_ATTACK_RANGE,
		},
	}
}

func (e *CarrotEnemy) Update(player *player.Player, dt float64) {
	e.MoveTowards(player.Pos, dt)

	e.attackTimer -= dt
	if e.Pos.Sub(player.Pos).Len() < e.AttackRange && e.attackTimer <= 0 {
		player.Health.Damage(e.Damage)
		e.attackTimer = e.AttackCooldown
	}
}

func (e *CarrotEnemy) Draw(screen *ebiten.Image, camX, camY float64) {
	x, y := e.Pos.X-camX, e.Pos.Y-camY
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(16), float32(16), color.RGBA{255, 128, 0, 255}, false)
}

type PeashooterEnemy struct {
	BaseEnemy
	RangedEnemyData
}
