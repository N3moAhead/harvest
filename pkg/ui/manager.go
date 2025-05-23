package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type UIManager struct {
	elements []UIElement
}

func NewUIManager() *UIManager {
	return &UIManager{
		elements: make([]UIElement, 0),
	}
}

func (m *UIManager) AddElement(element UIElement) {
	m.elements = append(m.elements, element)
}

func (m *UIManager) RemoveElement(element UIElement) {
	for i, e := range m.elements {
		if e == element {
			m.elements = append(m.elements[:i], m.elements[i+1:]...)
			return
		}
	}
}

func (m *UIManager) ClearElements() {
	m.elements = make([]UIElement, 0)
}

func (m *UIManager) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	inputState := &InputState{
		MouseX:                 mouseX,
		MouseY:                 mouseY,
		MouseButtonLeftPressed: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
	}

	for i := len(m.elements) - 1; i >= 0; i-- {
		el := m.elements[i]
		if el.IsVisible() {
			el.Update(inputState)
		}
	}

	for i := len(m.elements) - 1; i >= 0; i-- {
		el := m.elements[i]
		if el.IsVisible() && el.IsEnabled() {
			el.HandleInput(inputState)
		}
	}
}

func (m *UIManager) Draw(screen *ebiten.Image) {
	for _, element := range m.elements {
		if element.IsVisible() {
			element.Draw(screen)
		}
	}
}
