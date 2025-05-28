package enemy

import (
	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
)

type LeekEnemy struct {
	BaseMeleeEnemy
}

func NewLeekEnemy(pos component.Vector2D) *OnionEnemy {
	leekSprite, ok := assets.AssetStore.GetImage("leek")
	store := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(leekSprite, 32, 32, 0, 1*32, 7, 6, true)
		if err == nil {
			store.AddAnimation(WALK_RIGHT, walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(leekSprite, 32, 32, 0, 0*32, 7, 6, true)
		if err == nil {
			store.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackRight, err := animation.NewAnimation(leekSprite, 32, 32, 0, 3*32, 7, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_RIGHT, attackRight)
		}
		attackLeft, err := animation.NewAnimation(leekSprite, 32, 32, 0, 4*32, 7, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_LEFT, attackLeft)
		}
		spawnAnimation, err := animation.NewAnimation(leekSprite, 32, 32, 0, 2*32, 7, 10, false)
		if err == nil {
			store.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(leekSprite, 32, 32, 0, 5*32, 7, 10, false)
		if err == nil {
			store.AddAnimation(DEATH, deathAnimation)
		}
	}
	return &OnionEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypeLeek, pos, store, &BaseMeleeOptions{
			Speed:               config.LEEK_SPEED,
			MaxHealth:           config.LEEK_HEALTH,
			Damage:              config.LEEK_DAMAGE,
			AttackCooldown:      config.LEEK_ATTACK_COOLDOWN,
			DropProb:            config.LEEK_DROP_PROB,
			DropAmount:          config.LEEK_DROP_AMOUNT,
			DropAmountPerMinute: config.LEEK_DROP_AMOUNT_PER_MINUTE,
			AttackRange:         config.LEEK_ATTACK_RANGE,
			SpawnItem:           item.NewLeek,
		}),
	}
}

var _ EnemyInterface = (*OnionEnemy)(nil)
