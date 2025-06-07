package toast

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Toast struct {
	Width, Height int
	Text          string
	Font          font.Face
	DisplayUntil  time.Time
	Color         color.Color
}

func newToast(txt string, fnt font.Face, duration time.Duration) *Toast {
	txtBound := text.BoundString(fnt, txt)
	toast := &Toast{
		Width:        txtBound.Dx(),
		Height:       txtBound.Dy(),
		Text:         txt,
		Font:         fnt,
		DisplayUntil: time.Now().Add(duration),
		Color:        color.White,
	}
	return toast
}

func (toast *Toast) Update() (isAlive bool) {
	if time.Now().After(toast.DisplayUntil) {
		return false
	}
	return true
}

func (toast *Toast) Draw(screen *ebiten.Image, paddingTop int) {
	if toast.Text != "" && toast.Font != nil {
		// Center the text on Screen
		screenBounds := screen.Bounds()
		drawX := (screenBounds.Dx() - toast.Width) / 2.0
		drawY := 150 + paddingTop + int(toast.Font.Metrics().Ascent/64)
		text.Draw(screen, toast.Text, toast.Font, drawX, drawY, toast.Color)
	}
}
