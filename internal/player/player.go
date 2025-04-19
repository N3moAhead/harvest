package player

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/component"
	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	entity.Entity

	Speed           float64
	MagnetRadius    float64
	Health          component.Health
	FacingDirection component.Vector2D
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

func (p *Player) GetFacingDirection() component.Vector2D {
	// Fallback if the facing direction is not defined
	if p.FacingDirection.LengthSq() == 0 {
		return component.Vector2D{X: 0, Y: -1}
	}
	// Its important that this direction is normalized
	// but im making sure it is when setting it in the game.go file
	return p.FacingDirection
}

// TODO: implement a LoadPlayer function to get the saved
// game state from the past
func NewPlayer() *Player {
	baseEntity := entity.NewEntity(
		(config.WIDTH_IN_TILES*config.TILE_SIZE)/2,
		(config.HEIGHT_IN_TILES*config.TILE_SIZE)/2,
	)

	p := &Player{
		Entity:          *baseEntity,
		MagnetRadius:    config.INITIAL_PLAYER_MAGNET_RADIUS,
		Speed:           config.INITIAL_PLAYER_SPEED,
		Health:          component.NewHealth(100),
		FacingDirection: component.NewVector2D(0, -1), // Default looks up
	}
	return p
}
