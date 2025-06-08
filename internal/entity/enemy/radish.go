package enemy

import (
	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
)

type RadishEnemy struct {
	BaseMeleeEnemy
}

func NewRadishEnemy(pos component.Vector2D) *RadishEnemy {
	radishSprite, ok := assets.AssetStore.GetImage("radish")
	store := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(radishSprite, 32, 32, 0, 1*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_RIGHT, walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(radishSprite, 32, 32, 0, 0*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackRight, err := animation.NewAnimation(radishSprite, 32, 32, 0, 3*32, 8, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_RIGHT, attackRight)
		}
		attackLeft, err := animation.NewAnimation(radishSprite, 32, 32, 0, 4*32, 8, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_LEFT, attackLeft)
		}
		spawnAnimation, err := animation.NewAnimation(radishSprite, 32, 32, 0, 2*32, 8, 10, false)
		if err == nil {
			store.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(radishSprite, 32, 32, 0, 10*32, 4, 10, false)
		if err == nil {
			store.AddAnimation(DEATH, deathAnimation)
		}
	}
	return &RadishEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypeRadish, pos, store, &BaseMeleeOptions{
			Speed:               config.RADISH_SPEED,
			MaxHealth:           config.RADISH_HEALTH,
			Damage:              config.RADISH_DAMAGE,
			AttackCooldown:      config.RADISH_ATTACK_COOLDOWN,
			DropProb:            config.RADISH_DROP_PROB,
			DropAmount:          config.RADISH_DROP_AMOUNT,
			DropAmountPerMinute: config.RADISH_DROP_AMOUNT_PER_MINUTE,
			AttackRange:         config.RADISH_ATTACK_RANGE,
			SpawnItem:           item.NewRadish,
		}),
	}
}

var _ EnemyInterface = (*RadishEnemy)(nil)
