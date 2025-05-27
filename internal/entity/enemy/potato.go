package enemy

import (
	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
)

type PotatoEnemy struct {
	BaseMeleeEnemy
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
	}
	return &PotatoEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypePotato, pos, animationStore, &BaseMeleeOptions{
			Speed:               config.POTATO_SPEED,
			MaxHealth:           config.POTATO_HEALTH,
			Damage:              config.POTATO_DAMAGE,
			AttackRange:         config.POTATO_ATTACK_RANGE,
			AttackCooldown:      config.POTATO_ATTACK_COOLDOWN,
			DropProb:            config.POTATO_DROP_PROB,
			DropAmount:          config.POTATO_DROP_AMOUNT,
			DropAmountPerMinute: config.POTATO_DROP_AMOUNT_PER_MINUTE,
		}),
	}
}
