package scene

import (
	"image/color"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/internal/hud"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScoreScene struct {
	BaseScene
	uiManager *ui.UIManager
}

func NewScoreScene(score int) *ScoreScene {

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in new base scene")
	}

	text := ui.NewLabel(50, 50, "You Died!", fontFace, color.RGBA{R: 255, G: 255, B: 255, A: 255})
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

	endSceneButton := ui.NewButton(0, 0, 400, 50, "Back to Menu", fontFace, func() { newScoreScene.SetIsRunning(false) })
	btnWidth, btnHeight := endSceneButton.GetSize()
	btnDrawX := (config.SCREEN_WIDTH - btnWidth) / 2
	btnDrawY := (config.SCREEN_HEIGHT - btnHeight) / 2
	btnDrawY += 75 // Move the button below the menu text
	endSceneButton.SetPosition(btnDrawX, btnDrawY)

	newScoreDisplay := hud.NewScoreDisplay(&score, "Score")

	newUiManager.AddElement(text)
	newUiManager.AddElement(endSceneButton)
	newUiManager.AddElement(newScoreDisplay)

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
