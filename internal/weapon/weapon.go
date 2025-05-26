package weapon

import (
	"fmt"
	"time"

	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

type Weapon interface {
	Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) // Update the weopon
	Draw(
		screen *ebiten.Image,
		player *player.Player,
		mapOffsetX float64,
		mapOffsetY float64,
	) // Draws the weopon around the player if needed
	// Functions used for the ui hud
	GetType() itemtype.ItemType
	Description() string
	Name() string
	Level() int
	MaxLevel() int
}

// This struct can be used for multiple weopons
// for usage we just have to define the properties we need
// so for melee weopons theres no need to define ProjectileSpeed and stuff...
type WeaponStats struct {
	Damage          float64
	Cooldown        time.Duration
	AreaSize        float64
	ProjectileSpeed float64
	Knockback       float64
	ProjectileCount int
	Pierce          int
	Duration        time.Duration
}

// Can be embedded to get basic attributes in a weapon
// it already implements most of the functions needed for the
// weapon interface
type BaseWeapon struct {
	name          string
	description   string
	level         int
	maxLevel      int
	statsPerLevel []WeaponStats // level - 1  => Current Weapon Stats
	cooldownTimer time.Duration
	itemType      itemtype.ItemType
}

func (b *BaseWeapon) Name() string {
	return b.name
}

func (b *BaseWeapon) Level() int {
	return b.level
}

func (b *BaseWeapon) MaxLevel() int {
	return b.maxLevel
}

func (b *BaseWeapon) Description() string {
	return b.description
}

func (b *BaseWeapon) BaseStats() WeaponStats {
	if b.level > 0 && b.level <= len(b.statsPerLevel) {
		return b.statsPerLevel[b.level-1]
	}
	// Fallback if the level is out of the available options
	fmt.Printf("Warning: No base stats for weapon '%s' at level %d defined.\n", b.name, b.level)
	if len(b.statsPerLevel) > 0 {
		return b.statsPerLevel[0]
	}
	// We could also consider to panic here
	return WeaponStats{}
}

func (b *BaseWeapon) CurrentStats(player *player.Player) WeaponStats {
	stats := b.BaseStats()
	// Can be used for modifications by buffs in the player or something
	// A dmg buff or an area attack buff
	return stats
}

// Reduces the timer and returns true if the cooldown is finished.
func (b *BaseWeapon) UpdateCooldown(dt time.Duration) bool {
	if b.cooldownTimer > 0 {
		b.cooldownTimer -= dt
	}
	return b.cooldownTimer <= 0
}

func (b *BaseWeapon) GetType() itemtype.ItemType {
	return b.itemType
}

// Sets the timer based on the current stats.
func (b *BaseWeapon) ResetCooldown(player *player.Player) {
	b.cooldownTimer = b.CurrentStats(player).Cooldown
}
