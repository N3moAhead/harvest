package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct {
	BaseElement
	image *ebiten.Image
}

func NewImage(x, y float64, img *ebiten.Image) *Image {
	bounds := img.Bounds()

	newImage := &Image{
		BaseElement: *NewBaseElement(x, y, float64(bounds.Dx()), float64(bounds.Dy())),
		image:       img,
	}
	return newImage
}

func (img *Image) Update(input *InputState) {
	img.BaseElement.Update(input)
}

func (img *Image) Draw(screen *ebiten.Image) {
	if !img.Visible {
		return
	}

	if img.image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(img.X, img.Y)
		screen.DrawImage(img.image, op)
	}
	img.BaseElement.Draw(screen)
}

var _ UIElement = (*Image)(nil)
