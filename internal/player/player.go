package player

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	entity.Entity

	Speed  float64
	Health component.Health
}

// The player is currently just drawn as a rectangle.
// TODO: Draw the player with assets
func (p *Player) Draw(screen *ebiten.Image, mapOffsetX float64, mapOffsetY float64) {
	rectSize := 32
	halfRectSize := float32(rectSize / 2)
	vector.DrawFilledRect(
		screen,
		float32(p.Pos.X)-float32(mapOffsetX)-halfRectSize,
		float32(p.Pos.Y)-float32(mapOffsetY)-halfRectSize,
		float32(rectSize),
		float32(rectSize),
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		true,
	)
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	baseEntity := entity.NewEntity(config.SCREEN_WIDTH/2, config.SCREEN_HEIGHT/2)

	p := &Player{
		Entity: *baseEntity,
		Speed:  config.INITIAL_PLAYER_SPEED,
		Health: component.NewHealth(100),
	}
	return p
}
