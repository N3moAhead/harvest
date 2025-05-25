package enemy

import (
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type PotatoEnemy struct {
	Enemy
	MeleeEnemyData
}

func NewPotatoEnemy(pos component.Vector2D) *PotatoEnemy {
	potatoSprite, ok := assets.AssetStore.GetImage("potato")
	animationStore := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 1*32, 6, 6, true)
		if err == nil {
			animationStore.AddAnimation(WALK_RIGHT, walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 4*32, 6, 6, true)
		if err == nil {
			animationStore.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackRightAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 3*32, 7, 7, false)
		if err == nil {
			animationStore.AddAnimation(ATTACK_RIGHT, attackRightAnimation)
		}
		attackLeftAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 0*32, 7, 7, false)
		if err == nil {
			animationStore.AddAnimation(ATTACK_LEFT, attackLeftAnimation)
		}
		spawnAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 2*32, 8, 10, false)
		if err == nil {
			animationStore.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(potatoSprite, 32, 32, 0, 5*32, 7, 10, false)
		if err == nil {
			animationStore.AddAnimation(DEATH, deathAnimation)
		}

		ok := animationStore.SetCurrentAnimation(SPAWN)
		if !ok {
			fmt.Println("Warning: Unable to start the spawning animation")
		}
	}
	baseEntity := entity.NewEntity(pos.X, pos.Y)
	return &PotatoEnemy{
		Enemy: Enemy{
			Entity:         *baseEntity,
			Speed:          config.POTATO_SPEED,
			Health:         component.NewHealth(config.POTATO_HEALTH),
			Damage:         config.POTATO_DAMAGE,
			AttackCooldown: config.POTATO_ATTACK_COOLDOWN,
			attackTimer:    config.POTATO_ATTACK_START,
			animationStore: animationStore,
			DropProb:       config.POTATO_DROP_PROB,   // 80% chance to drop an item
			DropAmount:     config.POTATO_DROP_AMOUNT, // Drops 1 item
		},
		MeleeEnemyData: MeleeEnemyData{
			AttackRange: config.POTATO_ATTACK_RANGE,
		},
	}
}

func (e *PotatoEnemy) SetWalkingAnimation(player *player.Player) {
	dir := player.Pos.Sub(e.Pos)
	if dir.X > 0 {
		if ok := e.animationStore.SetCurrentAnimation(WALK_RIGHT); !ok {
			fmt.Println("Warning: Unable to start the potato walkRight animation")
		}
	} else {
		if ok := e.animationStore.SetCurrentAnimation(WALK_LEFT); !ok {
			fmt.Println("Warning: Unable to start the potato walkLeft animation")
		}
	}
}

func (e *PotatoEnemy) SetAttackAnimation(player *player.Player) {
	dir := player.Pos.Sub(e.Pos)
	if dir.X > 0 {
		if ok := e.animationStore.SetCurrentAnimation(ATTACK_RIGHT); !ok {
			fmt.Println("Warning: Unable to start the potato attackRight animation")
		}
	} else {
		if ok := e.animationStore.SetCurrentAnimation(ATTACK_LEFT); !ok {
			fmt.Println("Warning: Unable to start the potato attackLeft animation")
		}
	}
}

func (e *PotatoEnemy) Update(player *player.Player, dt float64) {
	e.animationStore.Update()

	// Only update the potato if it is alive
	if e.Health.HP > 0 {
		// Set current animation to walking if no animation is currently running
		// Or Update the type of running animation if currently the animation is running
		animationName := e.animationStore.GetCurrentAnimationName()
		updateRunningType := animationName == WALK_LEFT || animationName == WALK_RIGHT
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
				e.SetAttackAnimation(player)
			}
		}
	} else {
		// The potato is dead
		// So we start the death animation
		e.animationStore.SetCurrentAnimation(DEATH)
		e.UpdateKnockback()
	}
}

func (e *PotatoEnemy) Draw(screen *ebiten.Image, camX, camY float64) {
	frameImage := e.animationStore.GetImage()
	if frameImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(e.Pos.X-camX-16.0, e.Pos.Y-camY-16.0)
		screen.DrawImage(frameImage, op)
	} else {
		e.DefaultDraw(screen, camX, camY, config.POTATO_WIDTH, config.POTATO_HEIGHT,
			color.RGBA{R: config.POTATO_COLOR_R, G: config.POTATO_COLOR_G, B: config.POTATO_COLOR_B, A: 255})
	}
}

func (e *PotatoEnemy) IsAlive() bool {
	if e.Enemy.IsAlive() {
		return true
	} else {
		// The potato will be marked as dead after the death animation is finished
		return !e.animationStore.GetCurrentAnimation().IsFinished()
	}
}

func (e *PotatoEnemy) GetPosition() component.Vector2D {
	return e.Enemy.GetPosition()
}

func (e *PotatoEnemy) TryDrop(elapsedMinutes float32) []item.Item {
	prob := e.DropProb + elapsedMinutes*0.001 // +0.1% per minute, // so +1% per 10 minutes
	if prob > 1 {
		prob = 1
	}
	// Basis + rate * time
	amount := e.DropAmount + int(elapsedMinutes*config.POTATO_DROP_AMOUNT_PER_MINUTE)

	fmt.Printf("Potato Enemy TryDrop: prob=%.2f, amount=%d, %d min\n", prob, amount, int(elapsedMinutes))

	var drops []item.Item
	if rand.Float32() < prob {
		for i := 0; i < amount; i++ {
			drops = append(drops, *item.NewPotato(e.Pos.X, e.Pos.Y))
		}
	}
	return drops
}
