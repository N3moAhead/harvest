package enemy

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/entity/item"
)

type CabbageEnemy struct {
	BaseMeleeEnemy
}

func NewCabbageEnemy(pos component.Vector2D) *CabbageEnemy {
	cabbageSprite, ok := assets.AssetStore.GetImage("cabbage")
	store := animation.NewAnimationStore()
	if ok {
		walkRightAnimation, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 1*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_RIGHT, walkRightAnimation)
		} else {
			fmt.Println("Does not work here....")
			panic(err)
		}
		walkLeftAnimation, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 0*32, 8, 6, true)
		if err == nil {
			store.AddAnimation(WALK_LEFT, walkLeftAnimation)
		}
		attackRight, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 3*32, 8, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_RIGHT, attackRight)
		}
		attackLeft, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 4*32, 8, 6, false)
		if err == nil {
			store.AddAnimation(ATTACK_LEFT, attackLeft)
		}
		spawnAnimation, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 2*32, 8, 10, false)
		if err == nil {
			store.AddAnimation(SPAWN, spawnAnimation)
		}
		deathAnimation, err := animation.NewAnimation(cabbageSprite, 32, 32, 0, 6*32, 4, 10, false)
		if err == nil {
			store.AddAnimation(DEATH, deathAnimation)
		}
	}
	return &CabbageEnemy{
		BaseMeleeEnemy: *NewBaseMeleeEnemy(TypeCabbage, pos, store, &BaseMeleeOptions{
			Speed:               config.CABBAGE_SPEED,
			MaxHealth:           config.CABBAGE_HEALTH,
			Damage:              config.CABBAGE_DAMAGE,
			AttackCooldown:      config.CABBAGE_ATTACK_COOLDOWN,
			DropProb:            config.CABBAGE_DROP_PROB,
			DropAmount:          config.CABBAGE_DROP_AMOUNT,
			DropAmountPerMinute: config.CABBAGE_DROP_AMOUNT_PER_MINUTE,
			AttackRange:         config.CABBAGE_ATTACK_RANGE,
			SpawnItem:           item.NewCabbage,
		}),
	}
}

var _ EnemyInterface = (*OnionEnemy)(nil)
