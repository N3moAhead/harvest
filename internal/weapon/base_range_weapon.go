package weapon

import (
	"fmt"
	"time"

	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type RangeWeaponStats struct {
	// Weapon Settings
	ProjectileAmount int
	Cooldown         time.Duration
	BulletSpread     float64
	// Projectile Settings
	Speed     float64
	Damage    float64
	HitRadius float64
	Pierce    int
	Duration  time.Duration
	Knockback float64
}

type RangeBaseWeapon struct {
	name          string
	description   string
	level         int
	maxLevel      int
	statsPerLevel []RangeWeaponStats
	cooldownTimer time.Duration
	itemType      itemtype.ItemType
}

func (b *RangeBaseWeapon) Name() string {
	return b.name
}

func (b *RangeBaseWeapon) Level() int {
	return b.level
}

func (b *RangeBaseWeapon) MaxLevel() int {
	return b.maxLevel
}

func (b *RangeBaseWeapon) Description() string {
	return b.description
}

func (b *RangeBaseWeapon) BaseStats() RangeWeaponStats {
	if b.level > 0 && b.level <= len(b.statsPerLevel) {
		return b.statsPerLevel[b.level-1]
	}
	// Fallback if the level is out of the available options
	fmt.Printf("Warning: No base stats for weapon '%s' at level %d defined.\n", b.name, b.level)
	if len(b.statsPerLevel) > 0 {
		return b.statsPerLevel[0]
	}
	// We could also consider to panic here
	return RangeWeaponStats{}
}

func (b *RangeBaseWeapon) CurrentStats(player *player.Player) RangeWeaponStats {
	stats := b.BaseStats()
	// Can be used for modifications by buffs in the player or something
	// A dmg buff or an area attack buff
	return stats
}

// Reduces the timer and returns true if the cooldown is finished.
func (b *RangeBaseWeapon) UpdateCooldown(dt time.Duration) bool {
	if b.cooldownTimer > 0 {
		b.cooldownTimer -= dt
	}
	return b.cooldownTimer <= 0
}

func (b *RangeBaseWeapon) GetType() itemtype.ItemType {
	return b.itemType
}

// Sets the timer based on the current stats.
func (b *RangeBaseWeapon) ResetCooldown(player *player.Player) {
	b.cooldownTimer = b.CurrentStats(player).Cooldown
}

func (b *RangeBaseWeapon) Draw(
	screen *ebiten.Image,
	player *player.Player,
	mapOffsetX float64,
	mapOffsetY float64,
) {
	fmt.Println("Warning: Draw is not implemented in ", b.GetType().String())
}

func (b *RangeBaseWeapon) Update(
	player *player.Player,
	enemies []enemy.EnemyInterface,
	dt time.Duration,
) {
	fmt.Println("Warning: Update is not implemented in ", b.GetType().String())
}

var _ Weapon = (*RangeBaseWeapon)(nil)
