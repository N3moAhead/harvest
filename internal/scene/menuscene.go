package scene

import (
	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	BaseScene
	uiManager *ui.UIManager
	icon      *ebiten.Image
	isRunning bool
}

func NewMenuScene() *MenuScene {
	icon, ok := assets.AssetStore.GetImage("menu-icon")
	if !ok {
		panic("menu-icon nicht im AssetStore gefunden")
	}

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in new base scene")
	}
	newUiManager := ui.NewUIManager()
	newMenuScene := &MenuScene{
		BaseScene: *NewBaseScene(),
		uiManager: newUiManager,
		icon:      icon,
		isRunning: true,
	}

	startBtn := ui.NewButton(0, 0, 150, 40, "Start", fontFace, func() { newMenuScene.SetIsRunning(false) })
	container := ui.NewContainer((config.SCREEN_WIDTH-150)/2, 350, &ui.ContainerOptions{
		Direction: ui.Col,
		// Gap:       0,
	})
	container.AddChild(startBtn)
	newUiManager.AddElement(container)

	return newMenuScene
}

func (l *MenuScene) Draw(screen *ebiten.Image) {
	screenWidth := screen.Bounds().Dx()
	iconWidth := l.icon.Bounds().Dx()
	scale := 0.14 // Scale the icon
	scaledW := float64(iconWidth) * scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		(float64(screenWidth)-scaledW)/2, // horziontal center
		30,                               // 30px from top
	)
	screen.DrawImage(l.icon, op)

	l.uiManager.Draw(screen)
}

func (l *MenuScene) Update() error {
	l.uiManager.Update()
	return nil
}

func (m *MenuScene) IsRunning() bool {
	return m.isRunning
}

func (m *MenuScene) SetIsRunning(r bool) {
	m.isRunning = r
}

var _ Scene = (*MenuScene)(nil)
