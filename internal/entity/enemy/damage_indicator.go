package enemy

import (
	"fmt"
	"image/color"
	"time"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type DamageIndicator struct {
	Damage     float64
	Pos        component.Vector2D
	Offset     component.Vector2D
	Dir        component.Vector2D
	Speed      float64
	AliveUntil time.Time
	Font       font.Face
}

func NewDamageIndicator(pos component.Vector2D, dir component.Vector2D, dmg float64) *DamageIndicator {
	fontFace, _ := assets.AssetStore.GetFont("micro")
	return &DamageIndicator{
		Damage:     dmg,
		Pos:        pos,
		Offset:     component.NewVector2D(0, 0),
		Dir:        dir.Normalize(),
		Speed:      config.DAMAGE_INDICATOR_SPEED,
		AliveUntil: time.Now().Add(config.DAMAGE_INDICATOR_DURATION),
		Font:       fontFace,
	}
}

func (d *DamageIndicator) Update(pos component.Vector2D) (isAlive bool) {
	if time.Now().After(d.AliveUntil) {
		return false
	}
	d.Pos = pos
	d.Offset = d.Offset.Add(d.Dir.Mul(d.Speed))
	return true
}

func (d *DamageIndicator) Draw(screen *ebiten.Image, camX, camY float64) {
	if d.Font != nil {
		textY := d.Pos.Y + d.Offset.Y - camY + float64(d.Font.Metrics().Ascent/64)
		text.Draw(screen, fmt.Sprintf("%d", int(d.Damage)), d.Font, int(d.Pos.X-camX+d.Offset.X), int(textY), color.White)
	} else {
		fmt.Println("Warning: Missing font in damage indicator Draw")
	}
}
