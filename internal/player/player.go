package player

import (
	"image/color"

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
func (p *Player) Draw(screen *ebiten.Image) {
	rectSize := 32
	halfRectSize := float32(rectSize / 2)
	vector.DrawFilledRect(
		screen,
		float32(p.Pos.X)-halfRectSize,
		float32(p.Pos.Y)-halfRectSize,
		float32(rectSize),
		float32(rectSize),
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		true,
	)
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	p := &Player{
		Speed:  5,
		Pos:    component.NewVector2D(config.SCREEN_WIDTH/2, config.SCREEN_HEIGHT/2),
		Health: component.NewHealth(100),
	}
	return p
}
