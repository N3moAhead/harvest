package scene

import (
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	IsRunning() bool
	SetIsRunning(running bool)
}

type BaseScene struct {
	sceneRunning bool
}

func NewBaseScene() *BaseScene {
	return &BaseScene{
		sceneRunning: true,
	}
}

func (b *BaseScene) Update() error {
	return nil
}

func (b *BaseScene) Draw(screen *ebiten.Image) {}

func (b *BaseScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.SCREEN_WIDTH, config.SCREEN_HEIGHT
}

func (b *BaseScene) IsRunning() bool {
	return b.sceneRunning
}

func (b *BaseScene) SetIsRunning(running bool) {
	b.sceneRunning = running
}

var _ Scene = (*BaseScene)(nil)
var _ ebiten.Game = (*BaseScene)(nil)
