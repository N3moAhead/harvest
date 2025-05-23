package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputState struct {
	Up                     bool
	Right                  bool
	Down                   bool
	Left                   bool
	MouseX, MouseY         int
	MouseButtonLeftPressed bool
}

func GetInputState() *InputState {
	mouseX, mouseY := ebiten.CursorPosition()
	return &InputState{
		Up:                     ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp),
		Right:                  ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight),
		Down:                   ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown),
		Left:                   ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft),
		MouseX:                 mouseX,
		MouseY:                 mouseY,
		MouseButtonLeftPressed: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
	}
}
