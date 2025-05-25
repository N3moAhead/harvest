package scene

import (
	"image/color"
	"log"

	"github.com/N3moAhead/harvest/internal/assets"
	"github.com/N3moAhead/harvest/internal/config"
	"github.com/N3moAhead/harvest/pkg/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type LoadingScene struct {
	BaseScene
	uiManager   *ui.UIManager
	loadingDone chan struct{}
}

func NewLoadingScene() *LoadingScene {
	fontFace, ok := assets.AssetStore.GetFont("2p")
	if !ok {
		panic("Unable to load font '2p' in NewLoadingScene")
	}
	text := ui.NewLabel(50, 50, "Loading...", fontFace, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	textWidth, textHeight := text.GetSize()
	offsetX := (config.SCREEN_WIDTH - textWidth) / 2
	offsetY := (config.SCREEN_HEIGHT - textHeight) / 2
	drawX := offsetX
	drawY := offsetY - float64(text.Font.Metrics().Ascent/64)
	text.SetPosition(drawX, drawY)

	newUiManager := ui.NewUIManager()
	newUiManager.AddElement(text)

	newLoadingScene := &LoadingScene{
		BaseScene:   *NewBaseScene(),
		uiManager:   newUiManager,
		loadingDone: make(chan struct{}),
	}

	// Loading all assets in a goroutine to not block the ui updates
	go func() {
		assets.LoadAllAssets()
		log.Println("All assets loaded.")
		close(newLoadingScene.loadingDone)
	}()

	return newLoadingScene
}

func (l *LoadingScene) Draw(screen *ebiten.Image) {
	l.uiManager.Draw(screen)
}

func (l *LoadingScene) Update() error {
	l.uiManager.Update()

	if l.loadingDone != nil {
		select {
		case <-l.loadingDone:
			log.Println("Loading finished, setting LoadingScene.isRunning to false.")
			l.SetIsRunning(false)
			l.loadingDone = nil
		default:
		}
	}

	return nil
}

var _ Scene = (*LoadingScene)(nil)
