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
	baseRollingPinRadius  = 100.0   // Basic range of the hit in pixels
	rollingPinAttackAngle = math.Pi // 180 degrees
)

// OPTIONAL: Maybe instead of a copy of the spoon slash -> the rolling pin could roll in the direction
// the character is facing and then roll back after reaching a certain point?

type RollingPin struct {
	BaseWeapon
	rollImage    *ebiten.Image
	rollSound    []byte
	frameTimer   int
	currentFrame int
	displayRoll  bool
	hitDirection component.Vector2D
}

func NewRollingPin() *RollingPin {
	stats := []WeaponStats{
		// Level 1
		{
			Damage:    2,
			Cooldown:  2 * time.Second,
			AreaSize:  1.0,
			Pierce:    5,
			Knockback: 300.0,
		},
		// Level 2
		{
			Damage:    4,
			Cooldown:  1 * time.Second,
			AreaSize:  1.1,
			Pierce:    7,
			Knockback: 400.0,
		},
		// Level 3
		{
			Damage:    8,
			Cooldown:  1500 * time.Millisecond,
			AreaSize:  1.2,
			Pierce:    10,
			Knockback: 550.0,
		},
	}

	rollingPinIcon, ok := assets.AssetStore.GetImage("rolling_pin_icon")
	if !ok {
		fmt.Println("Warning: Rolling pin icon not found")
	}
	rollImage, ok := assets.AssetStore.GetImage("rolling_pin_roll")
	if !ok {
		fmt.Println("Warning: Rolling pin roll image not found")
	}
	rollSound, ok := assets.AssetStore.GetSFXData("rolling_pin_roll")
	if !ok {
		fmt.Println("Warning: Rolling pin roll sound not found")
	}
	return &RollingPin{
		BaseWeapon: BaseWeapon{
			name:          "Rolling Pin",
			description:   "Crushing your enemies with a mighty roll./Beating your enemies just like mama used to do with you.",
			cooldownTimer: 0,
			level:         1,
			maxLevel:      len(stats),
			statsPerLevel: stats,
			icon:          rollingPinIcon,
		},
		rollImage: rollImage,
		rollSound: rollSound,
	}
}

func (rp *RollingPin) Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) {
	// Update the cooldown
	canAttack := rp.UpdateCooldown(dt)

	// Update the animation frames
	rp.frameTimer++
	if rp.frameTimer >= animationSpeed {
		rp.frameTimer = 0 // Reset the frame timer
		if rp.currentFrame == frameCount-1 {
			rp.displayRoll = false
		} else {
			rp.currentFrame = (rp.currentFrame + 1) % frameCount
		}
	}

	if canAttack {
		rp.ResetCooldown(player) // Reset the cooldown before the attack

		stats := rp.CurrentStats(player)
		playerPos := player.Pos
		facingDir := player.GetFacingDirection()

		// Calculate the rolling pin range
		currentRadius := baseRollingPinRadius * stats.AreaSize

		// Find the enemies in rolling pin hit range
		hitEnemies := collision.FindEnemiesInArc(
			playerPos,
			currentRadius,
			facingDir,
			rollingPinAttackAngle,
			enemies,
		)

		// Lets deal some dmg but we have to keep track of the pierce value
		// -> TODO exchange pierce with something that makes sense (e.g. crush value or squeeze value)
		hits := 0
		for _, enemy := range hitEnemies {
			if hits >= stats.Pierce && stats.Pierce > 0 {
				break
			}
			enemy.TakeDamage(stats.Damage)
			enemy.AddKnockback(&player.Pos, stats.Knockback)
			hits++
		}

		// Start the animation
		rp.displayRoll = true
		rp.currentFrame = 0
		rp.frameTimer = 0
		rp.hitDirection = player.GetFacingDirection()

		// Play the sound
		sfxPlayer := assets.AudioContext.NewPlayerFromBytes(rp.rollSound)
		sfxPlayer.Play()
	}
}

func (rp *RollingPin) Draw(screen *ebiten.Image, player *player.Player, mapOffsetX float64, mapOffsetY float64) {
	if rp.rollImage == nil {
		return
	}

	if rp.displayRoll {
		sx := rp.currentFrame * frameWidth
		sy := 0
		frameRect := image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)

		bounds := rp.rollImage.Bounds()
		if !frameRect.In(bounds) {
			fmt.Printf("Error: Frame rect %v out of bounds %v\n", frameRect, bounds)
			return
		}
		frameImage := rp.rollImage.SubImage(frameRect).(*ebiten.Image)

		pivotX := float64(frameWidth) / 2.0
		pivotY := float64(frameHeight) / 2.0
		angle := math.Atan2(rp.hitDirection.Y, rp.hitDirection.X)

		stats := rp.CurrentStats(player)
		currentRadius := baseRollingPinRadius * stats.AreaSize
		calculatedScale := currentRadius / float64(frameWidth)

		finalScreen := component.NewVector2D(player.Pos.X-mapOffsetX, player.Pos.Y-mapOffsetY)
		// Moving the animation to the outside of the player
		finalScreen = finalScreen.Add(rp.hitDirection.Mul(frameWidth / 2)) // TODO i dont know if frameWidth / 2 is the correct value

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
