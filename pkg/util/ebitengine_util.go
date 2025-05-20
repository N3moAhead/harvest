package util

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// GetSubImage returns the sub-image corresponding to the given frameRect
func GetSubImage(img *ebiten.Image, subRect image.Rectangle) *ebiten.Image {
	bounds := img.Bounds()
	if !subRect.In(bounds) {
		fmt.Printf("subRect %v is not within bounds %v\n", subRect, bounds)
		return nil
	}
	subImg := img.SubImage(subRect).(*ebiten.Image)
	return subImg
}
