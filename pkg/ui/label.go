package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Label struct {
	BaseElement
	Text  string
	Font  font.Face
	Color color.Color
}

func NewLabel(x, y float64, txt string, fnt font.Face, clr color.Color) *Label {
	textBound := text.BoundString(fnt, txt)
	lbl := &Label{
		BaseElement: *NewBaseElement(x, y, float64(textBound.Dx()), float64(textBound.Dy())),
		Text:        txt,
		Font:        fnt,
		Color:       clr,
	}
	return lbl
}

func (lbl *Label) Update(input *InputState) {
	lbl.BaseElement.Update(input)
}

func (lbl *Label) Draw(screen *ebiten.Image) {
	if !lbl.Visible {
		return
	}

	if lbl.Text != "" && lbl.Font != nil {
		textY := lbl.Y + float64(lbl.Font.Metrics().Ascent/64)
		text.Draw(screen, lbl.Text, lbl.Font, int(lbl.X), int(textY), lbl.Color)
	}
	lbl.BaseElement.Draw(screen)
}

var _ UIElement = (*Label)(nil)
