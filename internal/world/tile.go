package world

import "github.com/hajimehoshi/ebiten/v2"

type TileType int

type Tile struct {
	Type       TileType
	FloorImage *ebiten.Image
	DecorImage *ebiten.Image
	IsWalkable bool
}

const (
	GrassMiddle TileType = iota
)

func (t *Tile) Draw(screen *ebiten.Image, posX float64, posY float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(posX, posY)
	if t.FloorImage != nil {
		screen.DrawImage(t.FloorImage, op)
	}
	if t.DecorImage != nil {
		screen.DrawImage(t.DecorImage, op)
	}
}
