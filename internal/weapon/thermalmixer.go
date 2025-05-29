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
	"github.com/N3moAhead/harvest/internal/entity/item/itemtype"
	"github.com/N3moAhead/harvest/internal/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	baseThermalmixerRadius  = 150.0   // Basic range of the hit in pixels
	thermalmixerAttackAngle = 2*math.Pi // 360 degrees
	thermalmixerFrameWidth              = 64      // The height of each slash animation frame
	thermalmixerFrameHeight             = 64      // The width of each slash animation frame
	thermalmixerFrameCount              = 8       // The amount of existing frames
	thermalmixerAnimationSpeed          = 6       // Ticks for each frame
)

type Thermalmixer struct {
	BaseWeapon
	slashImage   *ebiten.Image
	slashSound   []byte
	frameTimer   int
	currentFrame int
	displaySlash bool
	hitDirection component.Vector2D
}

func NewThermalmixer() *Thermalmixer {
	stats := []WeaponStats{
		// Level 1
		{
			Damage:    1,
			Cooldown:  2 * time.Second,
			AreaSize:  1.0,
			Pierce:    1000,
			Knockback: 350.0,
		},
		// Level 2
		{
			Damage:    2,
			Cooldown:  1500 * time.Millisecond,
			AreaSize:  1.1,
			Pierce:    1000,
			Knockback: 420.0,
		},
		// Level 3
		{
			Damage:    5,
			Cooldown:  650 * time.Millisecond,
			AreaSize:  1.2,
			Pierce:    1000,
			Knockback: 490.0,
		},
	}

	slashImage, ok := assets.AssetStore.GetImage("thermalmixer_slash")
	if !ok {
		fmt.Println("Warning: Thermalmixer slash image not found")
	}
	slashSound, ok := assets.AssetStore.GetSFXData("thermalmixer_slash")
	if !ok {
		fmt.Println("Warning: Thermalmixer slash sound not found")
	}
	return &Thermalmixer{
		BaseWeapon: BaseWeapon{
			name:          "Thermalmixer",
			description:   "A mighty thermal mixer! there are sayings Gordan Ramsey touched it once...",
			cooldownTimer: 0,
			level:         1,
			maxLevel:      len(stats),
			statsPerLevel: stats,
			itemType:      itemtype.Thermalmixer,
		},
		slashImage: slashImage,
		slashSound: slashSound,
	}
}

func (t *Thermalmixer) Update(player *player.Player, enemies []enemy.EnemyInterface, dt time.Duration) {
	// Update the cooldown
	canAttack := t.UpdateCooldown(dt)

	// Update the animation frames
	t.frameTimer++
	if t.frameTimer >= animationSpeed {
		t.frameTimer = 0 // Reset the frame timer
		if t.currentFrame == frameCount-1 {
			t.displaySlash = false
		} else {
			t.currentFrame = (t.currentFrame + 1) % frameCount
		}
	}

	if canAttack {
		t.ResetCooldown(player) // Reset the cooldown before the attack

		stats := t.CurrentStats(player)
		playerPos := player.Pos
		facingDir := player.GetFacingDirection()

		// Calculate the thermal mixer range
		currentRadius := baseThermalmixerRadius * stats.AreaSize

		// Find the enemies in thermal mixer hit range
		hitEnemies := collision.FindEnemiesInArc(
			playerPos,
			currentRadius,
			facingDir,
			thermalmixerAttackAngle,
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
			// TODO play thermalmixer_hit sound
		}

		// Start the animation
		t.displaySlash = true
		t.currentFrame = 0
		t.frameTimer = 0
		t.hitDirection = player.GetFacingDirection()

		// Play the sound
		sfxPlayer := assets.AudioContext.NewPlayerFromBytes(t.slashSound)
		sfxPlayer.Play()
	}
}

func (t *Thermalmixer) Draw(screen *ebiten.Image, player *player.Player, mapOffsetX float64, mapOffsetY float64) {
    if t.slashImage == nil {
        return
    }

    if t.displaySlash {
        sx := t.currentFrame * thermalmixerFrameWidth
        sy := 0
        frameRect := image.Rect(sx, sy, sx+thermalmixerFrameWidth, sy+thermalmixerFrameHeight)

        bounds := t.slashImage.Bounds()
        if !frameRect.In(bounds) {
            fmt.Printf("Error: Frame rect %v out of bounds %v\n", frameRect, bounds)
            return
        }
        frameImage := t.slashImage.SubImage(frameRect).(*ebiten.Image)

        pivotX := float64(thermalmixerFrameWidth) / 2.0
        pivotY := float64(thermalmixerFrameHeight) / 2.0
        angle := math.Atan2(t.hitDirection.Y, t.hitDirection.X)

        stats := t.CurrentStats(player)
        currentRadius := baseThermalmixerRadius * stats.AreaSize
        calculatedScale := currentRadius / float64(thermalmixerFrameWidth)

        // Calculate sprite position on the screen
        finalScreen := component.NewVector2D(player.Pos.X-mapOffsetX, player.Pos.Y-mapOffsetY)

        op := &ebiten.DrawImageOptions{}

        // Centralize the sprite
        op.GeoM.Translate(-pivotX, -pivotY)
        if calculatedScale != 1.0 {
            op.GeoM.Scale(calculatedScale, calculatedScale)
        }
        op.GeoM.Rotate(angle)
        op.GeoM.Translate(finalScreen.X, finalScreen.Y)

        screen.DrawImage(frameImage, op)
    }
}
