package player

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	Pos    component.Vector2D
	Speed  float64
	Health component.Health
}

// The player is currently just drawn as a rectangle.
// TODO: Draw the player with assets
func (p *Player) Draw(screen *ebiten.Image, assetStore *assets.Store, mapOffsetX float64, mapOffsetY float64) {
	// TODO move the player rect size to the config or somewhere else
	rectSize := 32.0
	var halfRectSize float64 = rectSize / 2

	if playerImg, ok := assetStore.GetImage("player"); ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(p.Pos.X-mapOffsetX-halfRectSize, p.Pos.Y-mapOffsetY-halfRectSize)
		screen.DrawImage(playerImg, op)
	} else {
		// Fallback if player image could not be loaded
		vector.DrawFilledRect(
			screen,
			float32(p.Pos.X)-float32(mapOffsetX)-float32(halfRectSize),
			float32(p.Pos.Y)-float32(mapOffsetY)-float32(halfRectSize),
			float32(rectSize),
			float32(rectSize),
			color.RGBA{R: 255, G: 255, B: 255, A: 255},
			true,
		)
	}
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	p := &Player{
		Speed:  config.INITIAL_PLAYER_SPEED,
		Pos:    component.NewVector2D(config.SCREEN_WIDTH/2, config.SCREEN_HEIGHT/2),
		Health: component.NewHealth(100),
	}
	return p
}
