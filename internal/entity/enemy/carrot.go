package enemy

import (
	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
)

type CarrotEnemy struct {
	BaseMeleeEnemy
}

func NewCarrotEnemy(pos component.Vector2D) *CarrotEnemy {
	carrotSprite, ok := assets.AssetStore.GetImage("carrot")
	animationStore := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 32, 8, 6, true)
		if err == nil {
			animationStore.AddAnimation(WALK_RIGHT, walkRightAnimation)
		}
		walkLeftAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 6*32, 8, 6, true)
		if err == nil {
			animationStore.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 4*32, 2, 6, false)
		if err == nil {
			animationStore.AddAnimation(ATTACK_LEFT, attackAnimation)
			animationStore.AddAnimation(ATTACK_RIGHT, attackAnimation)
		}
		spawnAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 2*32, 6, 10, false)
		if err == nil {
			animationStore.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(carrotSprite, 32, 32, 0, 7*32, 6, 10, false)
		if err == nil {
			animationStore.AddAnimation(DEATH, deathAnimation)
		}
	}
	return &CarrotEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypeCarrot, pos, animationStore, &BaseMeleeOptions{
			Speed:          config.CARROT_SPEED,
			MaxHealth:      config.CARROT_HEALTH,
			Damage:         config.CARROT_DAMAGE,
			AttackCooldown: config.CARROT_ATTACK_COOLDOWN,
			DropProb:       config.CARROT_DROP_PROB,
			DropAmount:     config.CARROT_DROP_AMOUNT,
			AttackRange:    config.CARROT_ATTACK_RANGE,
			SpawnItem:      item.NewCarrot,
		}),
	}
}
