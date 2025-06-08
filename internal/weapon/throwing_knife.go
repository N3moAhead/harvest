package weapon

import (
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/N3moAhead/harvest/internal/entity/projectile"
	"github.com/hajimehoshi/ebiten/v2"
)

type ThrowingKnife struct {
	RangeBaseWeapon
	// Maybe move projectiles out of the weapon
	knifes []*projectile.KnifeProjectile
}

func NewThrowingKnife() *ThrowingKnife {
	stats := []RangeWeaponStats{
		// Level 1
		{
			ProjectileAmount: 1,
			Cooldown:         2 * time.Second,
			Speed:            5,
			Damage:           2,
			HitRadius:        16, // Its the px size of the knife image
			Pierce:           5,
			Duration:         2 * time.Second,
			Knockback:        50,
			BulletSpread:     30.0,
		},
		// Level 2
		{
			ProjectileAmount: 3,
			Cooldown:         1500 * time.Millisecond,
			Speed:            5,
			Damage:           2,
			HitRadius:        16, // Its the px size of the knife image
			Pierce:           7,
			Duration:         4 * time.Second,
			Knockback:        50,
			BulletSpread:     30.0,
		},
		// Level 3
		{
			ProjectileAmount: 5,
			Cooldown:         1000 * time.Millisecond,
			Speed:            5,
			Damage:           2,
			HitRadius:        16, // Its the px size of the knife image
			Pierce:           10,
			Duration:         5 * time.Second,
			Knockback:        60,
			BulletSpread:     30.0,
		},
	}
	return &ThrowingKnife{
		RangeBaseWeapon: RangeBaseWeapon{
			name:          "Throwing Knifes",
			description:   "Let sharp kitchen knifes rain down upon some veggies",
			level:         1,
			maxLevel:      len(stats),
			statsPerLevel: stats,
			cooldownTimer: 0,
			itemType:      itemtype.ThrowingKnifes,
		},
	}
}

func (t *ThrowingKnife) Draw(
	screen *ebiten.Image,
	player *player.Player,
	mapOffsetX, mapOffsetY float64,
) {
	for _, knife := range t.knifes {
		knife.Draw(screen, mapOffsetX, mapOffsetY)
	}
}

func (t *ThrowingKnife) Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) {
	// Update all throwing knifes
	n := 0
	for i, knife := range t.knifes {
		active := knife.Update(enemies)
		if active {
			if n != i {
				t.knifes[n] = knife
			}
			n++
		}
	}
	t.knifes = t.knifes[:n] // Remove all not longer active knifes

	// Spawn new knifes
	canAttack := t.UpdateCooldown(dt)

	if canAttack {
		t.ResetCooldown(player)

		stats := t.CurrentStats(player)
		throwingDirections := calculateSpreadDirections(player.GetFacingDirection(), stats.ProjectileAmount, stats.BulletSpread)

		for _, dir := range throwingDirections {
			t.knifes = append(t.knifes, projectile.NewKnifeProjectile(
				player.GetPosition(),
				dir,
				stats.Speed,
				stats.Duration,
				stats.Damage,
				stats.HitRadius,
				stats.Pierce,
				stats.Knockback,
			))
		}

		knifeThrowSfx, ok := assets.AssetStore.GetSFXData("knife_throw")
		if ok {
			sfxPlayer := assets.AudioContext.NewPlayerFromBytes(knifeThrowSfx)
			sfxPlayer.Play()
		}
	}
}

var _ Weapon = (*ThrowingKnife)(nil)
