package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputState struct {
	MouseX, MouseY         int
	MouseButtonLeftPressed bool
}

type UIElement interface {
	Update(input *InputState)
	Draw(screen *ebiten.Image)
	SetPosition(x, y float64)
	GetPosition() (x, y float64)
	GetSize() (width, height float64)
	SetVisible(visible bool)
	IsVisible() bool
	SetEnabled(enabled bool)
	IsEnabled() bool
	AddChild(child UIElement)
	HandleInput(input *InputState)
	GetChildren() []UIElement
}

type BaseElement struct {
	X, Y          float64
	Width, Height float64
	Visible       bool
	Enabled       bool
	Children      []UIElement
	Parent        UIElement
}

func NewBaseElement(x, y, width, height float64) *BaseElement {
	return &BaseElement{
		X: x, Y: y,
		Width: width, Height: height,
		Visible:  true,
		Enabled:  true,
		Children: make([]UIElement, 0),
	}
}

func (b *BaseElement) SetPosition(x, y float64) {
	b.X = x
	b.Y = y
}

func (b *BaseElement) GetPosition() (float64, float64) {
	return b.X, b.Y
}

func (b *BaseElement) GetSize() (float64, float64) {
	return b.Width, b.Height
}

func (b *BaseElement) SetVisible(visible bool) {
	b.Visible = visible
}

func (b *BaseElement) IsVisible() bool {
	return b.Visible
}

func (b *BaseElement) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

func (b *BaseElement) IsEnabled() bool {
	return b.Enabled
}

func (b *BaseElement) Update(input *InputState) {
	if !b.Visible || !b.Enabled {
		return
	}
	for _, child := range b.Children {
		child.Update(input)
	}
}

func (b *BaseElement) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}
	for _, child := range b.Children {
		child.Draw(screen)
	}
}

func (b *BaseElement) HandleInput(input *InputState) {
	if !b.Visible || !b.Enabled {
		return
	}

	for _, child := range b.Children {
		child.HandleInput(input)
	}
}

func (b *BaseElement) AddChild(child UIElement) {
	b.Children = append(b.Children, child)
	if c, ok := child.(*BaseElement); ok {
		c.Parent = b
	}
}

func (b *BaseElement) GetChildren() []UIElement {
	return b.Children
}

func (b *BaseElement) IsMouseOver(mouseX, mouseY int) bool {
	return mouseX >= int(b.X) && mouseX < int(b.X+b.Width) &&
		mouseY >= int(b.Y) && mouseY < int(b.Y+b.Height)
}

// Check that BaseElement correctly implements UIElement
var _ UIElement = (*BaseElement)(nil)
