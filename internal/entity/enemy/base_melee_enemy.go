package enemy

import (
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/entity/item"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type BaseMeleeEnemy struct {
	Enemy
	AttackRange      float64
	spawnItem        func(x, y float64) *item.Item
	damageIndicators []*DamageIndicator
}

type BaseMeleeOptions struct {
	Speed               float64
	MaxHealth           float64
	Damage              float64
	AttackCooldown      float64
	DropProb            float32
	DropAmount          int
	DropAmountPerMinute float32
	AttackRange         float64
	SpawnItem           func(x, y float64) *item.Item
}

func NewBaseMeleeEnemy(enemyType EnemyType, pos component.Vector2D, store *animation.AnimationStore, op *BaseMeleeOptions) *BaseMeleeEnemy {
	// Check if all needed animations exist
	checkMeleeAnimations(store)
	ok := store.SetCurrentAnimation(SPAWN)
	if !ok {
		fmt.Println("Warning: Unable to start the spawning animation")
	}

	return &BaseMeleeEnemy{
		Enemy: Enemy{
			Entity:              *entity.NewEntity(pos.X, pos.Y),
			Speed:               op.Speed,
			Health:              component.NewHealth(op.MaxHealth),
			Damage:              op.Damage,
			AttackCooldown:      op.AttackCooldown,
			DropProb:            op.DropProb,
			DropAmount:          op.DropAmount,
			DropAmountPerMinute: op.DropAmountPerMinute,
			attackTimer:         0.0,
			animationStore:      store,
			enemyType:           enemyType,
		},
		AttackRange:      op.AttackRange,
		spawnItem:        op.SpawnItem,
		damageIndicators: make([]*DamageIndicator, 0),
	}
}

func (e *BaseMeleeEnemy) Update(player *player.Player, dt float64) {
	e.animationStore.Update()

	/// Remove dead damageIndicators
	n := 0
	for i, indicator := range e.damageIndicators {
		isAlive := indicator.Update(e.GetPosition())
		if isAlive {
			if n != i {
				e.damageIndicators[n] = indicator
			}
			n++
		}
	}
	e.damageIndicators = e.damageIndicators[:n]

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
		// The enemy is dead
		// So we start the death animation
		e.animationStore.SetCurrentAnimation(DEATH)
		// We still want to see the knockback happening
		// It just feels way better when playing :)
		e.UpdateKnockback()
	}
}

func (e *BaseMeleeEnemy) Draw(screen *ebiten.Image, camX, camY float64) {
	frameImage := e.animationStore.GetImage()
	assetSize := config.DEFAULT_ENEMY_ASSET_SIZE
	assetSizeHalf := assetSize / 2
	if frameImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(e.Pos.X-camX-assetSizeHalf, e.Pos.Y-camY-assetSizeHalf)
		screen.DrawImage(frameImage, op)
	} else {
		e.DefaultDraw(
			screen,
			camX,
			camY,
			int(assetSize),
			int(assetSize),
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
		)
	}
	for _, dmgIndicator := range e.damageIndicators {
		dmgIndicator.Draw(screen, camX, camY)
	}
}

func (e *BaseMeleeEnemy) TakeDamage(damage float64) {
	// Spawn a new damage indicator
	newDmgIndicator := NewDamageIndicator(e.GetPosition(), component.NewVector2D(0, -1), damage)
	e.damageIndicators = append(e.damageIndicators, newDmgIndicator)
	e.Health.Damage(damage)
}

func (e *BaseMeleeEnemy) TryDrop(elapsedMinutes float32) []item.Item {
	prob := e.DropProb + elapsedMinutes*0.001 // +0.1% per minute, // so +1% per 10 minutes
	if prob > 1 {
		prob = 1
	}
	// Basis + rate * time
	amount := e.DropAmount + int(elapsedMinutes*e.DropAmountPerMinute)

	var drops []item.Item
	if rand.Float32() < prob {
		for i := 0; i < amount; i++ {
			drops = append(drops, *e.spawnItem(e.Pos.X, e.Pos.Y))
		}
	}
	return drops
}

func (e *BaseMeleeEnemy) SetWalkingAnimation(player *player.Player) {
	dir := player.GetPosition().Sub(e.Pos)
	if dir.X > 0 {
		if ok := e.animationStore.SetCurrentAnimation(WALK_RIGHT); !ok {
			fmt.Println("Warning: Unable to start walkRight animation")
		}
	} else {
		if ok := e.animationStore.SetCurrentAnimation(WALK_LEFT); !ok {
			fmt.Println("Warning: Unable to start the walkLeft animation")
		}
	}
}

func (e *BaseMeleeEnemy) SetAttackAnimation(player *player.Player) {
	dir := player.Pos.Sub(e.Pos)
	if dir.X > 0 {
		if ok := e.animationStore.SetCurrentAnimation(ATTACK_RIGHT); !ok {
			fmt.Println("Warning: Unable to start the attackRight animation")
		}
	} else {
		if ok := e.animationStore.SetCurrentAnimation(ATTACK_LEFT); !ok {
			fmt.Println("Warning: Unable to start the attackLeft animation")
		}
	}
}

// Mark enemies as dead after the death animation is finished
func (e *BaseMeleeEnemy) IsAlive() bool {
	if e.Enemy.IsAlive() {
		return true
	} else {
		return !e.animationStore.GetCurrentAnimation().IsFinished()
	}
}

var _ EnemyInterface = (*BaseMeleeEnemy)(nil)

func doesAnimationExist(store *animation.AnimationStore, name string) {
	_, ok := store.GetAnimation(name)
	if !ok {
		panic("Warning: Missing Animation: " + name)
	}
}

func checkMeleeAnimations(store *animation.AnimationStore) {
	animNames := []string{
		WALK_RIGHT,
		WALK_LEFT,
		ATTACK_RIGHT,
		ATTACK_LEFT,
		SPAWN,
		DEATH,
	}
	for _, name := range animNames {
		doesAnimationExist(store, name)
	}
}
