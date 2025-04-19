package weapon

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/collision"
	"github.com/N3moAhead/harvest/internal/enemy"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
)

var whitePixelImage *ebiten.Image

const (
	baseSpoonRadius       = 100.0                  // Basic range of the hit in pixels
	spoonAttackAngle      = math.Pi                // 180 degrees
	attackDisplayDuration = 100 * time.Millisecond // How long the animation will be displayed
)

func init() {
	if whitePixelImage == nil {
		whitePixelImage = ebiten.NewImage(1, 1)
		whitePixelImage.Fill(color.White)
	}
}

type Spoon struct {
	BaseWeapon
	attackDisplayTimer time.Duration
}

func NewSpoon(assetStore *assets.Store) *Spoon {
	stats := []WeaponStats{
		// Level 1
		{Damage: 5, Cooldown: 700 * time.Millisecond, AreaSize: 1.0, Pierce: 3},
		// Level 2
		{Damage: 15, Cooldown: 650 * time.Millisecond, AreaSize: 1.1, Pierce: 4},
		// Level 3
		{Damage: 25, Cooldown: 650 * time.Millisecond, AreaSize: 1.1, Pierce: 5},
	}

	spoonIcon, ok := assetStore.GetImage("spoon_icon")
	if !ok {
		fmt.Println("Warning: Spoon icon not found")
	}
	return &Spoon{
		BaseWeapon: BaseWeapon{
			name:          "Spoon",
			description:   "A mighty spoon! there are sayings Gordan Ramsey touched it once...",
			cooldownTimer: 0,
			level:         1,
			maxLevel:      len(stats),
			statsPerLevel: stats,
			icon:          spoonIcon,
		},
		attackDisplayTimer: 0,
	}
}

func (s *Spoon) Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) {
	// Update the cooldown
	canAttack := s.UpdateCooldown(dt)

	// Update animation timer
	if s.attackDisplayTimer > 0 {
		s.attackDisplayTimer -= dt
	}

	if canAttack {
		s.ResetCooldown(player) // Reset the cooldown before the attack

		stats := s.CurrentStats(player)
		playerPos := player.Pos
		facingDir := player.GetFacingDirection()

		// Calculate the spoon range
		currentRadius := baseSpoonRadius * stats.AreaSize

		// Find the enemies in spoon hit range
		hitEnemies := collision.FindEnemiesInArc(
			playerPos,
			currentRadius,
			facingDir,
			spoonAttackAngle,
			enemies,
		)

		// Lets deal some dmg but we have to keep track of the pierce value
		hits := 0
		for _, enemy := range hitEnemies {
			if hits >= stats.Pierce && stats.Pierce > 0 {
				break
			}
			enemy.TakeDamage(stats.Damage)
			// TODO play spoon_hit sound
		}

		if hits > 0 {
			// TODO remove just debugging
			fmt.Printf("Spoon did hit %d enemies. \n", hits)
		}

		// Start the animation
		s.attackDisplayTimer = attackDisplayDuration
		// TODO Play spoon swing sound
	}
}

func (s *Spoon) Draw(screen *ebiten.Image, owner *player.Player, mapOffsetX float64, mapOffsetY float64) {
	// TODO implement the logic for drawing the swinging animation
	// currently its invisible
}
