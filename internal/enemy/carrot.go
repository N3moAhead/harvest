package enemy

import (
	"fmt"
	"image/color"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
)

type CarrotEnemy struct {
	Enemy
	MeleeEnemyData
}

func NewCarrotEnemy(pos component.Vector2D) *CarrotEnemy {
	carrotSprite, ok := assets.AssetStore.GetImage("carrot")
	animationStore := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 32, 8, 6, true)
		if err == nil {
			animationStore.AddAnimation("walkRight", walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 6*32, 8, 6, true)
		if err == nil {
			animationStore.AddAnimation("walkLeft", walkLeftAnimation)
		}
		attackAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 4*32, 2, 6, false)
		if err == nil {
			animationStore.AddAnimation("attack", attackAnimation)
		}
		spawnAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 2*32, 6, 10, false)
		if err == nil {
			animationStore.AddAnimation("spawn", spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 7*32, 6, 10, false)
		if err == nil {
			animationStore.AddAnimation("death", deathAnimation)
		}

		ok := animationStore.SetCurrentAnimation("spawn")
		if !ok {
			fmt.Println("Warning: Unable to start the spawning animation")
		}
	}
	baseEntity := entity.NewEntity(pos.X, pos.Y)
	return &CarrotEnemy{
		Enemy: Enemy{
			Entity:         *baseEntity,
			Speed:          config.CARROT_SPEED,
			Health:         component.NewHealth(config.CARROT_HEALTH),
			Damage:         config.CARROT_DAMAGE,
			AttackCooldown: config.CARROT_ATTACK_COOLDOWN,
			attackTimer:    config.CARROT_ATTACK_START,
			animationStore: animationStore,
		},
		MeleeEnemyData: MeleeEnemyData{
			AttackRange: config.CARROT_ATTACK_RANGE,
		},
	}
}

func (e *CarrotEnemy) SetWalkingAnimation(player *player.Player) {
	dir := player.Pos.Sub(e.Pos)
	if dir.X > 0 {
		if ok := e.animationStore.SetCurrentAnimation("walkRight"); !ok {
			fmt.Println("Warning: Unable to start the carrot walkRight animation")
		}
	} else {
		if ok := e.animationStore.SetCurrentAnimation("walkLeft"); !ok {
			fmt.Println("Warning: Unable to start the carrot walkLeft animation")
		}
	}
}

func (e *CarrotEnemy) Update(player *player.Player, dt float64) {
	e.animationStore.Update()

	// Only update the carrot if it is alive
	if e.Health.HP > 0 {
		// Set current animation to walking if no animation is currently running
		// Or Update the type of running animation if currently the animation is running
		animationName := e.animationStore.GetCurrentAnimationName()
		updateRunningType := animationName == "walkLeft" || animationName == "walkRight"
		if e.animationStore.GetCurrentAnimation().IsFinished() || updateRunningType {
			e.SetWalkingAnimation(player)
		}

		// The enemy does not move during the spawn animation
		if e.animationStore.GetCurrentAnimationName() != "spawn" {
			e.UpdateKnockback()
			e.MoveTowards(player.Pos, dt)

			e.attackTimer -= dt
			if e.Pos.Sub(player.Pos).Len() < e.AttackRange && e.attackTimer <= 0 {
				player.Health.Damage(e.Damage)
				e.attackTimer = e.AttackCooldown
				// Starting the attack animation
				if ok := e.animationStore.SetCurrentAnimation("attack"); !ok {
					fmt.Println("Warning: Unable to start the carrot attack animation")
				}
			}
		}
	} else {
		// The carrot is dead
		// So we start the death animation
		e.animationStore.SetCurrentAnimation("death")
		// We still want to see the knockback happening
		// It just feels way better when playing :)
		e.UpdateKnockback()
	}
}

func (e *CarrotEnemy) Draw(screen *ebiten.Image, camX, camY float64) {
	frameImage := e.animationStore.GetImage()
	if frameImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(e.Pos.X-camX-16.0, e.Pos.Y-camY-16.0)
		screen.DrawImage(frameImage, op)
	} else {
		e.DefaultDraw(screen, camX, camY, config.CARROT_WIDTH, config.CARROT_HEIGHT,
			color.RGBA{R: config.CARROT_COLOR_R, G: config.CARROT_COLOR_G, B: config.CARROT_COLOR_B, A: 255})
	}
}

func (e *CarrotEnemy) IsAlive() bool {
	if e.Enemy.IsAlive() {
		return true
	} else {
		// The carrot will be marked as dead after the death animation is finished
		return !e.animationStore.GetCurrentAnimation().IsFinished()
	}
}

func (e *CarrotEnemy) GetPosition() component.Vector2D {
	return e.Enemy.GetPosition()
}
