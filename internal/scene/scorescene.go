package scene

import (
	"fmt"
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

func NewScoreScene(stats PlayerStats) *ScoreScene {

	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font in score scene")
	}
	microFont, ok := assets.AssetStore.GetFont("micro")
	if !ok {
		panic("Unable to load font in score scene")
	}

	text := ui.NewLabel(50, 50, "You got Veggienated!", fontFace, color.RGBA{R: 255, G: 255, B: 255, A: 255})
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

	newScoreDisplay := hud.NewScoreDisplay(&stats.lastGameScore, "Score")

	// Stats Display
	statsContainer := ui.NewContainer(10, 10, &ui.ContainerOptions{
		Direction: ui.Col,
		Gap:       10,
	})
	xpEarnedDisplay := ui.NewLabel(0, 0, fmt.Sprintf("XP earned: %d", stats.lastGameXPEarned), microFont, color.White)
	levelEarned := ui.NewLabel(0, 0, fmt.Sprintf("Level earned: %d", uint(stats.lastGameXPEarned/10)), microFont, color.White)
	//lastWave := ui.NewLabel(0, 0, fmt.Sprintf("Made it to wave: %d", uint(stats.currentWaveIndex/10)), microFont, color.White) // Could be cool
	statsContainer.AddChild(xpEarnedDisplay)
	statsContainer.AddChild(levelEarned)
	newUiManager.AddElement(statsContainer)

	newUiManager.AddElement(text)
	newUiManager.AddElement(endSceneButton)
	newUiManager.AddElement(newScoreDisplay)

	sound, ok := assets.AssetStore.GetSFXData("veggienated")
	if ok {
		sfxPlayer := assets.AudioContext.NewPlayerFromBytes(sound)
		sfxPlayer.Play()
	} else {
		fmt.Println("Warning: Could not load Veggienated sound")
	}

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
