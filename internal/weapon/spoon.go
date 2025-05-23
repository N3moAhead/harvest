package weapon

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/collision"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	baseSpoonRadius  = 100.0   // Basic range of the hit in pixels
	spoonAttackAngle = math.Pi // 180 degrees
	frameWidth       = 32      // The height of each slash animation frame
	frameHeight      = 64      // The width of each slash animation frame
	frameCount       = 4       // The amount of existing frames
	animationSpeed   = 6       // Ticks for each frame
)

type Spoon struct {
	BaseWeapon
	slashImage   *ebiten.Image
	slashSound   []byte
	frameTimer   int
	currentFrame int
	displaySlash bool
	hitDirection component.Vector2D
}

func NewSpoon() *Spoon {
	stats := []WeaponStats{
		// Level 1
		{
			Damage:    1,
			Cooldown:  2 * time.Second,
			AreaSize:  1.0,
			Pierce:    1000,
			Knockback: 300.0,
		},
		// Level 2
		{
			Damage:    2,
			Cooldown:  1500 * time.Millisecond,
			AreaSize:  1.1,
			Pierce:    12,
			Knockback: 25.0,
		},
		// Level 3
		{
			Damage:    5,
			Cooldown:  650 * time.Millisecond,
			AreaSize:  1.2,
			Pierce:    15,
			Knockback: 30.0,
		},
	}

	spoonIcon, ok := assets.AssetStore.GetImage("spoon_icon")
	if !ok {
		fmt.Println("Warning: Spoon icon not found")
	}
	slashImage, ok := assets.AssetStore.GetImage("spoon_slash")
	if !ok {
		fmt.Println("Warning: Spoon slash image not found")
	}
	slashSound, ok := assets.AssetStore.GetSFXData("spoon_slash")
	if !ok {
		fmt.Println("Warning: Spoon slash sound not found")
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
		slashImage: slashImage,
		slashSound: slashSound,
	}
}

func (s *Spoon) Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) {
	// Update the cooldown
	canAttack := s.UpdateCooldown(dt)

	// Update the animation frames
	s.frameTimer++
	if s.frameTimer >= animationSpeed {
		s.frameTimer = 0 // Reset the frame timer
		if s.currentFrame == frameCount-1 {
			s.displaySlash = false
		} else {
			s.currentFrame = (s.currentFrame + 1) % frameCount
		}
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
			enemy.AddKnockback(&player.Pos, stats.Knockback)
			hits++
			// TODO play spoon_hit sound
		}

		// Start the animation
		s.displaySlash = true
		s.currentFrame = 0
		s.frameTimer = 0
		s.hitDirection = player.GetFacingDirection()

		// Play the sound
		sfxPlayer := assets.AudioContext.NewPlayerFromBytes(s.slashSound)
		sfxPlayer.Play()
	}
}

func (s *Spoon) Draw(screen *ebiten.Image, player *player.Player, mapOffsetX float64, mapOffsetY float64) {
	if s.slashImage == nil {
		return
	}

	if s.displaySlash {
		sx := s.currentFrame * frameWidth
		sy := 0
		frameRect := image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)

		bounds := s.slashImage.Bounds()
		if !frameRect.In(bounds) {
			fmt.Printf("Error: Frame rect %v out of bounds %v\n", frameRect, bounds)
			return
		}
		frameImage := s.slashImage.SubImage(frameRect).(*ebiten.Image)

		pivotX := float64(frameWidth) / 2.0
		pivotY := float64(frameHeight) / 2.0
		angle := math.Atan2(s.hitDirection.Y, s.hitDirection.X)

		stats := s.CurrentStats(player)
		currentRadius := baseSpoonRadius * stats.AreaSize
		calculatedScale := currentRadius / float64(frameWidth)

		finalScreen := component.NewVector2D(player.Pos.X-mapOffsetX, player.Pos.Y-mapOffsetY)
		// Moving the animation to the outside of the player
		finalScreen = finalScreen.Add(s.hitDirection.Mul(frameWidth / 2)) // TODO i dont know if frameWidth / 2 is the correct value

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-pivotX, -pivotY)
		if calculatedScale != 1.0 {
			op.GeoM.Scale(calculatedScale, calculatedScale)
		}
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(finalScreen.X, finalScreen.Y)

		screen.DrawImage(frameImage, op)
	}
}
