package enemy

import (
	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
)

type OnionEnemy struct {
	BaseMeleeEnemy
}

func NewOnionEnemy(pos component.Vector2D) *OnionEnemy {
	onionSprite, ok := assets.AssetStore.GetImage("onion")
	store := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(onionSprite, 32, 32, 0, 1*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_RIGHT, walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(onionSprite, 32, 32, 0, 0*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackAnimation, err := animation.NewAnimation(onionSprite, 32, 32, 0, 3*32, 8, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_LEFT, attackAnimation)
			store.AddAnimation(ATTACK_RIGHT, attackAnimation)
		}
		spawnAnimation, err := animation.NewAnimation(onionSprite, 32, 32, 0, 2*32, 8, 10, false)
		if err == nil {
			store.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(onionSprite, 32, 32, 0, 5*32, 4, 10, false)
		if err == nil {
			store.AddAnimation(DEATH, deathAnimation)
		}
	}
	return &OnionEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypeOnion, pos, store, &BaseMeleeOptions{
			Speed:               config.ONION_SPEED,
			MaxHealth:           config.ONION_HEALTH,
			Damage:              config.CARROT_DAMAGE,
			AttackCooldown:      config.ONION_ATTACK_COOLDOWN,
			DropProb:            config.ONION_DROP_PROB,
			DropAmount:          config.ONION_DROP_AMOUNT,
			DropAmountPerMinute: config.ONION_DROP_AMOUNT_PER_MINUTE,
			AttackRange:         config.ONION_ATTACK_RANGE,
			SpawnItem:           item.NewOnion,
		}),
	}
}

var _ EnemyInterface = (*OnionEnemy)(nil)
