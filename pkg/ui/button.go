package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Button struct {
	BaseElement
	Text                 string
	Font                 font.Face
	TextColor            color.Color
	BackgroundColor      color.Color
	HoverBackgroundColor color.Color
	OnClick              func()
	isHovered            bool
}

func NewButton(x, y, width, height float64, txt string, fnt font.Face, onClick func()) *Button {
	btn := &Button{
		BaseElement:          *NewBaseElement(x, y, width, height),
		Text:                 txt,
		Font:                 fnt,
		TextColor:            color.RGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor:      color.RGBA{R: 50, G: 50, B: 50, A: 255},
		HoverBackgroundColor: color.RGBA{R: 80, G: 80, B: 80, A: 255},
		OnClick:              onClick,
	}
	return btn
}

func (btn *Button) Update(input *InputState) {
	if !btn.Visible || !btn.Enabled {
		btn.isHovered = false
		return
	}

	btn.isHovered = btn.IsMouseOver(input.MouseX, input.MouseY)
	// A Button should not have any children. But a Button could have children
	// so im going to update it if there are any...
	btn.BaseElement.Update(input)
}

func (btn *Button) Draw(screen *ebiten.Image) {
	if !btn.Visible {
		return
	}

	bgColor := btn.BackgroundColor
	if !btn.Enabled {
		bgColor = color.RGBA{R: 30, G: 30, B: 30, A: 255}
	} else if btn.isHovered {
		bgColor = btn.HoverBackgroundColor
	}

	vector.DrawFilledRect(screen, float32(btn.X), float32(btn.Y), float32(btn.Width), float32(btn.Height), bgColor, false)

	if btn.Text != "" && btn.Font != nil {
		textBounds := text.BoundString(btn.Font, btn.Text)
		offsetX := (btn.Width - float64(textBounds.Dx())) / 2
		offsetY := (btn.Height - float64(textBounds.Dy())) / 2
		drawX := btn.X + offsetX
		drawY := btn.Y + offsetY - float64(textBounds.Min.Y)
		text.Draw(screen, btn.Text, btn.Font, int(drawX), int(drawY), btn.TextColor)
	}

	// A Button should again not have any children. But a Button could have children
	// so im going to draw them if there are any.
	btn.BaseElement.Draw(screen)
}

func (btn *Button) HandleInput(input *InputState) {
	if !btn.Visible || !btn.Enabled {
		return
	}

	if btn.IsMouseOver(input.MouseX, input.MouseY) && input.MouseButtonLeftPressed {
		if btn.OnClick != nil {
			btn.OnClick()
		}
	}
}

var _ UIElement = (*Button)(nil)
