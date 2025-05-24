package scene

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScoreScene struct {
	BaseScene
	uiManager *ui.UIManager
}

func NewScoreScene() *ScoreScene {

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in new base scene")
	}
	text := ui.NewLabel(50, 50, "Score", fontFace, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// Center the text
	textWidth, textHeight := text.GetSize()
	offsetX := (config.SCREEN_WIDTH - textWidth) / 2
	offsetY := (config.SCREEN_HEIGHT - textHeight) / 2
	drawX := offsetX
	drawY := offsetY - float64(text.Font.Metrics().Ascent/64)
	text.SetPosition(drawX, drawY)
	newUiManager := ui.NewUIManager()
	newScoreScene := &ScoreScene{
		BaseScene: *NewBaseScene(),
		uiManager: newUiManager,
	}
	endSceneButton := ui.NewButton(0, 0, 100, 50, "Next", fontFace, func() { newScoreScene.SetIsRunning(false) })

	newUiManager.AddElement(text)
	newUiManager.AddElement(endSceneButton)

	return newScoreScene
}

func (l *ScoreScene) Draw(screen *ebiten.Image) {
	l.uiManager.Draw(screen)
}

func (l *ScoreScene) Update() error {
	l.uiManager.Update()
	return nil
}

var _ Scene = (*LoadingScene)(nil)
