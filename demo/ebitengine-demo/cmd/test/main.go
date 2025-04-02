package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 800
	SquareSize   = 5
)

type Block struct {
	Img    *ebiten.Image
	X, Y   float64
	dx, dy float64
}

type Game struct {
	blocks []Block
}

func (g *Game) Update() error {
	for i := range g.blocks {
		block := &g.blocks[i]

		block.X += block.dx
		block.Y += block.dy

		if block.X <= 0 || block.X+SquareSize >= ScreenWidth {
			block.dx = -block.dx
		}
		if block.Y <= 0 || block.Y+SquareSize >= ScreenHeight {
			block.dy = -block.dy
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.blocks {
		block := &g.blocks[i]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(block.X, block.Y)
		screen.DrawImage(block.Img, op)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func randRange(min, max int) float64 {
	return float64(rand.IntN(max-min) + min)
}

func main() {
	game := &Game{
		blocks: []Block{},
	}
	for i := range 5000 {
		game.blocks = append(game.blocks, Block{
			Img: ebiten.NewImage(SquareSize, SquareSize),
			X:   randRange(0, ScreenWidth-SquareSize),
			Y:   randRange(0, ScreenHeight-SquareSize),
			dx:  (rand.Float64() + 1) * 2,
			dy:  (rand.Float64() + 1) * 2,
		})
		game.blocks[i].Img.Fill(color.RGBA{uint8(rand.Int()), uint8(rand.Int()), uint8(rand.Int()), 255})
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Harvest Ebiten Demo")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
