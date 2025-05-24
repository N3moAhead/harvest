package scene

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/pkg/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type LoadingScene struct {
	BaseScene
	uiManager *ui.UIManager
}

func NewLoadingScene() *LoadingScene {

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in new base scene")
	}
	text := ui.NewLabel(50, 50, "Insert Coin", fontFace, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// Center the text
	textWidth, textHeight := text.GetSize()
	offsetX := (config.SCREEN_WIDTH - textWidth) / 2
	offsetY := (config.SCREEN_HEIGHT - textHeight) / 2
	drawX := offsetX
	drawY := offsetY - float64(text.Font.Metrics().Ascent/64)
	text.SetPosition(drawX, drawY)
	newUiManager := ui.NewUIManager()
	newLoadingScene := &LoadingScene{
		BaseScene: *NewBaseScene(),
		uiManager: newUiManager,
	}
	endSceneButton := ui.NewButton(0, 0, 100, 50, "Next", fontFace, func() { newLoadingScene.SetIsRunning(false) })

	newUiManager.AddElement(text)
	newUiManager.AddElement(endSceneButton)

	return newLoadingScene
}

func (l *LoadingScene) Draw(screen *ebiten.Image) {
	l.uiManager.Draw(screen)
}

func (l *LoadingScene) Update() error {
	l.uiManager.Update()
	return nil
}

var _ Scene = (*LoadingScene)(nil)
