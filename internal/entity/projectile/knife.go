package projectile

import (
	"fmt"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
)

// This projectile will be used by the ThrowingKnifes weapon

type KnifeProjectile struct {
	BaseProjectile
}

func NewKnifeProjectile(
	pos component.Vector2D,
	dir component.Vector2D,
	speed float64,
	duration time.Duration,
	dmg float64,
	hitRadius float64,
	pierce int,
	knockback float64,
) *KnifeProjectile {
	knifeImg, ok := assets.AssetStore.GetImage("knife_projectile")
	if !ok {
		fmt.Println("Warning: Could not load knifeImg in NewKnifeProjectile")
	}
	return &KnifeProjectile{
		BaseProjectile: *NewBaseProjectile(
			pos,
			dir,
			speed,
			knifeImg,
			duration,
			dmg,
			hitRadius,
			pierce,
			knockback,
			"knife_throw_impact",
		),
	}
}

func (k *KnifeProjectile) PlayImpactSound() {

}

var _ Projectile = (*KnifeProjectile)(nil)
