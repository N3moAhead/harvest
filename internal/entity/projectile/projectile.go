package projectile

import (
	"fmt"
	"math"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/collision"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity/enemy"
	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile interface {
	Draw(screen *ebiten.Image, camX, camY float64)
	Update(enemies []enemy.EnemyInterface) (active bool)
	PlayImpactSound()
}

type BaseProjectile struct {
	Pos             component.Vector2D
	Dir             component.Vector2D
	Speed           float64
	Damage          float64
	HitRadius       float64
	Img             *ebiten.Image
	Pierce          int
	HittedEnemies   int
	Knockback       float64
	alreadyPierced  map[enemy.EnemyInterface]bool // making sure enemies are just getting pierced once
	FlyUntil        time.Time
	ImpactSoundName string
}

func NewBaseProjectile(
	pos,
	dir component.Vector2D,
	speed float64,
	img *ebiten.Image,
	duration time.Duration,
	dmg float64,
	hitRadius float64,
	pierce int,
	knockback float64,
	impactSoundName string,
) *BaseProjectile {
	return &BaseProjectile{
		Pos:             pos,
		Dir:             dir.Normalize(),
		Speed:           speed,
		Img:             img,
		Damage:          dmg,
		HitRadius:       hitRadius,
		Knockback:       knockback,
		Pierce:          pierce,
		HittedEnemies:   0,
		alreadyPierced:  make(map[enemy.EnemyInterface]bool),
		FlyUntil:        time.Now().Add(duration),
		ImpactSoundName: impactSoundName,
	}
}

func (b *BaseProjectile) Update(enemies []enemy.EnemyInterface) (active bool) {
	if time.Now().After(b.FlyUntil) {
		return false // The Projectile has run out of time
	}

	hitEnemies := collision.FindEnemiesInCircle(b.Pos, b.HitRadius, enemies)
	for _, enemy := range hitEnemies {
		// Do not pierce more enmies then allowed
		if b.HittedEnemies >= b.Pierce {
			return false
		}

		// Projectiles should pierce enemies only once
		if _, ok := b.alreadyPierced[enemy]; ok {
			continue
		}

		enemy.TakeDamage(b.Damage)
		enemy.AddKnockback(&b.Pos, b.Knockback)
		impactSound, ok := assets.AssetStore.GetSFXData(b.ImpactSoundName)
		if ok {
			sfxPlayer := assets.AudioContext.NewPlayerFromBytes(impactSound)
			sfxPlayer.Play()
		}

		b.HittedEnemies++
		b.alreadyPierced[enemy] = true
	}

	b.Pos = b.Pos.Add(b.Dir.Mul(b.Speed))
	return true
}

func (b *BaseProjectile) Draw(screen *ebiten.Image, camX, camY float64) {
	op := &ebiten.DrawImageOptions{}
	rotation := math.Atan2(b.Dir.Y, b.Dir.X)
	bounds := b.Img.Bounds()
	pivotX := bounds.Dx() / 2
	pivotY := bounds.Dy() / 2
	drawX := b.Pos.X - camX
	drawY := b.Pos.Y - camY
	op.GeoM.Translate(-float64(pivotX), -float64(pivotY))
	op.GeoM.Rotate(rotation)
	op.GeoM.Translate(drawX, drawY)
	screen.DrawImage(b.Img, op)
}

func (b *BaseProjectile) PlayImpactSound() {
	fmt.Println("Warning PlayImpactSound is not yet implemented")
}

var _ Projectile = (*BaseProjectile)(nil)
