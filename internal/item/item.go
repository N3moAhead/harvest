package item

import (
	"image/color"
	"math/rand/v2"

	"github.com/N3moAhead/harvest/internal/entity"
	"github.com/N3moAhead/harvest/internal/player"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Item struct {
	entity.Entity

	Type ItemType
}

func (i *Item) Update(player *player.Player) {
	// distance to player
	diff := player.Pos.Sub(i.Pos)
	len := diff.Len()
	if len < player.MagnetRadius {
		// Move towards player
		dir := diff.Normalize()
		i.Pos = i.Pos.Add(dir.Mul(config.PLAYER_MAGNET_ATTRACTION_SPEED))
	}
	// Check if the item is in the magnet radius of the player
	// if so move the item closer to the player
	// if the item is in the pickup radius of the player
	// move it into his inventory
}

func (i *Item) Draw(screen *ebiten.Image, mapOffsetX float64, mapOffsetY float64) {
	vector.DrawFilledRect(
		screen,
		float32(i.Pos.X)-float32(mapOffsetX)-16.0,
		float32(i.Pos.Y)-float32(mapOffsetY)-16.0,
		32.0,
		32.0,
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		true,
	)
}

func newItemBase() *Item {
	// For testing items will spawn randomly on the map
	// random float64 between 0 and 1280
	posX := rand.Float64() * config.WIDTH_IN_TILES * config.TILE_SIZE
	posY := rand.Float64() * config.HEIGHT_IN_TILES * config.TILE_SIZE
	baseClass := entity.NewEntity(posX, posY)
	return &Item{
		Entity: *baseClass,
	}
}
